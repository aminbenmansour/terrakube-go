package config

import (
	"log/slog"
	"testing"
	"time"
)

var allKeys = []string{
	"TERRAKUBE_ENV",
	"TERRAKUBE_PUBLIC_URL",
	"TERRAKUBE_LOG_LEVEL",
	"TERRAKUBE_DATABASE_URL",
	"DATABASE_URL",
	"TERRAKUBE_AUTO_MIGRATE",
	"TERRAKUBE_REDIS_URL",
	"REDIS_URL",
	"TERRAKUBE_AUTH_TOKEN",
	"TERRAKUBE_DEX_ENABLED",
	"TERRAKUBE_DEX_AUDIENCE_VALIDATION",
	"TERRAKUBE_DEX_EXTRA_AUDIENCES",
	"TERRAKUBE_INTERNAL_TOKEN_SECRET",
	"IO_TERRAKUBE_TOKEN_INTERNAL",
	"TERRAKUBE_STATE_MAX_SIZE_BYTES",
	"IO_TERRAKUBE_STATE_MAX_SIZE_BYTES",
	"TERRAKUBE_SHUTDOWN_TIMEOUT",
	"TERRAKUBE_HTTP_ADDR",
	"TERRAKUBE_EXECUTOR_HTTP_ADDR",
	"TERRAKUBE_REGISTRY_HTTP_ADDR",
	"TERRAKUBE_STORAGE_TYPE",
	"TERRAKUBE_STORAGE_LOCAL_PATH",
	"TERRAKUBE_STORAGE_S3_BUCKET",
	"TERRAKUBE_STORAGE_S3_REGION",
	"TERRAKUBE_STORAGE_S3_ENDPOINT",
	"TERRAKUBE_STORAGE_S3_ACCESS_KEY_ID",
	"TERRAKUBE_STORAGE_S3_SECRET_ACCESS_KEY",
	"TERRAKUBE_STORAGE_S3_USE_PATH_STYLE",
	"TERRAKUBE_ENCRYPTION_KEYS",
	"TERRAKUBE_ENCRYPTION_ACTIVE_KEY",
}

func clearEnv(t *testing.T) {
	t.Helper()
	for _, key := range allKeys {
		t.Setenv(key, "")
	}
}

func TestLoadDefaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Env != "dev" {
		t.Errorf("Env = %q, want dev", cfg.Env)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Errorf("LogLevel = %v, want info", cfg.LogLevel)
	}
	if cfg.Storage.Type != StorageTypeLocal {
		t.Errorf("Storage.Type = %q, want local", cfg.Storage.Type)
	}
	if cfg.AutoMigrate != true {
		t.Errorf("AutoMigrate = %v, want true", cfg.AutoMigrate)
	}
	if cfg.ShutdownTimeout != 15*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 15s", cfg.ShutdownTimeout)
	}
	if cfg.PublicURL != "http://localhost:8080" {
		t.Errorf("PublicURL = %q, want http://localhost:8080", cfg.PublicURL)
	}
	if cfg.StateMaxSizeBytes != 104857600 {
		t.Errorf("StateMaxSizeBytes = %d, want 104857600", cfg.StateMaxSizeBytes)
	}
	if cfg.ExecutorJobTimeout != 2*time.Hour {
		t.Errorf("ExecutorJobTimeout = %v, want 2h", cfg.ExecutorJobTimeout)
	}
	if cfg.DexEnabled || cfg.DexAudienceValidation || len(cfg.DexExtraAudiences) != 0 {
		t.Errorf("Dex settings = %v/%v/%v, want disabled with no audience checks",
			cfg.DexEnabled, cfg.DexAudienceValidation, cfg.DexExtraAudiences)
	}
	if cfg.VCSAllowLocal {
		t.Error("VCSAllowLocal = true by default, want false")
	}
	if cfg.PATSecret != "" || cfg.AdminGroup != "TERRAKUBE_ADMIN" {
		t.Errorf("PAT settings = %q/%q, want empty/TERRAKUBE_ADMIN", cfg.PATSecret, cfg.AdminGroup)
	}
	if cfg.TerraformLogin || cfg.TerraformLoginDefaultDays != 30 || cfg.TerraformLoginMaxDays != 365 {
		t.Errorf("login settings = %v/%d/%d, want false/30/365",
			cfg.TerraformLogin, cfg.TerraformLoginDefaultDays, cfg.TerraformLoginMaxDays)
	}
}

