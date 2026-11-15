// Package config loads the runtime configuration for every Terrakube Go
// binary from the environment.
//
// Terrakube Go aims to be deployment-compatible with the upstream Java
// implementation, so where an upstream environment variable exists we accept
// it as an alias for our TERRAKUBE_* name.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	StorageTypeLocal = "local"
	StorageTypeS3    = "s3"
)

// Config is the configuration shared by the api, executor and registry
// binaries. Individual binaries only read the fields they need.
type Config struct {
	Env                 string
	PublicURL           string
	LogLevel            slog.Level
	DatabaseURL         string
	AutoMigrate         bool
	RedisURL            string
	AuthToken           string
	InternalTokenSecret string
	StateMaxSizeBytes   int64
	ShutdownTimeout     time.Duration

	HTTPAddr         string
	ExecutorHTTPAddr string
	RegistryHTTPAddr string

	Storage Storage
}

// Storage configures where state, run logs and uploaded configuration live.
type Storage struct {
	Type      string
	LocalPath string
	S3        S3
}

// S3 holds the settings for an S3-compatible object store (AWS, MinIO, ...).
type S3 struct {
	Bucket          string
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UsePathStyle    bool
}

// Load reads configuration from the environment, applying development-friendly
// defaults so that `go run ./cmd/api` works against the compose stack.
func Load() (Config, error) {
	level, err := ParseLogLevel(env("TERRAKUBE_LOG_LEVEL", "info"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Env:                 env("TERRAKUBE_ENV", "dev"),
		PublicURL:           strings.TrimRight(env("TERRAKUBE_PUBLIC_URL", "http://localhost:8080"), "/"),
		LogLevel:            level,
		DatabaseURL:         envFirst([]string{"TERRAKUBE_DATABASE_URL", "DATABASE_URL"}, "postgres://terrakube@localhost:5432/terrakube?sslmode=disable"),
		AutoMigrate:         envBool("TERRAKUBE_AUTO_MIGRATE", true),
		RedisURL:            envFirst([]string{"TERRAKUBE_REDIS_URL", "REDIS_URL"}, "redis://localhost:6379/0"),
		AuthToken:           env("TERRAKUBE_AUTH_TOKEN", ""),
		InternalTokenSecret: envFirst([]string{"TERRAKUBE_INTERNAL_TOKEN_SECRET", "IO_TERRAKUBE_TOKEN_INTERNAL"}, ""),
		StateMaxSizeBytes:   envInt64([]string{"TERRAKUBE_STATE_MAX_SIZE_BYTES", "IO_TERRAKUBE_STATE_MAX_SIZE_BYTES"}, 104857600),
		ShutdownTimeout:     envDuration("TERRAKUBE_SHUTDOWN_TIMEOUT", 15*time.Second),
		HTTPAddr:            env("TERRAKUBE_HTTP_ADDR", ":8080"),
		ExecutorHTTPAddr:    env("TERRAKUBE_EXECUTOR_HTTP_ADDR", ":8081"),
		RegistryHTTPAddr:    env("TERRAKUBE_REGISTRY_HTTP_ADDR", ":8082"),
		Storage: Storage{
			Type:      strings.ToLower(env("TERRAKUBE_STORAGE_TYPE", StorageTypeLocal)),
			LocalPath: env("TERRAKUBE_STORAGE_LOCAL_PATH", ".data"),
			S3: S3{
				Bucket:          env("TERRAKUBE_STORAGE_S3_BUCKET", ""),
				Region:          env("TERRAKUBE_STORAGE_S3_REGION", "us-east-1"),
				Endpoint:        env("TERRAKUBE_STORAGE_S3_ENDPOINT", ""),
				AccessKeyID:     env("TERRAKUBE_STORAGE_S3_ACCESS_KEY_ID", ""),
				SecretAccessKey: env("TERRAKUBE_STORAGE_S3_SECRET_ACCESS_KEY", ""),
				UsePathStyle:    envBool("TERRAKUBE_STORAGE_S3_USE_PATH_STYLE", false),
			},
		},
	}

	switch cfg.Storage.Type {
	case StorageTypeLocal:
		if cfg.Storage.LocalPath == "" {
			return Config{}, fmt.Errorf("config: TERRAKUBE_STORAGE_LOCAL_PATH must not be empty")
		}
	case StorageTypeS3:
		if cfg.Storage.S3.Bucket == "" {
			return Config{}, fmt.Errorf("config: TERRAKUBE_STORAGE_S3_BUCKET is required when storage type is %q", StorageTypeS3)
		}
	default:
		return Config{}, fmt.Errorf("config: unknown storage type %q (want %q or %q)", cfg.Storage.Type, StorageTypeLocal, StorageTypeS3)
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("config: database URL must not be empty")
	}
	if cfg.RedisURL == "" {
		return Config{}, fmt.Errorf("config: redis URL must not be empty")
	}

	return cfg, nil
}

// ParseLogLevel maps a textual level onto slog.Level.
func ParseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "", "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("config: invalid log level %q (want debug, info, warn or error)", value)
	}
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func envFirst(keys []string, fallback string) string {
	for _, key := range keys {
		if v, ok := os.LookupEnv(key); ok && v != "" {
			return v
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt64(keys []string, fallback int64) int64 {
	for _, key := range keys {
		v, ok := os.LookupEnv(key)
		if !ok || v == "" {
			continue
		}
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}
