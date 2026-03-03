package ocr

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/otiai10/gosseract/v2"
)

// Processor performs OCR using Tesseract.
type Processor struct {
	Logger *slog.Logger
}

// NewProcessor creates OCR processor.
func NewProcessor() *Processor {
	return &Processor{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// ProcessFile extracts text from an image/PDF page.
func (p *Processor) ProcessFile(ctx context.Context, filePath string) (string, error) {
	p.Logger.Info("ocr started", "file", filePath)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImage(filePath); err != nil {
		return "", fmt.Errorf("set image failed: %w", err)
	}

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("tesseract failed: %w", err)
	}

	p.Logger.Info("ocr completed", "file", filePath)
	return text, nil
}
