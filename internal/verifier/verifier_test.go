package verifier_test

import (
	"context"
	"net/smtp"
	"testing"

	"github.com/thelazydo/email-verifier/internal/verifier"
	"github.com/thelazydo/email-verifier/types"
)

func TestParseSMTPError(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		err  error
		want types.VerificationResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verifier.ParseSMTPError(tt.err)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ParseSMTPError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDialBest(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		domain  string
		want    *smtp.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := verifier.DialBest(tt.domain)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DialBest() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DialBest() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("DialBest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifySMTP(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		domain  string
		address string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := verifier.VerifySMTP(tt.domain, tt.address)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("VerifySMTP() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("VerifySMTP() succeeded unexpectedly")
			}
		})
	}
}

func TestIsCatchAll(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		domain  string
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := verifier.IsCatchAll(tt.domain)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("IsCatchAll() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("IsCatchAll() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("IsCatchAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasMX(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		domain  string
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := verifier.HasMX(context.Background(), tt.domain)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HasMX() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("HasMX() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("HasMX() = %v, want %v", got, tt.want)
			}
		})
	}
}
