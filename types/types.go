package types

import (
	"time"

	"github.com/google/uuid"
)

type (
	VerificationResult string
	JobStatus          string

	Job struct {
		ID         uuid.UUID           `json:"id"`
		Email      string              `json:"email"`
		Domain     string              `json:"domain"`
		Status     JobStatus           `json:"status"`
		Result     *VerificationResult `json:"result"`
		WebhookURL *string             `json:"webhook_url"`
		CreatedAt  time.Time           `json:"created_at"`
		UpdatedAt  time.Time           `json:"updated_at"`
	}

	Result struct {
		JobID  uuid.UUID `json:"job_id"`
		Status string    `json:"status"`
	}

	Response struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
		Data    any    `json:"data"`
	}

	JobData struct {
		Address    string `json:"address"`
		WebhookURL string `json:"webhookURL"`
	}
)

const (
	Completed  JobStatus = "completed"
	Pending    JobStatus = "pending"
	Processing JobStatus = "processing"
	Failed     JobStatus = "failed"
)

const (
	ResultValid      VerificationResult = "valid"
	ResultInvalid    VerificationResult = "invalid"
	ResultDisposable VerificationResult = "disposable"
	ResultCatchAll   VerificationResult = "catch-all"
	ResultUnknown    VerificationResult = "unknown"
	ResultError      VerificationResult = "error"
)
