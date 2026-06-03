package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encrypt encrypts plaintext using AES-256-GCM with the given hex-encoded key.
// Returns hex-encoded ciphertext (nonce + encrypted data).
// If key is empty, returns plaintext as-is (graceful fallback for development).
func Encrypt(plaintext string, hexKey string) (string, error) {
	if hexKey == "" {
		return plaintext, nil
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("crypto: invalid hex key: %w", err)
	}
	if len(key) != 32 {
		return "", errors.New("crypto: key must be 32 bytes (64 hex chars)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("crypto: new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("crypto: new gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("crypto: generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts hex-encoded ciphertext using AES-256-GCM with the given hex-encoded key.
// If key is empty, returns ciphertext as-is (graceful fallback for development).
func Decrypt(hexCiphertext string, hexKey string) (string, error) {
	if hexKey == "" {
		return hexCiphertext, nil
	}

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("crypto: invalid hex key: %w", err)
	}
	if len(key) != 32 {
		return "", errors.New("crypto: key must be 32 bytes (64 hex chars)")
	}

	ciphertext, err := hex.DecodeString(hexCiphertext)
	if err != nil {
		// ถ้า decode hex ไม่ได้ แสดงว่า token เป็น plain text เดิม (migration period)
		// คืน plain text กลับไปเพื่อ backward compatibility
		return hexCiphertext, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("crypto: new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("crypto: new gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		// ข้อมูลสั้นเกินไป → น่าจะเป็น plain text เดิม (migration period)
		return hexCiphertext, nil
	}

	nonce, ciphertextBytes := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		// Decrypt ไม่สำเร็จ → น่าจะเป็น plain text เดิม (migration period)
		return hexCiphertext, nil
	}

	return string(plaintext), nil
}
