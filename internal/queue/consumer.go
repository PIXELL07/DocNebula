package queue

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

const (
	MaxRetry         = 3
	ProcessingSuffix = ":processing"
)

type Consumer struct {
	Client      *redis.Client
	DLQ         string
	WorkerCount int
	Logger      *slog.Logger
}

func (c *Consumer) Start(ctx context.Context, queue string, handler func(context.Context, Message) error, producer *Producer) {
	if c.WorkerCount <= 0 {
		c.WorkerCount = 4
	}

	for i := 0; i < c.WorkerCount; i++ {
		go c.workerLoop(ctx, queue, handler, producer, i)
	}

	<-ctx.Done()
}
