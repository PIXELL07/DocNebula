package ocr

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// Processor simulates OCR work.
type Processor struct {
	Logger *slog.Logger
}

// ProcessFile performs OCR on a single file.
func (p *Processor) ProcessFile(ctx context.Context, filePath string) (string, error) {
	p.Logger.Info("ocr started", "file", filePath)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(700 * time.Millisecond):
	}

	// simulated extracted text
	text := fmt.Sprintf("extracted text from %s", filePath)

	p.Logger.Info("ocr completed", "file", filePath)
	return text, nil
}

// NewProcessor creates OCR processor.
func NewProcessor() *Processor {
	return &Processor{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}
