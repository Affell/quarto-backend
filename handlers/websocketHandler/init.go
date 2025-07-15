package websocketHandler

import (
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func All() []models.Route {
	handler := NewWebSocketHandler()

	return []models.Route{
		{
			Path:    "/ws",
			Method:  echo.GET,
			Handler: handler.HandleWebSocket,
		},
	}
}
