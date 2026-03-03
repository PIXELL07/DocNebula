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

	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})

	consumer := &queue.Consumer{
		Client:      rdb,
		DLQ:         "ocr_dlq",
		WorkerCount: 4,
		Logger:      logger,
	}
	producer := &queue.Producer{Client: rdb}

	handler := func(ctx context.Context, msg queue.Message) error {
		logger.Info("OCR start", "job_id", msg.JobID)

		// simulate OCR
		time.Sleep(700 * time.Millisecond)

		return producer.Publish(ctx, "vector_queue", msg)
	}

	consumer.Start(ctx, "ocr_queue", handler, producer)
}
