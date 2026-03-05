// Old worker: cmd/ocr-worker
// job → OCR whole file
// in future, delete this file and move the logic to
// cmd/ocr-page-worker/main.go, and use a single worker for all tasks

package main

import (
	"DocNebula/internal/queue"
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("redis connection failed", "err", err)
		os.Exit(1)
	}

	logger.Info("ocr worker started")

	consumer := &queue.Consumer{
		Client:      rdb,
		DLQ:         "ocr_dlq",
		WorkerCount: 4,
		Logger:      logger,
	}

	producer := &queue.Producer{Client: rdb}

	handler := func(ctx context.Context, msg queue.Message) error {

		logger.Info("OCR start",
			"job_id", msg.JobID,
			"attempt", msg.Attempt,
		)

		// TODO: replace with real OCR engine
		time.Sleep(700 * time.Millisecond)

		logger.Info("OCR finished",
			"job_id", msg.JobID,
		)

		return producer.Publish(ctx, "vector_queue", msg)
	}

	consumer.Start(ctx, "ocr_queue", handler, producer)
}
