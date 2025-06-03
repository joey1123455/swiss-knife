package swissknife

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncryptAESGCM(t *testing.T) {
	key := make([]byte, 32) // AES-256 key size
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}

	plaintext := []byte("This is a secret message")

	ciphertext, nonce, err := EncryptAESGCM(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAESGCM failed: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Fatal("ciphertext is empty")
	}

	if len(nonce) == 0 {
		t.Fatal("nonce is empty")
	}

	decrypted, err := DecryptAESGCM(ciphertext, nonce, key)
	if err != nil {
		t.Fatalf("DecryptAESGCM failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decrypted text does not match original plaintext. got %q want %q", decrypted, plaintext)
	}
}

func TestDecryptAESGCM_WrongKey(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("failed to generate random key: %v", err)
	}
	plaintext := []byte("Message")
	ciphertext, nonce, err := EncryptAESGCM(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAESGCM failed: %v", err)
	}

	wrongKey := make([]byte, 32)
	if _, err := rand.Read(wrongKey); err != nil {
		t.Fatalf("failed to generate wrong key: %v", err)
	}

	_, err = DecryptAESGCM(ciphertext, nonce, wrongKey)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key, got none")
	}
}
