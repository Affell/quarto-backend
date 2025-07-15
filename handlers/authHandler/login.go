package authHandler

import (
	"errors"
	"net/http"
	"quarto/models/user"

	"github.com/labstack/echo/v4"
)

type LoginForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Token    string `form:"token" json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int64  `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"user"`
}

type LoginError400 struct {
	Message string `json:"message" example:"Please fully fill in the login form"`
}

type LoginError403 struct {
	Message string `json:"message" example:"Invalid email or password"`
}

// login handles the user login process.
// @Summary User login
// @Description Logs in a user using email/password or token.
// @Tags auth
// @Accept json
// @Produce json
// @Param loginForm body LoginForm true "Login form"
// @Success 200 {object} LoginResponse "token and user details"
// @Failure 400 {object} LoginError400
// @Failure 403 {object} LoginError403
// @Router /auth/login [post]
func login(c echo.Context) error {

	var loginForm LoginForm
	if err := c.Bind(&loginForm); err != nil || ((loginForm.Password == "" || loginForm.Email == "") && loginForm.Token == "") {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Please fully fill in the login form"})
	}

	var CurrentUserToken user.UserToken
	var err error
	if loginForm.Token == "" && loginForm.Password != "" && loginForm.Email != "" {
		CurrentUserToken, err = user.GetSQLUserToken(loginForm.Email, loginForm.Password)
		if err != nil {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Invalid email or password"})
		}
	} else if loginForm.Token != "" && loginForm.Password == "" && loginForm.Email == "" {
		CurrentUserToken, err = user.GetUserToken(loginForm.Token)
		if err != nil {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Invalid token"})
		}
	} else {
		return c.NoContent(http.StatusBadRequest)
	}

	TokenID := CurrentUserToken.Store()
	if TokenID == "" {
		return errors.New("error during token storage")
	}

	return c.JSON(http.StatusOK, map[string]any{"token": TokenID, "user": CurrentUserToken.User.ToSelfWebDetail()})
}
