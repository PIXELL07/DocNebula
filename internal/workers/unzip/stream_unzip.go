package unzip

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// StreamUnzip extracts a zip in a memory-safe streaming way.
// Suitable for large archives.

func StreamUnzip(
	ctx context.Context,
	logger *slog.Logger,
	zipPath string,
	destDir string,
) ([]string, error) {

	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open zip failed: %w", err)
	}
	defer zr.Close()

	var extracted []string

	for _, f := range zr.File {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		cleanName := filepath.Clean(f.Name)
		outPath := filepath.Join(destDir, cleanName)

		// ** zip-slip protection **
		if !filepath.HasPrefix(outPath, filepath.Clean(destDir)) {
			logger.Warn("zip slip detected, skipping", "file", f.Name)
			continue
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return nil, err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return nil, err
		}

		src, err := f.Open()
		if err != nil {
			return nil, err
		}

		dst, err := os.Create(outPath)
		if err != nil {
			src.Close()
			return nil, err
		}

		if _, err := io.Copy(dst, src); err != nil {
			dst.Close()
			src.Close()
			return nil, err
		}

		dst.Close()
		src.Close()

		extracted = append(extracted, outPath)
		logger.Info("file extracted", "path", outPath)
	}

	return extracted, nil
}
