// New worker: cmd/ocr-page-worker
// New version
// page → OCR
// This worker will replace the old worker in the future, and handle all tasks related to OCR,
// including whole file OCR and page OCR. For now, it only handles page OCR, and the logic for whole file OCR will be moved here later.

package ocrpageworker

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"DocNebula/internal/queue"
	"DocNebula/internal/repository"
	"DocNebula/internal/workers/ocr"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// --- Postgres ---
	pg, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		logger.Error("postgres connect failed", "err", err)
		os.Exit(1)
	}

	pageRepo := &repository.PageRepo{DB: pg}
	fileRepo := &repository.FileRepo{DB: pg}
	producer := &queue.Producer{}

	// --- Redis ---
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	producer.Client = rdb

	// --- OCR processor ---
	proc := ocr.NewProcessor()

	logger.Info("ocr-page-worker started")

	for {
		// block for work
		res, err := rdb.BRPop(ctx, 0, "ocr_page_queue").Result()
		if err != nil {
			logger.Error("queue pop failed", "err", err)
			continue
		}

		var msg queue.PageMessage
		if err := json.Unmarshal([]byte(res[1]), &msg); err != nil {
			logger.Error("bad message", "err", err)
			continue
		}

		logger.Info("ocr page start",
			"job_id", msg.JobID,
			"file_id", msg.FileID,
			"page", msg.PageNum,
		)

		// --- timeout guard per page ---
		pageCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)

		text, err := proc.ProcessFile(pageCtx, msg.ImagePath)
		cancel()

		if err != nil {
			logger.Error("ocr failed",
				"job_id", msg.JobID,
				"page", msg.PageNum,
				"err", err,
			)
			continue
		}

		if err := pageRepo.MarkDone(ctx, msg.FileID, msg.PageNum, text); err != nil {
			logger.Error("mark page failed", "err", err)
			continue
		}

		// barrier check
		allDone, err := pageRepo.AllDone(ctx, msg.FileID)
		if err != nil {
			logger.Error("barrier check failed", "err", err)
			continue
		}

		// if file finished → trigger vector stage
		if allDone {
			logger.Info("all pages complete", "file_id", msg.FileID)

			vectorMsg := queue.Message{
				JobID:     msg.JobID,
				Attempt:   0,
				Timestamp: time.Now(),
			}

			if err := producer.Publish(ctx, "vector_queue", vectorMsg); err != nil {
				logger.Error("vector enqueue failed", "err", err)
			}

			_ = fileRepo.UpdateStatus(ctx, msg.FileID, "OCR_DONE")
		}

		logger.Info("ocr page done",
			"job_id", msg.JobID,
			"page", msg.PageNum,
			"text_len", len(text),
		)
	}
}
