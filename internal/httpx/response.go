package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  []any  `json:"errors,omitempty"`
	Success bool   `json:"success"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, message string, data any) {
	resp := APIResponse{
		Status:  getStatusText(statusCode),
		Message: message,
		Data:    data,
		Success: statusCode >= 200 && statusCode < 300,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, statusCode int, message string, details []any) {
	resp := APIResponse{
		Status:  getStatusText(statusCode),
		Message: message,
		Errors:  details,
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(resp)
}

func getStatusText(statusCode int) string {
	switch statusCode {
		case http.StatusOK:
			return "OK"
		case http.StatusCreated:
			return "Created"
		case http.StatusBadRequest:
			return "BadRequest"
		case http.StatusUnauthorized:
			return "Unauthorized"
		case http.StatusForbidden:
			return "Forbidden"
		case http.StatusNotFound:
			return "NotFound"
		case http.StatusInternalServerError:
			return "InternalServerError"
		default:
			return "UnknownStatus"
	}
}
