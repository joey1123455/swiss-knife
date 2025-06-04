package swissknife

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const (
	letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GenerateUniqueString(length int) (string, error) {
	if length < 1 {
		return "", errors.New("length must be at least 1")
	}

	// Calculate entropy needed for uniqueness (128-bit minimum recommended)
	entropyBytes := make([]byte, length)
	_, err := rand.Read(entropyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	// Convert entropy to our character set
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random index: %w", err)
		}
		result[i] = letterBytes[num.Int64()]
	}

	return string(result), nil
}

func GenerateTimestampedUniqueStringID(length int) (string, error) {
	randPart, err := GenerateUniqueString(length)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), randPart), nil
}
