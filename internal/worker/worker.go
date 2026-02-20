package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/thelazydo/email-verifier/internal/database"
	"github.com/thelazydo/email-verifier/internal/verifier"
	"github.com/thelazydo/email-verifier/types"
	"github.com/thelazydo/email-verifier/utils"
)

// WORKER
func SendWebhook(url string, payload types.Result) {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal webhook payload", "job_id", payload.JobID, "error", err)
		return
	}

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			slog.Error("failed to create webhook request", "job_id", payload.JobID, "error", err)
			time.Sleep(time.Second * time.Duration(2*i+1))
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			slog.Info("webhook delivered", "job_id", payload.JobID)
			resp.Body.Close()
			return
		}
		if err != nil {
			slog.Error("webhook post error", "job_id", payload.JobID, "error", err)
		} else {
			slog.Error("webhook failed with non-2xx response", "job_id", payload.JobID, "status_code", resp.StatusCode)
			resp.Body.Close()
		}
		time.Sleep(time.Second * time.Duration(2*i+1))
	}
	slog.Error("webhook failed after 3 retries", "job_id", payload.JobID)
}

func Start(ctx context.Context, numWorkers int, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			slog.Info("starting worker", "worker_id", workerID)
			for {
				processJobs(workerID)
				select {
				case <-ctx.Done():
					slog.Info("worker shutting down", "worker_id", workerID)
					return

				case <-time.After(5 * time.Second):
				}
				// time.Sleep(5 * time.Second)

			}
		}(i + 1)
	}
}

func processJobs(workerID int) {
	ctx, cancle := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancle()

	jobs, err := database.GetPendingJobs(ctx, 10)
	if err != nil {
		slog.Error("failed to get pending jobs", "worker_id", workerID, "error", err)
		return
	}
	if len(jobs) == 0 {
		return
	}

	slog.Info("worker has jobs", "worker_id", workerID, "num_jobs", len(jobs))
	for _, job := range jobs {
		processJob(ctx, job)
	}
}

func processJob(ctx context.Context, job types.Job) {
	result := performVerification(job.Email)
	slog.Info("job completed", "job_id", job.ID, "result", result)
	if err := database.UpdateJobStatus(ctx, job.ID, result); err != nil {
		slog.Error("failed to update job status", "job_id", job.ID, "error", err)
		return
	}

	// notify client if webhook url exists
	if *job.WebhookURL != "" {
		go SendWebhook(*job.WebhookURL, types.Result{
			JobID:  job.ID,
			Status: string(result),
		})
	}

	slog.Info("process completed", "job_id", job.ID, "result", result)
}

func performVerification(address string) types.VerificationResult {
	domain, err := utils.GetDomain(address)
	if err != nil {
		return verifier.ParseSMTPError(err)
	}

	hasMX, err := verifier.HasMX(context.Background(), domain)
	if err != nil {
		slog.Error("mx lookup failed", "domain", domain, "error", err)
		return types.ResultUnknown
	}
	if !hasMX {
		return types.ResultInvalid
	}

	isCatchAll, err := verifier.IsCatchAll(domain)
	if err != nil {
		slog.Warn("catch-all check failed", "domain", domain, "error", err)
		return types.ResultError
	}

	if isCatchAll {
		return types.ResultCatchAll
	}

	if err := verifier.VerifySMTP(domain, address); err != nil {
		return verifier.ParseSMTPError(err)
	}
	return types.ResultValid
}
