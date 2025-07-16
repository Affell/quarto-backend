package challengeHandler

import (
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func All(prefix string) []models.Route {

	return []models.Route{
		{
			Path:    prefix + "/send",
			Method:  echo.POST,
			Handler: sendChallenge,
		},
		{
			Path:    prefix + "/respond",
			Method:  echo.POST,
			Handler: respondToChallenge,
		},
		{
			Path:    prefix + "/my",
			Method:  echo.GET,
			Handler: getMyChallenges,
		},
	}
}
