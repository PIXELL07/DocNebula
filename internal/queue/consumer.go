package queue

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

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

func (c *Consumer) workerLoop(ctx context.Context, queue string, handler func(context.Context, Message) error, producer *Producer, workerID int) {
	processingQueue := queue + ProcessingSuffix

	for {
		// Visibility-timeout style pop
		res, err := c.Client.BRPopLPush(ctx, queue, processingQueue, 0).Result()
		if err != nil {
			c.Logger.Error("queue pop failed", "err", err)
			continue
		}

		var msg Message
		if err := json.Unmarshal([]byte(res), &msg); err != nil {
			c.Logger.Error("bad message -> DLQ")
			c.Client.LPush(ctx, c.DLQ, res)
			c.Client.LRem(ctx, processingQueue, 1, res)
			continue
		}

		// Per-job timeout protection
		jobCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		err = handler(jobCtx, msg)
		cancel()

		if err != nil {
			msg.Attempt++

			if msg.Attempt >= MaxRetry {
				c.Logger.Error("max retry exceeded -> DLQ", "job_id", msg.JobID)
				b, _ := json.Marshal(msg)
				c.Client.LPush(ctx, c.DLQ, b)
				c.Client.LRem(ctx, processingQueue, 1, res)
				continue
			}

			// exponential backoff
			time.Sleep(time.Duration(msg.Attempt*2) * time.Second)
			_ = producer.Publish(ctx, queue, msg)
			c.Client.LRem(ctx, processingQueue, 1, res)
			continue
		}

		// ACK (remove from processing queue)
		c.Client.LRem(ctx, processingQueue, 1, res)
		c.Logger.Info("message processed", "job_id", msg.JobID, "worker", workerID)
	}
}
