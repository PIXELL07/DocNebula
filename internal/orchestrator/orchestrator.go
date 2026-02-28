// This is used to orchestrate the workflow of processing uploaded ZIP files.
// It sends messages to the appropriate queues to trigger each step of the pipeline(stp) (unzip → OCR → vectorization).
// The Orchestrator is responsible for starting the job by publishing a message to the "unzip_queue" with the job ID and timestamp.
package orchestrator

import (
	"DocNebula/internal/queue"
	"context"
	"time"
)

type Orchestrator struct{ Producer *queue.Producer }

func (o *Orchestrator) StartJob(ctx context.Context, jobID string) error {
	return o.Producer.Publish(ctx, "unzip_queue", queue.Message{
		JobID: jobID,
		TS:    time.Now(),
	})
}
