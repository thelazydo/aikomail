package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/thelazydo/email-verifier/internal/middlewares"
)

func TestLimitMiddleware(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		next              http.HandlerFunc
		requestsPerSecond int
		want              http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := middlewares.LimitMiddleware(tt.next, tt.requestsPerSecond)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("LimitMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
