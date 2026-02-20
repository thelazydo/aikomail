package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thelazydo/email-verifier/internal/config"
	"github.com/thelazydo/email-verifier/types"
)

var DB *pgxpool.Pool

type CacheEntry struct {
	Domain  string    `json:"domain"`
	Records []*net.MX `json:"records"`
	Expiry  time.Time `json:"expiry"`
}

func InitDB(cfg *config.Config) error {
	var err error
	DB, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	if err := DB.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}
	slog.Info("database connection successful")
	StartCleanupTask(DB)

	return migrate()
}

func StartCleanupTask(db *pgxpool.Pool) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			tag, err := db.Exec(ctx, "DELETE FROM verification_jobs WHERE created_at < NOW()- INTERVAL '30 days'")
			if err != nil {
				slog.Error("cleanup failed", "error", err)
			} else {
				slog.Info("cleanup successful", "rows_deleted", tag.RowsAffected())
			}
		}
	}()
}

func UpdateJobStatus(ctx context.Context, id uuid.UUID, status types.VerificationResult) error {
	query := "UPDATE verification_jobs SET status='completed',updated_at=NOW(), result=$1 WHERE id=$2;"

	_, err := DB.Exec(ctx,
		query,
		status,
		id,
	)
	return err
}

func GetJob(ctx context.Context, id uuid.UUID) (*types.Job, error) {
	var job types.Job
	query := "SELECT * FROM verification_jobs WHERE id=$1;"

	err := DB.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.Email,
		&job.Domain,
		&job.Status,
		&job.Result,
		&job.WebhookURL,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, err
}

func GetPendingJobs(ctx context.Context, limit int) ([]types.Job, error) {
	query := `
	UPDATE verification_jobs
	SET status = 'processing'
	WHERE id IN (
		SELECT id
		FROM verification_jobs
		WHERE status = 'pending'
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED
		LIMIT $1
	)
	RETURNING id, email, domain, status, result, webhook_url, created_at, updated_at;
	`

	rows, err := DB.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Job])
}

func SaveJob(ctx context.Context, entry *types.Job) (uuid.UUID, error) {
	query := `
	INSERT INTO verification_jobs (email, domain, webhook_url)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	var id uuid.UUID
	err := DB.QueryRow(ctx, query, &entry.Email, &entry.Domain, &entry.WebhookURL).Scan(&id)

	return id, err
}

func SaveMXCache(ctx context.Context, entry CacheEntry) error {
	query := `
	INSERT INTO mx_cache (domain, records, expiry)
	VALUES ($1, $2, $3)
	ON CONFLICT (domain)
	DO UPDATE SET records = $2, expiry = $3;
	`
	_, err := DB.Exec(ctx, query, entry.Domain, entry.Records, entry.Expiry)
	return err
}

func CachedLookupMx(ctx context.Context, domain string) (*CacheEntry, error) {
	var entry CacheEntry
	var records []byte
	query := "SELECT domain, records, expiry FROM mx_cache WHERE domain = $1 AND expiry > NOW()"

	err := DB.QueryRow(ctx, query, domain).Scan(&entry.Domain, &records, &entry.Expiry)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(records, &entry.Records); err != nil {
		return nil, err
	}

	return &entry, nil
}

func migrate() error {
	// DROP TABLE IF EXISTS verification_jobs;
	// DROP TABLE IF EXISTS mx_cache;
	query := `
	CREATE TABLE IF NOT EXISTS verification_jobs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) NOT NULL,
		domain VARCHAR(255) NOT NULL,
		status VARCHAR(255) DEFAULT 'pending',
		result TEXT,
		webhook_url TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS mx_cache (
		domain TEXT PRIMARY KEY,
		records JSONB NOT NULL,
		expiry TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_status ON verification_jobs(status);
	CREATE INDEX IF NOT EXISTS idx_mx_created_at ON verification_jobs(created_at);
	CREATE INDEX IF NOT EXISTS idx_mx_cache_expiry ON mx_cache(expiry);
	`
	_, err := DB.Exec(context.Background(), query)
	return err
}
