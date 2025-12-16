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

type SuccessResponseData struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

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

	// Admin-only context (set by JWT middleware)
	if adminID, ok := GetAdminID(c); ok {
		fields["admin_id"] = adminID
	}
	if role, ok := GetRole(c); ok {
		fields["role"] = role
	}

	if code >= 500 {
		logger.WithFields(fields).Error(message)
	} else if code >= 400 {
		logger.WithFields(fields).Warn(message)
	} else {
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

func SuccessResponse(c echo.Context, message string, data interface{}) error {
	return sendResponse(c, http.StatusOK, "success", message, data)
}

func CreatedResponse(c echo.Context, message string, data interface{}) error {
	return sendResponse(c, http.StatusCreated, "success", message, data)
}

func NoContentResponse(c echo.Context) error {
	return sendResponse(c, http.StatusNoContent, "success", "No Content", nil)
}

func BadRequestResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusBadRequest, "error", message, nil)
}

func UnauthorizedResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusUnauthorized, "error", message, nil)
}

func ForbiddenResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusForbidden, "error", message, nil)
}

func NotFoundResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusNotFound, "error", message, nil)
}

func ConflictResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusConflict, "error", message, nil)
}

func UnprocessableEntityResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusUnprocessableEntity, "error", message, nil)
}

func InternalServerErrorResponse(c echo.Context, message string) error {
	return sendResponse(c, http.StatusInternalServerError, "error", message, nil)
}
