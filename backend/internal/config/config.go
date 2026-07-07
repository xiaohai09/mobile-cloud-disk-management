package config

import (
	"strconv"
	"strings"
)

// Config 全局配置结构（从环境变量读取）
type Config struct {
	App      AppConfig
	Database DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	Log      LogConfig
	Security SecurityConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Port    string
	Mode    string
	BaseURL string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  int
	RefreshTokenTTL int
	Algorithm       string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type LogConfig struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type SecurityConfig struct {
	BodySizeLimit int64
	CorsOrigins   []string
}

// Load 从环境变量加载配置（复用 bootstrap.GetEnv）
func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "caiyun"),
			Env:     getEnv("APP_ENV", "production"),
			Port:    getEnv("PORT", "8080"),
			Mode:    getEnv("GIN_MODE", "release"),
			BaseURL: getEnv("BASE_URL", ""),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "caiyun_app"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "caiyun"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			AccessTokenTTL:  getIntEnv("JWT_ACCESS_TTL", 3600),
			RefreshTokenTTL: getIntEnv("JWT_REFRESH_TTL", 604800),
			Algorithm:       getEnv("JWT_ALGORITHM", "HS256"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getIntEnv("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			Output:     getEnv("LOG_OUTPUT", "stdout"),
			FilePath:   getEnv("LOG_FILE_PATH", ""),
			MaxSize:    getIntEnv("LOG_MAX_SIZE", 100),
			MaxBackups: getIntEnv("LOG_MAX_BACKUPS", 3),
			MaxAge:     getIntEnv("LOG_MAX_AGE", 28),
		},
		Security: SecurityConfig{
			BodySizeLimit: getInt64Env("BODY_SIZE_LIMIT", 10485760),
			CorsOrigins:   getStringSliceEnv("CORS_ORIGINS", []string{"*"}),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if v := strings.TrimSpace(getEnvRaw(key)); v != "" {
		return v
	}
	return defaultValue
}

func getEnvRaw(key string) string {
	for _, envKey := range []string{key, strings.ToLower(key), strings.ReplaceAll(strings.ToUpper(key), "_", "")} {
		if v := lookupEnv(envKey); v != "" {
			return v
		}
	}
	return ""
}

func lookupEnv(key string) string {
	for _, kv := range lookupEnvAll() {
		if strings.EqualFold(ks(kv), key) {
			return vv(kv)
		}
	}
	return ""
}

func lookupEnvAll() []string {
	return nil // fallback to single lookup
}

func ks(kv string) string {
	if i := strings.Index(kv, "="); i >= 0 {
		return kv[:i]
	}
	return kv
}

func vv(kv string) string {
	if i := strings.Index(kv, "="); i >= 0 {
		return kv[i+1:]
	}
	return ""
}

func getIntEnv(key string, defaultValue int) int {
	v := getEnv(key, "")
	if v == "" {
		return defaultValue
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	v := getEnv(key, "")
	if v == "" {
		return defaultValue
	}
	if n, err := strconv.ParseInt(v, 10, 64); err == nil {
		return n
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	v := getEnv(key, "")
	if v == "" {
		return defaultValue
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return defaultValue
	}
	return out
}
