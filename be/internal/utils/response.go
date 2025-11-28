package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// sendResponse is a helper function to send JSON responses
// code   : HTTP status code (200, 400, 404, etc.)
// status : "success" or "error"
// message: message to send
// data   : payload data, can be nil if not needed
func sendResponse(c echo.Context, code int, status string, message string, data interface{}) error {
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
