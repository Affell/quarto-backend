package userHandler

import (
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func All(prefix string) (routes []models.Route) {
	userHandler := NewUserHandler()

	routes = append(routes, models.Route{
		Path:    prefix,
		Method:  echo.GET,
		Handler: userHandler.GetUsers,
	})

	routes = append(routes, models.Route{
		Path:    prefix + "/:id",
		Method:  echo.GET,
		Handler: userHandler.GetUser,
	})

	return
}
