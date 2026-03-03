package ocr

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
)

// PDFToImages converts a PDF into page images using pdftoppm.

func PDFToImages(
	ctx context.Context,
	logger *slog.Logger,
	pdfPath string,
	outDir string,
) ([]string, error) {

	prefix := filepath.Join(outDir, "page")

	// Ex output: page-1.png, page-2.png...
	cmd := exec.CommandContext(
		ctx,
		"pdftoppm",
		"-png",
		pdfPath,
		prefix,
	)

	logger.Info("pdf to image start", "pdf", pdfPath)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pdftoppm failed: %w", err)
	}

	// collect generated files
	matches, err := filepath.Glob(prefix + "-*.png")
	if err != nil {
		return nil, err
	}

	logger.Info("pdf split complete", "pages", len(matches))
	return matches, nil
}
