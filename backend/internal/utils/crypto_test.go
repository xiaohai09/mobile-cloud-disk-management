package utils

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func generateTestKey() []byte {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	return key
}

func TestCrypto_EncryptDecrypt(t *testing.T) {
	crypto, err := NewCrypto(generateTestKey())
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	plaintext := "test secret value"
	ciphertext, err := crypto.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if ciphertext == plaintext {
		t.Error("Encrypted value should differ from plaintext")
	}

	if !hasPrefix(ciphertext, encryptionPrefix) {
		t.Errorf("Encrypted value should have prefix %q, got %q", encryptionPrefix, ciphertext)
	}

	decrypted, err := crypto.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted value mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestCrypto_DecryptPlaintext(t *testing.T) {
	crypto, err := NewCrypto(generateTestKey())
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	plaintext := "plaintext value without prefix"
	decrypted, err := crypto.Decrypt(plaintext)
	if err != nil {
		t.Fatalf("Decrypt should not fail on plaintext: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Plaintext should be returned unchanged: got %q, want %q", decrypted, plaintext)
	}
}

func TestCrypto_EmptyString(t *testing.T) {
	crypto, err := NewCrypto(generateTestKey())
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	enc, err := crypto.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt empty string failed: %v", err)
	}
	if enc != "" {
		t.Errorf("Empty string should return empty string, got %q", enc)
	}

	dec, err := crypto.Decrypt("")
	if err != nil {
		t.Fatalf("Decrypt empty string failed: %v", err)
	}
	if dec != "" {
		t.Errorf("Empty string should return empty string, got %q", dec)
	}
}

func TestCrypto_WrongKey(t *testing.T) {
	crypto1, err := NewCrypto(generateTestKey())
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	crypto2, err := NewCrypto(generateTestKey())
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	ciphertext, err := crypto1.Encrypt("secret")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = crypto2.Decrypt(ciphertext)
	if err == nil {
		t.Error("Decrypt with wrong key should fail")
	}
}

func TestCrypto_Base64RoundTrip(t *testing.T) {
	crypto, err := NewCrypto(generateTestKey())
	if err != nil {
		t.Fatalf("NewCrypto failed: %v", err)
	}

	plaintext := "Basic dGVzdDp0ZXN0OnRlc3Q=" // sample auth-like value
	ciphertext, err := crypto.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Verify it's valid base64 after prefix
	encoded := ciphertext[len(encryptionPrefix):]
	_, err = base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("Encrypted value should be valid base64: %v", err)
	}

	decrypted, err := crypto.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Round-trip failed: got %q, want %q", decrypted, plaintext)
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
