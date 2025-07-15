package gameHandler

import (
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func All(prefix string) []models.Route {
	gameHandler := NewGameHandler()

	return []models.Route{
		{
			Path:    prefix + "/:id",
			Method:  echo.GET,
			Handler: gameHandler.GetGame,
		},
		{
			Path:    prefix + "/:id/select-piece",
			Method:  echo.POST,
			Handler: gameHandler.SelectPiece,
		},
		{
			Path:    prefix + "/:id/place-piece",
			Method:  echo.POST,
			Handler: gameHandler.PlacePiece,
		},
		{
			Path:    prefix + "/:id/forfeit",
			Method:  echo.POST,
			Handler: gameHandler.ForfeitGame,
		},
		{
			Path:    prefix + "/my",
			Method:  echo.GET,
			Handler: gameHandler.GetMyGames,
		},
	}
}
