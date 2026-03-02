package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// ReadyHandler checks critical dependencies (DB + Redis).
type ReadyHandlerDeps struct {
	DB    *sql.DB
	Redis *redis.Client
}

func (d ReadyHandlerDeps) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	type check struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	}

	resp := struct {
		Status string  `json:"status"`
		Checks []check `json:"checks"`
	}{Status: "ok"}

	// DB check
	if err := d.DB.PingContext(ctx); err != nil {
		resp.Status = "degraded"
		resp.Checks = append(resp.Checks, check{"postgres", "fail", err.Error()})
	} else {
		resp.Checks = append(resp.Checks, check{"postgres", "ok", ""})
	}

	// Redis check
	if err := d.Redis.Ping(ctx).Err(); err != nil {
		resp.Status = "degraded"
		resp.Checks = append(resp.Checks, check{"redis", "fail", err.Error()})
	} else {
		resp.Checks = append(resp.Checks, check{"redis", "ok", ""})
	}

	code := http.StatusOK
	if resp.Status != "ok" {
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}
