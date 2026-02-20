package utils

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

var (
	disposableDomains = make(map[string]struct{})
	emailRegex        = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func GetDomain(address string) (string, error) {
	_, domain, found := strings.Cut(address, "@")
	if !found {
		slog.Info("invalid address", "address", address)
		return "", fmt.Errorf("invalid adress: %s", address)
	}
	return domain, nil
}

func LoadDisposableDomains() {
	// In production, download or read from a local file
	list := []string{
		"mailinator.com", "10minuteemail.com", "guerrillamail.com", "yopmail.com",
		"temp-mail.org", "getnada.com",
	}
	for _, d := range list {
		disposableDomains[d] = struct{}{}
	}
}

func IsDisposable(domain string) bool {
	_, exists := disposableDomains[strings.ToLower(domain)]
	return exists
}

func RandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		slog.Error("failed while generating random string", "error", err)
		return "abcdef"
	}
	return fmt.Sprintf("%x", b)
}

func IsValidEmail(address string) (string, bool) {
	clean := strings.ToLower(strings.TrimSpace(address))

	if len(clean) < 3 || len(clean) > 254 {
		return "", false
	}
	return clean, emailRegex.MatchString(clean)
}
