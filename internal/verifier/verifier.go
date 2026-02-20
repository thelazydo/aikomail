package verifier

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thelazydo/email-verifier/internal/database"
	"github.com/thelazydo/email-verifier/types"
)

const fromAddress = "noreply@example.com"

// MX
func HasMX(ctx context.Context, domain string) (bool, error) {
	entry, err := database.CachedLookupMx(ctx, domain)
	if err == nil && entry != nil && len(entry.Records) > 0 {
		return true, nil
	}
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
			return false, nil
		}
		return false, fmt.Errorf("mx lookup failed: %v", err)
	}

	if len(mxRecords) == 0 {
		return false, nil
	}

	// cache successful lookup
	database.SaveMXCache(ctx, database.CacheEntry{
		Domain:  domain,
		Records: mxRecords,
		Expiry:  time.Now().Add(1 * time.Hour),
	})

	return true, nil
}

// Detect
func IsCatchAll(domain string) (bool, error) {
	randomEmail := fmt.Sprintf("%s@%s", uuid.New().String(), domain)
	err := VerifySMTP(domain, randomEmail)
	if err == nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "550") {
		return false, nil
	}
	return false, err
}

func VerifySMTP(domain string, address string) error {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return fmt.Errorf("dns lookup failed: %w", err)
	}

	if len(mxRecords) == 0 {
		return fmt.Errorf("no mx records found")
	}

	var lastErr error
	for _, mx := range mxRecords {
		host := strings.TrimSuffix(mx.Host, ".")
		addr := fmt.Sprintf("%s:%d", host, 25)

		err := tryHost(addr, host, address)
		if err == nil {
			return nil
		}
		lastErr = err
		if smtpErr, ok := err.(*textproto.Error); ok && smtpErr.Code >= 500 {
			return lastErr
		}
	}

	return lastErr
}

// HANDSHAKE
func DialBest(domain string) (*smtp.Client, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for _, mxRecord := range mxRecords {
		address := net.JoinHostPort(mxRecord.Host, "25")

		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			lastErr = err
			continue
		}

		client, err := smtp.NewClient(conn, mxRecord.Host)
		if err != nil {
			conn.Close()
			lastErr = err
			continue
		}
		return client, nil

	}
	return nil, fmt.Errorf("all mx servers failed: %v", lastErr)
}

func tryHost(addr string, host, address string) error {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)

	if err := client.Hello("localhost"); err != nil {
		return err
	}

	if err := client.Mail(fromAddress); err != nil {
		return err
	}

	if err := client.Rcpt(address); err != nil {
		return err
	}

	return nil
}

func ParseSMTPError(err error) types.VerificationResult {
	if smtpErr, ok := err.(*textproto.Error); ok {
		if smtpErr.Code >= 500 {
			return types.ResultInvalid
		}
		if smtpErr.Code >= 400 {
			return types.ResultUnknown
		}
	}
	return types.ResultError
}
