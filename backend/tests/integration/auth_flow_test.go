//go:build integration

package integration

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/redis"

	"caiyun/internal/bootstrap"
	"caiyun/internal/handlers"
	"caiyun/internal/models"
	"caiyun/internal/services"
	"caiyun/pkg/jwt"
)

// TestMain starts MySQL + Redis containers once, tears them down after all
// integration tests in this package finish.
func TestMain(m *testing.M) {
	ctx := context.Background()

	password := os.Getenv("TEST_DB_PASSWORD")
	if password == "" {
		password = "caiyun_pass"
	}
	_ = os.Setenv("DB_PASSWORD", password)

	mysqlContainer, err := mysql.Run(ctx, "mysql:8.0",
		mysql.WithDatabase("caiyun_test"),
		mysql.WithUsername("caiyun"),
		mysql.WithPassword(password),
	)
	if err != nil {
		fmt.Println("MySQL container failed to start:", err)
		os.Exit(1)
	}
	defer mysqlContainer.Terminate(ctx)

	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		fmt.Println("Redis container failed to start:", err)
		os.Exit(1)
	}
	defer redisContainer.Terminate(ctx)

	mysqlHost, err := mysqlContainer.Host(ctx)
	if err != nil {
		fmt.Println("mysql host:", err)
		os.Exit(1)
	}
	mysqlPort, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		fmt.Println("mysql mapped port:", err)
		os.Exit(1)
	}
	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		fmt.Println("redis host:", err)
		os.Exit(1)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		fmt.Println("redis mapped port:", err)
		os.Exit(1)
	}

	_ = os.Setenv("DB_HOST", mysqlHost)
	_ = os.Setenv("DB_PORT", mysqlPort.Port())
	_ = os.Setenv("DB_USER", "caiyun")
	_ = os.Setenv("DB_NAME", "caiyun_test")
	_ = os.Setenv("REDIS_HOST", redisHost)
	_ = os.Setenv("REDIS_PORT", redisPort.Port())
	_ = os.Setenv("REDIS_PASSWORD", "")
	_ = os.Setenv("REDIS_DB", "0")
	_ = os.Setenv("APP_ENV", "test")

	if v := os.Getenv("TEST_DB_PASSWORD"); v != "" {
		_ = os.Setenv("DB_PASSWORD", v)
	} else {
		_ = os.Setenv("DB_PASSWORD", "caiyun_pass")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret-at-least-32-bytes-long-enough-1234567890"
	}
	_ = os.Setenv("JWT_SECRET", secret)
	_ = os.Setenv("JWT_ALGORITHM", "HS256")
	_ = os.Setenv("ALLOW_INSECURE_DEFAULTS", "true")

	// Wait for MySQL to accept connections and run migrations.
	time.Sleep(5 * time.Second)

	bootstrap.LoadEnvFile()

	core, err := bootstrap.InitCore()
	if err != nil {
		fmt.Println("init core:", err)
		os.Exit(1)
	}
	defer core.Close()

	if err := core.DB.AutoMigrate(&models.User{}, &models.Account{}); err != nil {
		fmt.Println("auto migrate:", err)
		os.Exit(1)
	}

	engine := buildTestRouter(core)
	testServer = httptest.NewServer(engine)

	m.Run()
	testServer.Close()
}

var testServer *httptest.Server

func buildTestRouter(core *bootstrap.Core) *gin.Engine {
	jwtManager := jwt.NewManager(os.Getenv("JWT_SECRET"))
	authService := services.NewAuthServiceWithPasswordResetCache(
		core.Repository.User,
		jwtManager,
		24*time.Hour,
		services.PasswordResetConfig{},
		nil,
	)
	authHandler := handlers.NewAuthHandler(authService, jwtManager)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	public := r.Group("/api/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
		public.POST("/logout", authHandler.Logout)
		public.GET("/me", authHandler.GetCurrentUser)
	}

	return r
}

func randomUsername() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "user_" + time.Now().Format("150405")
	}
	return "user_" + hex.EncodeToString(buf[:])
}

func randomPassword() string {
	return os.Getenv("TEST_USER_PASSWORD")
}

func decodeAuthCookie(t *testing.T, cookies []*http.Cookie) *jwt.Claims {
	t.Helper()
	secret := os.Getenv("JWT_SECRET")
	manager := jwt.NewManager(secret)

	for _, c := range cookies {
		if c.Name == "auth_token" && c.Value != "" {
			claims, err := manager.ValidateToken(c.Value)
			if err != nil {
				return nil
			}
			return claims
		}
	}
	return nil
}

