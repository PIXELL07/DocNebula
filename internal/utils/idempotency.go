package utils

import (
	// SHA-256 cryptographic hash function
	// No matter the input size, output is always:
	// 64 hex characters (256 bits)
	// For example: hello , output: 2cf24dba5fb0a30e26e83b2ac5b9e29e...
	// if we hello! completely different hash..

	"crypto/sha256" // better than just a Raw String..
	"encoding/hex"
	"io"
	"net/http"
)

// FromString creates deterministic key from string
func FromString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// FromRequest creates fallback key from request metadata
func FromRequest(r *http.Request) string {
	base := r.Method + ":" + r.URL.Path + ":" + r.RemoteAddr
	hash := sha256.Sum256([]byte(base))
	return hex.EncodeToString(hash[:])
}

// FromReader creates content-based key (for file dedup)
func FromReader(reader io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Prefer client Idempotency-Key header, fallback otherwise
func FromHeaderOrRequest(r *http.Request) string {
	if key := r.Header.Get("Idempotency-Key"); key != "" {
		return key
	}
	return FromRequest(r)
}

// -------------------------------

// When You Should Use SHA-256 (Real Systems)
// Use it when you need:

// ** Idempotency keys **
// Most common in:
// payment APIs
// upload deduplication
// job dedup
// retry safety

// ** Content fingerprinting **
// Example:
// detect duplicate files
// cache keys
// integrity checks

// ** Cache key normalization **
// Large systems often hash keys to keep them uniform

// ----------------------------------------

// When NOT to Use SHA-256 ❌
// encryption (it is NOT encryption)
// reversible data
// password storage without salt/bcrypt
// random tokens
