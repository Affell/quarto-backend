package authHandler

import (
	"errors"
	"net/http"
	"quarto/models/user"

	"github.com/labstack/echo/v4"
)

type SignupForm struct {
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"Password123!"`
}

type SignupResponse struct {
	Message string `json:"message" example:"Signup successful"`
}

type SignupError400 struct {
	Message string `json:"message" example:"Please fully fill in the signup form"`
}

type SignupError409 struct {
	Message string `json:"message" example:"Email not available"`
}

// @Summary Signup a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param signupForm body SignupForm true "Signup form"
// @Success 201 {object} SignupResponse
// @Failure 400 {object} SignupError400
// @Failure 409 {object} SignupError409
// @Router /auth/signup [post]
func signup(c echo.Context) error {

	var signupForm SignupForm

	if err := c.Bind(&signupForm); err != nil && signupForm.Username == "" || signupForm.Email == "" || signupForm.Password == "" {
		return c.JSON(400, SignupError400{Message: "Please fully fill in the signup form"})
	}

	if !user.ValidEmail(signupForm.Email) {
		return c.JSON(http.StatusBadRequest, SignupError400{Message: "Invalid Email"})
	}

	if !user.CheckEmailAvailability(signupForm.Email) {
		return c.JSON(http.StatusConflict, SignupError409{Message: "Email not available"})
	}

	if err := user.ValidPassword(signupForm.Password); err != "" {
		return c.JSON(http.StatusBadRequest, SignupError400{Message: err})
	}

	var id int64
	if id = user.CreateAccount(signupForm.Email, signupForm.Username, signupForm.Password); id == -1 {
		return errors.New("error creating account")
	}

	return c.JSON(201, SignupResponse{Message: "Signup successful"})
}
