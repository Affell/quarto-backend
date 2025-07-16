package gameHandler

import (
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func All(prefix string) []models.Route {
	return []models.Route{
		{
			Path:    prefix + "/:id",
			Method:  echo.GET,
			Handler: getGame,
		},
		{
			Path:    prefix + "/:id/select-piece",
			Method:  echo.POST,
			Handler: selectPiece,
		},
		{
			Path:    prefix + "/:id/place-piece",
			Method:  echo.POST,
			Handler: placePiece,
		},
		{
			Path:    prefix + "/:id/forfeit",
			Method:  echo.POST,
			Handler: forfeitGame,
		},
		{
			Path:    prefix + "/my",
			Method:  echo.GET,
			Handler: getMyGames,
		},
	}
}
