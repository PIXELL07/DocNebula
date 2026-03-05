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

	logger.Info("unzip worker started")

	consumer := &queue.Consumer{
		Client:      rdb,
		DLQ:         "unzip_dlq",
		WorkerCount: 4,
		Logger:      logger,
	}

	producer := &queue.Producer{Client: rdb}

	handler := func(ctx context.Context, msg queue.Message) error {

		logger.Info("UNZIP start",
			"job_id", msg.JobID,
			"attempt", msg.Attempt,
		)

		// TODO: real unzip logic
		time.Sleep(500 * time.Millisecond)

		logger.Info("UNZIP finished",
			"job_id", msg.JobID,
		)

		return producer.Publish(ctx, "ocr_queue", msg)
	}

	consumer.Start(ctx, "unzip_queue", handler, producer)
}
