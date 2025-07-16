package challengeHandler

import (
	"net/http"
	"quarto/models/challenge"
	"quarto/models/user"

	"github.com/labstack/echo/v4"
)

// sendChallenge envoie un défi à un autre joueur
// @Summary Send challenge
// @Description Send a challenge to another player
// @Tags challenges
// @Accept json
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param request body challenge.SendChallengeRequest true "Send challenge request"
// @Success 201 {object} challenge.Challenge
// @Failure 400 {object} map[string]string
// @Router /challenge/send [post]
func sendChallenge(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	var req challenge.SendChallengeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Données invalides")
	}

	newChallenge, err := challenge.SendChallenge(userToken.User.ID, req.ChallengedID, req.Message)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, newChallenge.ToWeb())
}

// respondToChallenge répond à un défi (accepter ou refuser)
// @Summary Respond to challenge
// @Description Accept or decline a challenge
// @Tags challenges
// @Accept json
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param request body challenge.RespondToChallengeRequest true "Response to challenge"
// @Success 200 {object} challenge.ChallengeResponse
// @Failure 400 {object} map[string]string
// @Router /challenge/respond [post]
func respondToChallenge(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	var req challenge.RespondToChallengeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Données invalides")
	}

	if req.Accept {
		// Accepter le défi
		updatedChallenge, newGame, err := challenge.AcceptChallenge(req.ChallengeID, userToken.User.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		response := challenge.ChallengeResponse{
			Challenge: updatedChallenge,
			Game:      newGame.ToWeb(),
		}

		return c.JSON(http.StatusOK, response)
	} else {
		// Refuser le défi
		updatedChallenge, err := challenge.DeclineChallenge(req.ChallengeID, userToken.User.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		response := challenge.ChallengeResponse{
			Challenge: updatedChallenge,
			Game:      nil,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// getMyChallenges récupère tous les défis de l'utilisateur
// @Summary Get my challenges
// @Description Get all challenges sent and received by the user
// @Tags challenges
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Success 200 {object} challenge.ChallengeListResponse
// @Failure 401 {object} map[string]string
// @Router /challenge/my [get]
func getMyChallenges(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	challenges, err := challenge.GetMyChallenges(userToken.User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, challenges)
}
