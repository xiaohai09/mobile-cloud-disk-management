package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// encryptionPrefix marks encrypted values for backward compatibility with existing plaintext data.
	encryptionPrefix = "enc:"
	// nonceSize is the AES-GCM nonce size (12 bytes).
	nonceSize = 12
)

// Crypto provides AES-256-GCM encryption/decryption for sensitive fields.
type Crypto struct {
	key []byte
}

// NewCrypto creates a new Crypto instance with the given 32-byte key.
func NewCrypto(key []byte) (*Crypto, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes, got %d", len(key))
	}
	return &Crypto{key: key}, nil
}

// NewCryptoFromEnv creates a new Crypto instance using the DB_ENCRYPTION_KEY environment variable.
// The key is derived from DB_ENCRYPTION_KEY using SHA-256 to ensure it's always 32 bytes.
func NewCryptoFromEnv() (*Crypto, error) {
	keyStr := strings.TrimSpace(os.Getenv("DB_ENCRYPTION_KEY"))
	if keyStr == "" {
		return nil, fmt.Errorf("DB_ENCRYPTION_KEY environment variable is required for field encryption")
	}
	// Use SHA-256 to derive a 32-byte key from any length input
	keyHash := sha256.Sum256([]byte(keyStr))
	return NewCrypto(keyHash[:])
}

// Encrypt encrypts plaintext using AES-256-GCM and returns a base64-encoded string prefixed with "enc:".
// Returns the plaintext unchanged if it's empty.
func (c *Crypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return plaintext, nil
	}

	// Generate random nonce
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	// Encrypt: nonce + ciphertext + tag
	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64 encode and prefix
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encryptionPrefix + encoded, nil
}

// Decrypt decrypts a string that was encrypted with Encrypt.
// If the value is empty or doesn't have the encryption prefix, it's returned as-is (backward compatibility).
func (c *Crypto) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return ciphertext, nil
	}

	// Check if the value is encrypted (has the prefix)
	if !strings.HasPrefix(ciphertext, encryptionPrefix) {
		// Plaintext value - return as-is for backward compatibility with existing data
		return ciphertext, nil
	}

	// Remove prefix
	encoded := strings.TrimPrefix(ciphertext, encryptionPrefix)

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode encrypted value: %w", err)
	}

	if len(decoded) < nonceSize {
		return "", fmt.Errorf("encrypted value too short: %d bytes", len(decoded))
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	// Split nonce and ciphertext
	nonce, cipherData := decoded[:nonceSize], decoded[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EncryptField encrypts a single field value, used as a helper for repository methods.
// Returns the encrypted string or original if empty.
func (c *Crypto) EncryptField(value string) (string, error) {
	if value == "" {
		return value, nil
	}
	return c.Encrypt(value)
}

// DecryptField decrypts a single field value, used as a helper for repository methods.
// Returns the decrypted string or original if empty/not encrypted.
func (c *Crypto) DecryptField(value string) (string, error) {
	if value == "" {
		return value, nil
	}
	return c.Decrypt(value)
}

// encryptOrEmpty encrypts a string if non-empty, returns empty string unchanged.
func EncryptOrEmpty(c *Crypto, value string) string {
	if value == "" {
		return value
	}
	encrypted, err := c.Encrypt(value)
	if err != nil {
		return value
	}
	return encrypted
}

// decryptOrEmpty decrypts a string if non-empty, returns empty string unchanged.
func DecryptOrEmpty(c *Crypto, value string) string {
	if value == "" {
		return value
	}
	decrypted, err := c.Decrypt(value)
	if err != nil {
		return value
	}
	return decrypted
}