func TestLoadOverrides(t *testing.T) {
	clearEnv(t)
	t.Setenv("TERRAKUBE_ENV", "prod")
	t.Setenv("TERRAKUBE_ENCRYPTION_KEYS", "dev:configured")
	t.Setenv("TERRAKUBE_LOG_LEVEL", "debug")
	t.Setenv("TERRAKUBE_AUTO_MIGRATE", "false")
	t.Setenv("TERRAKUBE_STORAGE_TYPE", "s3")
	t.Setenv("TERRAKUBE_STORAGE_S3_BUCKET", "state")
	t.Setenv("TERRAKUBE_PUBLIC_URL", "https://terrakube.example.com/")
	t.Setenv("TERRAKUBE_EXECUTOR_JOB_TIMEOUT", "45m")
	t.Setenv("TERRAKUBE_DEX_ENABLED", "true")
	t.Setenv("TERRAKUBE_DEX_AUDIENCE_VALIDATION", "true")
	t.Setenv("TERRAKUBE_DEX_EXTRA_AUDIENCES", "terrakube, second-client")
	t.Setenv("TERRAKUBE_PAT_SECRET", "secret")
	t.Setenv("TERRAKUBE_ADMIN_GROUP", "platform-admins")
	t.Setenv("TERRAKUBE_TERRAFORM_LOGIN_DEFAULT_DAYS", "14")
	t.Setenv("TERRAKUBE_TERRAFORM_LOGIN_MAX_DAYS", "90")
	t.Setenv("TERRAKUBE_VCS_ALLOW_LOCAL", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.LogLevel != slog.LevelDebug {
		t.Errorf("LogLevel = %v, want debug", cfg.LogLevel)
	}
	if cfg.AutoMigrate {
		t.Error("AutoMigrate = true, want false")
	}
	if cfg.Storage.Type != StorageTypeS3 {
		t.Errorf("Storage.Type = %q, want s3", cfg.Storage.Type)
	}
	if cfg.PublicURL != "https://terrakube.example.com" {
		t.Errorf("PublicURL = %q, want trailing slash trimmed", cfg.PublicURL)
	}
	if cfg.ExecutorJobTimeout != 45*time.Minute {
		t.Errorf("ExecutorJobTimeout = %v, want 45m", cfg.ExecutorJobTimeout)
	}
	if !cfg.DexEnabled || !cfg.DexAudienceValidation {
		t.Errorf("Dex enabled/audience = %v/%v, want true/true", cfg.DexEnabled, cfg.DexAudienceValidation)
	}
	if len(cfg.DexExtraAudiences) != 2 || cfg.DexExtraAudiences[1] != "second-client" {
		t.Errorf("DexExtraAudiences = %v, want [terrakube second-client]", cfg.DexExtraAudiences)
	}
	if cfg.PATSecret != "secret" || cfg.AdminGroup != "platform-admins" {
		t.Errorf("PATSecret/AdminGroup = %q/%q, want secret/platform-admins", cfg.PATSecret, cfg.AdminGroup)
	}
	if cfg.TerraformLoginDefaultDays != 14 || cfg.TerraformLoginMaxDays != 90 {
		t.Errorf("login days = %d/%d, want 14/90", cfg.TerraformLoginDefaultDays, cfg.TerraformLoginMaxDays)
	}
	if !cfg.VCSAllowLocal {
		t.Error("VCSAllowLocal = false, want true")
	}
}

func TestProductionRequiresEncryptionKeys(t *testing.T) {
	t.Setenv("TERRAKUBE_ENV", "production")
	t.Setenv("TERRAKUBE_ENCRYPTION_KEYS", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() with production and no keys = nil error, want error")
	}
}

func TestProductionAliasRequiresEncryptionKeys(t *testing.T) {
	t.Setenv("TERRAKUBE_ENV", "prod")
	t.Setenv("TERRAKUBE_ENCRYPTION_KEYS", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() with prod and no keys = nil error, want error")
	}
}

func TestDevelopmentAllowsMissingEncryptionKeys(t *testing.T) {
	t.Setenv("TERRAKUBE_ENV", "dev")
	t.Setenv("TERRAKUBE_ENCRYPTION_KEYS", "")
	if _, err := Load(); err != nil {
		t.Fatalf("Load() = %v, want nil", err)
	}
}

func TestLoadStateSettings(t *testing.T) {
	clearEnv(t)
	t.Setenv("TERRAKUBE_INTERNAL_TOKEN_SECRET", "secret")
	t.Setenv("TERRAKUBE_STATE_MAX_SIZE_BYTES", "2048")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.InternalTokenSecret != "secret" {
		t.Errorf("InternalTokenSecret = %q, want secret", cfg.InternalTokenSecret)
	}
	if cfg.StateMaxSizeBytes != 2048 {
		t.Errorf("StateMaxSizeBytes = %d, want 2048", cfg.StateMaxSizeBytes)
	}
}

func TestLoadRejectsUnknownStorage(t *testing.T) {
	clearEnv(t)
	t.Setenv("TERRAKUBE_STORAGE_TYPE", "gcs")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want unknown storage type error")
	}
}

func TestLoadRejectsS3WithoutBucket(t *testing.T) {
	clearEnv(t)
	t.Setenv("TERRAKUBE_STORAGE_TYPE", "s3")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want missing bucket error")
	}
}

func TestParseLogLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"":        slog.LevelInfo,
		"info":    slog.LevelInfo,
		"WARN":    slog.LevelWarn,
		"warning": slog.LevelWarn,
		"error":   slog.LevelError,
	}
	for input, want := range cases {
		got, err := ParseLogLevel(input)
		if err != nil {
			t.Fatalf("ParseLogLevel(%q) error = %v", input, err)
		}
		if got != want {
			t.Errorf("ParseLogLevel(%q) = %v, want %v", input, got, want)
		}
	}

	if _, err := ParseLogLevel("loud"); err == nil {
		t.Fatal("ParseLogLevel(loud) error = nil, want error")
	}
}

func TestDatabaseURLFallback(t *testing.T) {
	clearEnv(t)
	t.Setenv("DATABASE_URL", "postgres://legacy@db/terrakube")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DatabaseURL != "postgres://legacy@db/terrakube" {
		t.Errorf("DatabaseURL = %q, want legacy fallback", cfg.DatabaseURL)
	}

	t.Setenv("TERRAKUBE_DATABASE_URL", "postgres://new@db/terrakube")
	cfg, err = Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DatabaseURL != "postgres://new@db/terrakube" {
		t.Errorf("DatabaseURL = %q, want TERRAKUBE_DATABASE_URL to win", cfg.DatabaseURL)
	}
}
