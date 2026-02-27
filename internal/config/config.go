// Private application code

package config

import "os"

type Config struct {
	PostgresDSN string
	RedisAddr   string
	MinioURL    string
	MinioKey    string
	MinioSecret string
	Bucket      string
}

func Load() Config {
	return Config{
		PostgresDSN: get("POSTGRES_DSN", "postgres://postgres:postgres@postgres:5432/docflow?sslmode=disable"),
		RedisAddr:   get("REDIS_ADDR", "redis:6379"),
		MinioURL:    get("MINIO_URL", "minio:9000"),
		MinioKey:    get("MINIO_ROOT_USER", "minioadmin"),
		MinioSecret: get("MINIO_ROOT_PASSWORD", "minioadmin"),
		Bucket:      get("MINIO_BUCKET", "uploads"),
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
