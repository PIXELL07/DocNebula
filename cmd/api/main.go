package main

import (
	httpx "DocNebula/internal/http"
	"DocNebula/internal/queue"
	"DocNebula/internal/repository"
	"DocNebula/internal/utils"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// PostgreSQL
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		logger.Error("POSTGRES_DSN not set")
		os.Exit(1)
	}

	pg, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("postgres open failed", "err", err)
		os.Exit(1)
	}

	if err := pg.Ping(); err != nil {
		logger.Error("postgres ping failed", "err", err)
		os.Exit(1)
	}

	logger.Info("postgres connected")

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("redis connection failed", "err", err)
		os.Exit(1)
	}

	logger.Info("redis connected")

	// Repositories
	jobRepo := &repository.JobRepo{DB: pg}
	userRepo := &repository.UserRepo{DB: pg}

	producer := &queue.Producer{Client: rdb}

	// http handlers
	authHandler := &httpx.AuthHandler{
		UserRepo: userRepo,
	}

	resetHandler := &httpx.ResetHandler{
		UserRepo: userRepo,
	}

	http.HandleFunc("/signup", authHandler.Signup)
	http.HandleFunc("/login", authHandler.Login)

	http.HandleFunc("/forgot-password", resetHandler.ForgotPassword)
	http.HandleFunc("/reset-password", resetHandler.ResetPassword)

	// Upload endpoint
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		idemKey := utils.FromHeaderOrRequest(r)

		logger.Info("upload request",
			"idempotency_key", idemKey,
			"remote", r.RemoteAddr,
		)

		job, err := jobRepo.Create(ctx, idemKey)
		if err != nil {
			logger.Error("job create failed", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := queue.Message{
			JobID:     job.ID,
			Attempt:   0,
			Timestamp: time.Now(),
		}

		if err := producer.Publish(ctx, "unzip_queue", msg); err != nil {
			logger.Error("enqueue failed", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Idempotency-Key", idemKey)

		json.NewEncoder(w).Encode(job)

		logger.Info("upload processed",
			"job_id", job.ID,
			"latency_ms", time.Since(start).Milliseconds(),
		)
	})

	// Health endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DocNebula API running"))
	})

	logger.Info("api starting", "port", 9000)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		logger.Error("server failed", "err", err)
	}
}

// Kubernetes / Docker can now probe these endpoints.

// additional endpoints:
// - GET /status?id=job_id -> job status
// - GET /results?id=job_id -> job results (once ready)
// - POST /retry?id=job_id -> trigger retry (for failed jobs)
// - GET /metrics -> Prometheus metrics (job counts, latencies, etc)
// - POST /cancel?id=job_id -> cancel a job (if still processing)
// fix: handle graceful shutdown, close db and redis connections, etc.. in internal/http/health.go
