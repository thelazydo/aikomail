package main

import (
	"net/http"
	"testing"

	"github.com/thelazydo/email-verifier/internal/config"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_run(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := run()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("run() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("run() succeeded unexpectedly")
			}
		})
	}
}

func Test_routes(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg  *config.Config
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := routes(tt.cfg)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("routes() = %v, want %v", got, tt.want)
			}
		})
	}
}
