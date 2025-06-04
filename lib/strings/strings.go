package swissknife

import (
	"fmt"
	"strings"
)

func ExtractMainDomain(domain string) (string, error) {
	parts := strings.Split(domain, ".")
	if len(parts) < 3 {
		return "", fmt.Errorf("domain must have at least 3 parts")
	}
	return parts[1], nil // the middle part
}
