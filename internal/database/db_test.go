package database_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thelazydo/email-verifier/internal/config"
	"github.com/thelazydo/email-verifier/internal/database"
	"github.com/thelazydo/email-verifier/types"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     *config.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := database.InitDB(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("InitDB() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("InitDB() succeeded unexpectedly")
			}
		})
	}
}

func TestStartCleanupTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		db *pgxpool.Pool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database.StartCleanupTask(tt.db)
		})
	}
}

func TestUpdateJobStatus(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      uuid.UUID
		status  types.VerificationResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := database.UpdateJobStatus(context.Background(), tt.id, tt.status)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateJobStatus() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateJobStatus() succeeded unexpectedly")
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		id      uuid.UUID
		want    *types.Job
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.GetJob(context.Background(), tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetJob() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetJob() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPendingJobs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		limit   int
		want    []types.Job
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.GetPendingJobs(context.Background(), tt.limit)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetPendingJobs() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetPendingJobs() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetPendingJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveJob(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		entry   *types.Job
		want    uuid.UUID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.SaveJob(context.Background(), tt.entry)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SaveJob() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SaveJob() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("SaveJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveMXCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		entry   database.CacheEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := database.SaveMXCache(context.Background(), tt.entry)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SaveMXCache() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SaveMXCache() succeeded unexpectedly")
			}
		})
	}
}

func TestCachedLookupMx(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		domain  string
		want    *database.CacheEntry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.CachedLookupMx(context.Background(), tt.domain)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CachedLookupMx() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CachedLookupMx() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CachedLookupMx() = %v, want %v", got, tt.want)
			}
		})
	}
}
