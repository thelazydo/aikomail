package transport_test

import (
	"net/http"
	"testing"

	"github.com/thelazydo/email-verifier/internal/transport"
)

func TestHandleGetStatus(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w http.ResponseWriter
		r *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport.HandleGetStatus(tt.w, tt.r)
		})
	}
}

func TestHandlePostVerify(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w http.ResponseWriter
		r *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport.HandlePostVerify(tt.w, tt.r)
		})
	}
}

func TestHandleBaseRoute(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w http.ResponseWriter
		r *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport.HandleBaseRoute(tt.w, tt.r)
		})
	}
}
