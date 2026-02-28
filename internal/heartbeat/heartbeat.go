// This package implements a simple heartbeat mechanism to track active workers in the system.
// Each worker periodically updates its last seen timestamp in the database, allowing the system to identify which workers are currently active and responsive.

package heartbeat

import (
	"context"
	"database/sql"
	"time"
)

type Heartbeater struct {
	DB       *sql.DB
	WorkerID string
}

func (h *Heartbeater) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				h.beat(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (h *Heartbeater) beat(ctx context.Context) {
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO worker_heartbeats(worker_id, last_seen)
		VALUES ($1, NOW())
		ON CONFLICT (worker_id)
		DO UPDATE SET last_seen = NOW()
	`, h.WorkerID)
}
