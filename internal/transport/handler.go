package transport

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thelazydo/email-verifier/internal/database"
	"github.com/thelazydo/email-verifier/types"
	"github.com/thelazydo/email-verifier/utils"
)

// Response is a generic API response container.
type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

// JobStatusResponse represents the data returned for a job's status.
type JobStatusResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Domain     string    `json:"domain"`
	Status     string    `json:"status"`
	Result     string    `json:"result,omitempty"`
	WebhookURL string    `json:"webhook_url,omitempty"`
}

// VerifyResponseData represents the data returned when a verification job is created.
type VerifyResponseData struct {
	JobID uuid.UUID `json:"job_id"`
}

// @Summary Get verification status of a job
// @Description Get the status of a verification job by its ID.
// @Tags status
// @Produce json
// @Param id path string true "Job ID" format(uuid)
// @Success 200 {object} Response{data=JobStatusResponse} "Job status retrieved successfully."
// @Failure 400 {object} Response{data=nil} "Invalid job ID format."
// @Failure 404 {object} Response{data=nil} "Job not found."
// @Failure 500 {object} Response{data=nil} "Failed to retrieve job."
// @Router /status/{id} [get]
func HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	idstring := r.PathValue("id")
	id, err := uuid.Parse(idstring)
	if err != nil {
		respondWithError(w, "invalid job id format", http.StatusBadRequest)
		return
	}

	job, err := database.GetJob(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, "job not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to get job", "error", err)
		respondWithError(w, "failed to retrieve job", http.StatusInternalServerError)
		return
	}
	response := JobStatusResponse{
		ID:     job.ID,
		Email:  job.Email,
		Domain: job.Domain,
		Status: string(job.Status),
	}
	if job.WebhookURL != nil {
		response.WebhookURL = *job.WebhookURL
	}

	if job.Result != nil {
		response.Result = string(*job.Result)
	}

	respondWithJSON(w, Response{
		Message: "job fetched successfully",
		Success: true,
		Data:    response,
	}, http.StatusOK)
}

// @Summary Verify an email address
// @Description Queues a job to verify an email address.
// @Tags verification
// @Accept json
// @Produce json
// @Param request body types.JobData true "Email verification payload"
// @Success 202 {object} Response{data=VerifyResponseData} "Job accepted for processing."
// @Failure 400 {object} Response{data=nil} "Invalid request payload or email address."
// @Failure 500 {object} Response{data=nil} "Failed to queue job."
// @Router /verify [post]
func HandlePostVerify(w http.ResponseWriter, r *http.Request) {
	var data types.JobData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	email, valid := utils.IsValidEmail(data.Address)
	if !valid {
		respondWithError(w, "invalid email address provided", http.StatusBadRequest)
		return
	}

	domain, _ := utils.GetDomain(data.Address)
	if utils.IsDisposable(domain) {
		respondWithError(w, "disposable email addresses are not allowed", http.StatusBadRequest)
		return
	}
	job := &types.Job{
		Email:      email,
		Domain:     domain,
		WebhookURL: &data.WebhookURL,
	}

	jobID, err := database.SaveJob(r.Context(), job)
	if err != nil {
		slog.Error("failed to save job", "error", err)
		respondWithError(w, "failed to queue job", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, Response{Message: "job received successfully", Success: true, Data: VerifyResponseData{JobID: jobID}}, http.StatusAccepted)
}

// @Summary API Information
// @Description Provides information about the available API endpoints.
// @Tags info
// @Produce json
// @Success 200 {object} Response "API information retrieved successfully."
// @Router / [get]
func HandleBaseRoute(w http.ResponseWriter, r *http.Request) {
	response := types.Response{
		Message: "Email Verifier API",
		Success: true,
		Data: map[string]string{
			"/verify": "Allowed methods - [POST]",
			"/status": "Allowed methods - [GET]",
		},
	}
	respondWithJSON(w, response, http.StatusOK)
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	respondWithJSON(w, Response{Message: message, Success: false}, statusCode)
}

func respondWithJSON(w http.ResponseWriter, payload any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}
