package authHandler

import (
	"net/http"
	"quarto/models/user"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type MeError struct {
	Message string `json:"message" example:"Incorrect token"`
}

// @Summary Get user details
// @Description Get details of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Success 200 {object} UserResponse
// @Failure 401 {object} MeError
// @Router /auth/me [get]
func me(c echo.Context) error {

	var token user.UserToken
	if t := c.Get("userToken"); t != nil {
		token = t.(user.UserToken)
	} else {
		return c.JSON(http.StatusUnauthorized, MeError{Message: "Incorrect token"})
	}

	u, err := user.GetUserById(token.User.ID)
	if err == pgx.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, MeError{Message: "Incorrect token"})
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, u.ToSelfWebDetail())
}
