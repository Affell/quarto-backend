package handlers

import (
	"net/http"
	"quarto/config"
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func health(c echo.Context) error {

	return c.JSON(http.StatusOK, models.HealthModel{
		Status:  "Healthy !",
		Version: config.Version,
	})
}
