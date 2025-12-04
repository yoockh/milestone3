package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
}

// Response represents the standard API response structure
type Response struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse represents a successful API response
type SuccessResponseData struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message"`
}

// sendResponse is a helper function to send JSON responses with logging
func sendResponse(c echo.Context, code int, status string, message string, data interface{}) error {
	fields := logrus.Fields{
		"method": c.Request().Method,
		"path":   c.Request().URL.Path,
		"status": code,
	}

	if userID, ok := GetUserID(c); ok {
		fields["user_id"] = userID
	}

	if code >= 500 {
		logger.WithFields(fields).Error(message)
	} else if code >= 400 {
		logger.WithFields(fields).Warn(message)
	} else if code >= 200 && code < 300 {
		logger.WithFields(fields).Info(message)
	}

	resp := map[string]interface{}{
		"status":  status,
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}
	return c.JSON(code, resp)
}

// SuccessResponse sends a standard success response with HTTP status 200 OK
// Example usage:
// return utils.SuccessResponse(c, "Data fetched successfully", data)
func SuccessResponse(c echo.Context, message string, data interface{}) error {
	return sendResponse(c, http.StatusOK, "success", message, data)
}

// CreatedResponse sends a standard success response with HTTP status 201 Created
// Example usage:
// return utils.CreatedResponse(c, "User created successfully", user)
func CreatedResponse(c echo.Context, message string, data interface{}) error {
	return sendResponse(c, http.StatusCreated, "success", message, data)
}

// NoContentResponse sends a standard success response with HTTP status 204 No Content
// Example usage:
// return utils.NoContentResponse(c)
func NoContentResponse(c echo.Context) error {
	return sendResponse(c, http.StatusNoContent, "success", "No Content", nil)
}

// BadRequestResponse sends a standard error response with HTTP status 400 Bad Request
// Example usage:
// return utils.BadRequestResponse(c, "Invalid request parameters")
func BadRequestResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusBadRequest, "error", message, nil)
}

// UnauthorizedResponse sends a standard error response with HTTP status 401 Unauthorized
// Example usage:
// return utils.UnauthorizedResponse(c, "Unauthorized access")
func UnauthorizedResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusUnauthorized, "error", message, nil)
}

// ForbiddenResponse sends a standard error response with HTTP status 403 Forbidden
// Example usage:
// return utils.ForbiddenResponse(c, "Forbidden access")
func ForbiddenResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusForbidden, "error", message, nil)
}

// NotFoundResponse sends a standard error response with HTTP status 404 Not Found
// Example usage:
// return utils.NotFoundResponse(c, "Resource not found")
func NotFoundResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusNotFound, "error", message, nil)
}

// ConflictResponse sends a standard error response with HTTP status 409 Conflict
// Example usage:
// return utils.ConflictResponse(c, "Conflict occurred")
func ConflictResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusConflict, "error", message, nil)
}

// UnprocessableEntityResponse sends a standard error response with HTTP status 422 Unprocessable Entity
// Example usage:
// return utils.UnprocessableEntityResponse(c, "Unprocessable entity")
func UnprocessableEntityResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusUnprocessableEntity, "error", message, nil)
}

// InternalServerErrorResponse sends a standard error response with HTTP status 500 Internal Server Error
// Example usage:
// return utils.InternalServerErrorResponse(c, "Internal server error")
func InternalServerErrorResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusInternalServerError, "error", message, nil)
}
