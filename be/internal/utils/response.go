package utils

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse sends a standard success response with status 200 OK
// Examples of usage:
// SuccessResponse(w, "Operation successful", data)
func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// CreatedResponse sends a standard created response with status 201 Created
// Examples of usage:
// CreatedResponse(w, "Resource created successfully", data)
func CreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// NoContentResponse sends a standard no content response with status 204 No Content
// Examples of usage:
// NoContentResponse(w)
func NoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "No Content",
	})
}

// BadRequestResponse sends a standard bad request response with status 400 Bad Request
// Examples of usage:
// BadRequestResponse(w, "Invalid request parameters")
func BadRequestResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// UnauthorizedResponse sends a standard unauthorized response with status 401 Unauthorized
// Examples of usage:
// UnauthorizedResponse(w, "Unauthorized access")
func UnauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// ForbiddenResponse sends a standard forbidden response with status 403 Forbidden
// Examples of usage:
// ForbiddenResponse(w, "Forbidden access")
func ForbiddenResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// NotFoundResponse sends a standard not found response with status 404 Not Found
// Examples of usage:
// NotFoundResponse(w, "Resource not found")
func NotFoundResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// ConflictResponse sends a standard conflict response with status 409 Conflict
// Examples of usage:
// ConflictResponse(w, "Conflict occurred")
func ConflictResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// UnprocessableEntityResponse sends a standard unprocessable entity response with status 422 Unprocessable Entity
// Examples of usage:
// UnprocessableEntityResponse(w, "Unprocessable entity")
func UnprocessableEntityResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// InternalServerErrorResponse sends a standard internal server error response with status 500 Internal Server Error
// Examples of usage:
// InternalServerErrorResponse(w, "Internal server error")
func InternalServerErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}
