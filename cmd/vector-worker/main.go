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

	logger.Info("vector worker started")

	consumer := &queue.Consumer{
		Client:      rdb,
		DLQ:         "vector_dlq",
		WorkerCount: 4,
		Logger:      logger,
	}

	producer := &queue.Producer{Client: rdb}

	handler := func(ctx context.Context, msg queue.Message) error {

		logger.Info("VECTOR processing",
			"job_id", msg.JobID,
			"attempt", msg.Attempt,
		)

		// TODO: add embedding generation here

		logger.Info("VECTOR finished",
			"job_id", msg.JobID,
		)

		return nil
	}

	consumer.Start(ctx, "vector_queue", handler, producer)
}
