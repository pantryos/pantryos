// helpers/response.go
package helpers

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIError struct {
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

func sendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Failed to encode response to JSON", http.StatusInternalServerError)
	}
}

func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	sendJSON(w, statusCode, response)
}

func Error(w http.ResponseWriter, statusCode int, message string, errDetails APIError) {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   &errDetails,
	}
	sendJSON(w, statusCode, response)
}

type GetUserSuccessData struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	AccountID int       `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
}
