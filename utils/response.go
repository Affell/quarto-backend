package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse crée une réponse de succès standardisée
func SuccessResponse(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// CreatedResponse crée une réponse de création standardisée
func CreatedResponse(c echo.Context, data any) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse crée une réponse d'erreur standardisée
func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// MessageResponse crée une réponse avec un message simple
func MessageResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
	})
}

// ValidationErrorResponse crée une réponse d'erreur de validation
func ValidationErrorResponse(c echo.Context, errors map[string]string) error {
	return c.JSON(http.StatusBadRequest, map[string]any{
		"success": false,
		"error":   "Erreurs de validation",
		"details": errors,
	})
}
