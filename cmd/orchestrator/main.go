// Pipeline controller

package main

import (
	"DocNebula/internal/queue"
	"DocNebula/internal/repository"
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// ** For reference ** , The orchestrator: is responsible for kicking off jobs that are in UPLOADED state.
// In many real systems this would be event-driven (S3/SNS/etc), but this
// polling version is simple and production-reasonable for practice.

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
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})

	jobRepo := &repository.JobRepo{DB: pg}
	producer := &queue.Producer{Client: rdb}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	logger.Info("orchestrator started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			kickJobs(ctx, logger, jobRepo, producer)
		}
	}
}

// this function finds newly uploaded jobs and pushes them into the pipeline.......
func kickJobs(ctx context.Context, logger *slog.Logger, jobRepo *repository.JobRepo, producer *queue.Producer) {
	rows, err := jobRepo.DB.QueryContext(ctx, `
		SELECT id FROM jobs WHERE status='UPLOADED' LIMIT 50
	`)
	if err != nil {
		logger.Error("orchestrator query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var jobID string
		if err := rows.Scan(&jobID); err != nil {
			continue
		}

		msg := queue.Message{
			JobID:     jobID,
			Attempt:   0,
			Timestamp: time.Now(),
		}

		if err := producer.Publish(ctx, "unzip_queue", msg); err != nil {
			logger.Error("enqueue failed", "job_id", jobID, "err", err)
			continue
		}

		// mark job running
		_ = jobRepo.UpdateStatus(ctx, jobID, "RUNNING")

		logger.Info("job dispatched", "job_id", jobID)
	}
}
