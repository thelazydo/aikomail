package worker_test

import (
	"context"
	"sync"
	"testing"

	"github.com/thelazydo/email-verifier/internal/worker"
	"github.com/thelazydo/email-verifier/types"
)

func TestSendWebhook(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		url     string
		payload types.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker.SendWebhook(tt.url, tt.payload)
		})
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		numWorkers int
		wg         *sync.WaitGroup
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker.Start(context.Background(), tt.numWorkers, tt.wg)
		})
	}
}
