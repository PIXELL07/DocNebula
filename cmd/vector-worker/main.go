package main

import (
	"DocNebula/internal/queue"
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})

	consumer := &queue.Consumer{
		Client:      rdb,
		DLQ:         "vector_dlq",
		WorkerCount: 4,
		Logger:      logger,
	}

	producer := &queue.Producer{Client: rdb}

	handler := func(ctx context.Context, msg queue.Message) error {
		logger.Info("VECTOR done", "job_id", msg.JobID)
		return nil
	}

	consumer.Start(ctx, "vector_queue", handler, producer)
}
