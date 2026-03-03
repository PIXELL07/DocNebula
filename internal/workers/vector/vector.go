package vector

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"os"
	"time"
)

// Vectorizer converts text into a fake embedding.
// In real systems its known as embedding model.
type Vectorizer struct {
	Logger *slog.Logger
}

func NewVectorizer() *Vectorizer {
	return &Vectorizer{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// Embed generates a deterministic fake vector from text.
func (v *Vectorizer) Embed(ctx context.Context, text string) (string, error) {
	v.Logger.Info("vectorizing text")

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(300 * time.Millisecond):
	}

	hash := sha256.Sum256([]byte(text))
	vectorID := hex.EncodeToString(hash[:])

	v.Logger.Info("vectorization complete")
	return vectorID, nil
}
