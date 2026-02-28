package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Message struct {
	JobID     string    `json:"job_id"`
	Attempt   int       `json:"attempt"`
	Timestamp time.Time `json:"ts"`
}

type Producer struct {
	Client *redis.Client
}

func (p *Producer) Publish(ctx context.Context, queue string, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.Client.LPush(ctx, queue, b).Err()
}
