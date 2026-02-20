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

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}

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
	response := map[string]string{
		"id":     job.ID.String(),
		"email":  job.Email,
		"domain": job.Domain,
		"status": string(job.Status),
	}
	if job.WebhookURL != nil {
		response["webhook_url"] = *job.WebhookURL
	}

	if job.Result != nil {
		response["result"] = string(*job.Result)
	}

	respondWithJSON(w, Response{
		Message: "job fetched successfully",
		Success: true,
		Data:    response,
	}, http.StatusOK)
}

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
	respondWithJSON(w, Response{Message: "job received successfully", Success: true, Data: map[string]uuid.UUID{"job_id": jobID}}, http.StatusAccepted)
}

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
