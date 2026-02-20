package config_test

import (
	"testing"

	"github.com/thelazydo/email-verifier/internal/config"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    *config.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := config.LoadConfig()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("LoadConfig() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("LoadConfig() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
