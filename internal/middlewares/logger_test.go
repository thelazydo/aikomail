package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/thelazydo/email-verifier/internal/middlewares"
)

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		next http.HandlerFunc
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := middlewares.LoggingMiddleware(tt.next)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("LoggingMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
