package main

import (
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

	// Postgres
	pg, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		logger.Error("postgres connect failed", "err", err)
		os.Exit(1)
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	jobRepo := &repository.JobRepo{DB: pg}
	producer := &queue.Producer{Client: rdb}

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

	// health
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	logger.Info("api starting", "port", 8080)
	http.ListenAndServe(":8080", nil)
}
