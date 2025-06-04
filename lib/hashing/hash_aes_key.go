package swissknife

import "crypto/sha256"

func HashToAESKey(key []byte) []byte {
	hash := sha256.Sum256(key) // Always 32 bytes (AES-256)
	return hash[:]
}
