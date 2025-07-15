package challengeHandler

import (
	"quarto/models"

	"github.com/labstack/echo/v4"
)

func All(prefix string) []models.Route {
	challengeHandler := NewChallengeHandler()

	return []models.Route{
		{
			Path:    prefix + "/send",
			Method:  echo.POST,
			Handler: challengeHandler.SendChallenge,
		},
		{
			Path:    prefix + "/respond",
			Method:  echo.POST,
			Handler: challengeHandler.RespondToChallenge,
		},
		{
			Path:    prefix + "/my",
			Method:  echo.GET,
			Handler: challengeHandler.GetMyChallenges,
		},
	}
}
