package handlers

import (
	"quarto/handlers/authHandler"
	"quarto/handlers/challengeHandler"
	"quarto/handlers/gameHandler"
	"quarto/handlers/userHandler"
	"quarto/handlers/websocketHandler"
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func All() (routes []models.Route) {
	routes = append(routes, models.Route{
		Path:    "/",
		Method:  echo.GET,
		Handler: health,
	})

	routes = append(routes, authHandler.All("/auth")...)
	routes = append(routes, gameHandler.All("/game")...)
	routes = append(routes, challengeHandler.All("/challenge")...)
	routes = append(routes, userHandler.All("/users")...)
	routes = append(routes, websocketHandler.All()...)

	return
}
