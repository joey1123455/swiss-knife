package swissknife

import (
	"crypto/x509"
	"fmt"
	"os"
)

func loadCertPool(certFile string) (*x509.CertPool, error) {
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert %s: %w", certFile, err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certData)
	return pool, nil
}

// ====================================== logger ===============================
