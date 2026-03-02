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
		DLQ:         "unzip_dlq",
		WorkerCount: 4,
		Logger:      logger,
	}
	producer := &queue.Producer{Client: rdb}

	handler := func(ctx context.Context, msg queue.Message) error {
		logger.Info("UNZIP start", "job_id", msg.JobID)

		// TODO: real unzip work
		time.Sleep(500 * time.Millisecond)

		return producer.Publish(ctx, "ocr_queue", msg)
	}

	consumer.Start(ctx, "unzip_queue", handler, producer)
}
