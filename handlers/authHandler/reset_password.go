package authHandler

import (
	"net/http"
	"quarto/models/user"

	"github.com/labstack/echo/v4"
)

type PasswordResetForm struct {
	RecoverToken string `json:"token"`
	Password     string `json:"password"`
}

type PasswordResetResponse struct {
	Message string `json:"message" example:"Password reset successfully"`
}

type PasswordResetError400 struct {
	Message string `json:"message" example:"Please fully fill in the password reset form"`
}

type PasswordResetError403 struct {
	Message string `json:"message" example:"Invalid recover token"`
}

// @Summary Reset user password
// @Description Reset the password for a user using a recovery token
// @Tags auth
// @Accept json
// @Produce json
// @Param passwordResetForm body PasswordResetForm true "Password Reset Form"
// @Success 201 {object} PasswordResetResponse
// @Failure 400 {object} PasswordResetError400
// @Failure 403 {object} PasswordResetError403
// @Router /auth/reset_password [post]
func resetPassword(c echo.Context) error {

	var passwordResetForm PasswordResetForm
	if err := c.Bind(&passwordResetForm); err != nil || passwordResetForm.RecoverToken == "" || passwordResetForm.Password == "" {
		return c.JSON(http.StatusBadRequest, PasswordResetError400{Message: "Please fully fill in the password reset form"})
	}

	if msg := user.ValidPassword(passwordResetForm.Password); msg != "" {
		return c.JSON(http.StatusBadRequest, PasswordResetError400{Message: msg})
	}

	if ok := user.ResetPassword(passwordResetForm.RecoverToken, passwordResetForm.Password); !ok {
		return c.JSON(http.StatusForbidden, PasswordResetError403{Message: "Invalid recover token"})
	}

	return c.JSON(http.StatusCreated, PasswordResetResponse{Message: "Password reset successfully"})
}