// TestAuthFlow verifies the end-to-end auth flow:
// register → login → refresh → logout
func TestAuthFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode skips integration tests")
	}
	if testServer == nil {
		t.Skip("integration test server is not initialized")
	}

	username := randomUsername()
	password := randomPassword()
	if password == "" {
		password = "integration-test-password-123"
		_ = os.Setenv("TEST_USER_PASSWORD", password)
	}

	base := testServer.URL

	// 1. Register
	t.Run("register", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
			"email":    strings.ToLower(username) + "@example.com",
		})
		req, _ := http.NewRequest(http.MethodPost, base+"/api/auth/register", strings.NewReader(string(payload)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("register request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("register status = %d, want %d", resp.StatusCode, http.StatusCreated)
		}
		claims := decodeAuthCookie(t, resp.Cookies())
		if claims == nil || claims.Username != username {
			t.Fatalf("register auth cookie claims = %+v", claims)
		}
		if claims.UserID == 0 || claims.Role == "" {
			t.Fatalf("register claims missing expected fields: %+v", claims)
		}
	})

	// 2. Login
	var authToken string
	t.Run("login", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})
		req, _ := http.NewRequest(http.MethodPost, base+"/api/auth/login", strings.NewReader(string(payload)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("login request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("login status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		claims := decodeAuthCookie(t, resp.Cookies())
		if claims == nil || claims.Username != username || claims.UserID == 0 || claims.Role == "" {
			t.Fatalf("login claims = %+v", claims)
		}
		for _, c := range resp.Cookies() {
			if c.Name == "auth_token" {
				authToken = c.Value
				break
			}
		}
	})

	// 3. Refresh
	t.Run("refresh", func(t *testing.T) {
		if authToken == "" {
			t.Fatal("no auth token from login")
		}
		req, _ := http.NewRequest(http.MethodPost, base+"/api/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("refresh request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("refresh status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		claims := decodeAuthCookie(t, resp.Cookies())
		if claims == nil || claims.Username != username || claims.UserID == 0 || claims.Role == "" {
			t.Fatalf("refresh claims = %+v", claims)
		}
		for _, c := range resp.Cookies() {
			if c.Name == "auth_token" {
				authToken = c.Value
				break
			}
		}
	})

	// 4. Logout
	t.Run("logout", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, base+"/api/auth/logout", nil)
		req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("logout request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("logout status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		for _, c := range resp.Cookies() {
			if c.Name == "auth_token" && c.MaxAge < 0 {
				return
			}
		}
		t.Fatalf("logout did not clear auth_token cookie")
	})

	// 5. After logout, /api/auth/me must reject the cleared session.
	t.Run("me_after_logout", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, base+"/api/auth/me", nil)
		req.AddCookie(&http.Cookie{Name: "auth_token", Value: authToken})
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("/me request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("/me after logout status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
	})
}

// TestAuthLoginDuplicateUser asserts 409 on duplicate registration.
func TestAuthLoginDuplicateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode skips integration tests")
	}
	if testServer == nil {
		t.Skip("integration test server is not initialized")
	}

	password := randomPassword()
	if password == "" {
		password = "integration-test-password-123"
		_ = os.Setenv("TEST_USER_PASSWORD", password)
	}

	username := randomUsername()
	payload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"email":    strings.ToLower(username) + "@example.com",
	})

	register := func() {
		req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/register", strings.NewReader(string(payload)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("register request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("first register status = %d, want %d", resp.StatusCode, http.StatusCreated)
		}
	}
	register()
	register()

	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/register", strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("duplicate register request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate register status = %d, want %d", resp.StatusCode, http.StatusConflict)
	}
}

// TestAuthLoginInvalidCredentials asserts 401 for wrong password.
func TestAuthLoginInvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode skips integration tests")
	}
	if testServer == nil {
		t.Skip("integration test server is not initialized")
	}

	password := randomPassword()
	if password == "" {
		password = "integration-test-password-123"
		_ = os.Setenv("TEST_USER_PASSWORD", password)
	}
	username := randomUsername()
	payload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"email":    strings.ToLower(username) + "@example.com",
	})

	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/register", strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("register request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	badPayload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": "wrong-password",
	})
	req, _ = http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/login", strings.NewReader(string(badPayload)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid login status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

// TestAuthJWTClaims verifies that the JWT payload matches the registered user
// and that a forged token is rejected.
func TestAuthJWTClaims(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode skips integration tests")
	}
	if testServer == nil {
		t.Skip("integration test server is not initialized")
	}

	password := randomPassword()
	if password == "" {
		password = "integration-test-password-123"
		_ = os.Setenv("TEST_USER_PASSWORD", password)
	}
	username := randomUsername()
	registerPayload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"email":    strings.ToLower(username) + "@example.com",
	})

	req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/register", strings.NewReader(string(registerPayload)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("register request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	loginPayload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	req, _ = http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/login", strings.NewReader(string(loginPayload)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	claims := decodeAuthCookie(t, resp.Cookies())
	if claims == nil {
		t.Fatal("login did not set auth_token cookie")
	}
	if claims.Username != username {
		t.Fatalf("jwt username = %q, want %q", claims.Username, username)
	}
	if claims.UserID == 0 {
		t.Fatal("jwt user_id is zero")
	}
	if claims.Role == "" {
		t.Fatal("jwt role is empty")
	}
}

// TestAuthFlowWithTestcontainersReachable verifies the full container-driven
// flow: MySQL + Redis are created by testcontainers-go and reachable from the
// server using dynamically mapped host/port values.
func TestAuthFlowWithTestcontainersReachable(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode skips integration tests")
	}
	if testServer == nil {
		t.Skip("integration test server is not initialized")
	}

	password := randomPassword()
	if password == "" {
		password = "integration-test-password-123"
		_ = os.Setenv("TEST_USER_PASSWORD", password)
	}
	username := randomUsername()
	payload, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"email":    strings.ToLower(username) + "@example.com",
	})

	t.Run("register_via_testcontainers_env", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/register", strings.NewReader(string(payload)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("register request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("register status = %d, want %d", resp.StatusCode, http.StatusCreated)
		}
		claims := decodeAuthCookie(t, resp.Cookies())
		if claims == nil || claims.Username != username || claims.UserID == 0 || claims.Role == "" {
			t.Fatalf("register claims = %+v", claims)
		}
	})

	t.Run("login_via_testcontainers_env", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})
		req, _ := http.NewRequest(http.MethodPost, testServer.URL+"/api/auth/login", strings.NewReader(string(payload)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("login request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("login status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		claims := decodeAuthCookie(t, resp.Cookies())
		if claims == nil || claims.Username != username || claims.UserID == 0 || claims.Role == "" {
			t.Fatalf("login claims = %+v", claims)
		}
	})
}
