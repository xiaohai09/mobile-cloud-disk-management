package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

// RefreshTokenPair 持有 access/refresh token 及其元信息。
type RefreshTokenPair struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  int64
	RefreshTokenExpiresAt int64
}

// HashedRefreshToken 将原始 refresh token 哈希后入库，避免泄露。
type HashedRefreshToken struct {
	Hash     string
	UserID   uint
	ExpiresAt time.Time
	Revoked  bool
}

// hashRefreshToken 计算 refresh token 的 SHA256 哈希（base64url）。
func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// GenerateRefreshTokenPair 生成双 token 对（access + refresh）。
// access token 由外部传入已生成的字符串；refresh token 为高熵随机串。
func GenerateRefreshTokenPair(
	accessToken string,
	accessTTL, refreshTTL time.Duration,
) (*RefreshTokenPair, error) {
	now := time.Now()

	// refresh token: 32 bytes random, base64url encoded
	rbi := make([]byte, 32)
	if _, err := rand.Read(rbi); err != nil {
		return nil, err
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(rbi)

	return &RefreshTokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  now.Add(accessTTL).Unix(),
		RefreshTokenExpiresAt: now.Add(refreshTTL).Unix(),
	}, nil
}

// VerifyRefreshTokenHash 验证原始 refresh token 是否与哈希匹配。
func VerifyRefreshTokenHash(raw, expectedHash string) bool {
	return hashRefreshToken(raw) == expectedHash
}
