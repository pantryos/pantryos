// helpers/response.go
package helpers

import (
	"encoding/json"
	"net/http"
	"time"
)

// APIError mendefinisikan struktur error yang lebih detail.
// Ini membantu frontend untuk melakukan penanganan error yang lebih spesifik.
type APIError struct {
	Code    string `json:"code,omitempty"`    // Kode error spesifik aplikasi, cth: "VALIDATION_ERROR"
	Details string `json:"details,omitempty"` // Pesan error yang lebih teknis untuk debugging
}

// APIResponse adalah struktur standar untuk semua respons JSON dari API.
type APIResponse struct {
	Success bool        `json:"success"`         // Menandakan apakah request berhasil atau tidak
	Message string      `json:"message"`         // Pesan yang bisa dibaca manusia
	Data    interface{} `json:"data,omitempty"`  // Payload data, diabaikan jika tidak ada
	Error   *APIError   `json:"error,omitempty"` // Detail error, diabaikan jika berhasil
}

// sendJSON adalah fungsi internal untuk mengirim respons JSON.
// Mencegah duplikasi kode pada fungsi Success dan Error.
func sendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Jika payload-nya nil, kita tidak perlu mengirim body
	if payload == nil {
		return
	}

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		// Jika ada error saat encoding JSON, log error tersebut
		// dan kirim respons HTTP 500 sebagai fallback.
		http.Error(w, "Failed to encode response to JSON", http.StatusInternalServerError)
	}
}

// Success mengirimkan respons berhasil (HTTP 2xx).
// `data` bisa berupa apa saja (struct, map, string, dll) atau nil.
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	sendJSON(w, statusCode, response)
}

// Error mengirimkan respons error (HTTP 4xx atau 5xx).
func Error(w http.ResponseWriter, statusCode int, message string, errDetails APIError) {
	response := APIResponse{
		Success: false,
		Message: message,
		Error:   &errDetails,
	}
	sendJSON(w, statusCode, response)
}

// getUserSuccessData
type GetUserSuccessData struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	AccountID int       `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
}
