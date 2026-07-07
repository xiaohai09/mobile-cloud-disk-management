package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("DB_HOST")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.App.Name != "caiyun" {
		t.Errorf("App.Name = %q, want %q", cfg.App.Name, "caiyun")
	}
	if cfg.App.Port != "8080" {
		t.Errorf("App.Port = %q, want %q", cfg.App.Port, "8080")
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host = %q, want %q", cfg.Database.Host, "localhost")
	}
	if cfg.Redis.DB != 0 {
		t.Errorf("Redis.DB = %d, want 0", cfg.Redis.DB)
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	_ = os.Setenv("APP_NAME", "test-app")
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("DB_HOST", "db-host")
	_ = os.Setenv("REDIS_DB", "2")
	_ = os.Setenv("CORS_ORIGINS", "https://a.com,https://b.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.App.Name != "test-app" {
		t.Errorf("App.Name = %q, want %q", cfg.App.Name, "test-app")
	}
	if cfg.App.Port != "9090" {
		t.Errorf("App.Port = %q, want %q", cfg.App.Port, "9090")
	}
	if cfg.Database.Host != "db-host" {
		t.Errorf("Database.Host = %q, want %q", cfg.Database.Host, "db-host")
	}
	if cfg.Redis.DB != 2 {
		t.Errorf("Redis.DB = %d, want 2", cfg.Redis.DB)
	}
	if len(cfg.Security.CorsOrigins) != 2 {
		t.Errorf("CorsOrigins len = %d, want 2", len(cfg.Security.CorsOrigins))
	}

	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("DB_HOST")
	_ = os.Unsetenv("REDIS_DB")
	_ = os.Unsetenv("CORS_ORIGINS")
}

func TestGetStringSliceEnv(t *testing.T) {
	_ = os.Setenv("CORS_ORIGINS", "https://a.com,https://b.com")
	got := getStringSliceEnv("CORS_ORIGINS", []string{"*"})
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
	_ = os.Unsetenv("CORS_ORIGINS")
}
