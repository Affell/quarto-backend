package authHandler

import (
	"net/http"
	"quarto/config"
	"quarto/email"
	"quarto/models/user"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

type AskRecoverForm struct {
	Email string `json:"email"`
}

type RecoverResponse struct {
	Message string `json:"message" example:"Recovery token created successfully"`
}

type RecoverError400 struct {
	Message string `json:"message" example:"Please fully fill in the account recovery form"`
}

// @Summary Account Recovery
// @Description Generates a recovery token for the user to reset their password.
// @Tags auth
// @Accept json
// @Produce json
// @Param askRecoverForm body AskRecoverForm true "Account recovery form"
// @Success 201 {object} RecoverResponse
// @Failure 400 {object} RecoverError400
// @Router /auth/recover [post]
func recover(c echo.Context) error {

	var askRecoverForm AskRecoverForm
	if err := c.Bind(&askRecoverForm); err != nil || askRecoverForm.Email == "" {
		return c.JSON(http.StatusBadRequest, RecoverError400{Message: "Please fully fill in the account recovery form"})
	}

	recoverToken, err := user.CreateRecoverToken(askRecoverForm.Email)
	if err != nil {
		return err
	}

	if recoverToken != "" {
		// Send email with recover token
		btnURL := config.Config.FrontURL + "/reset_password/?recover_token=" + recoverToken
		subject := "Quarto - Password recovery"
		text := "A password reset request has been made for your Quarto account using your email address. " +
			"To complete this process, click the button below :"
		btnText := "RESET PASSWORD"

		err = email.New(askRecoverForm.Email, subject, subject, text, btnText, btnURL).Send(config.Config.Email)
		if err != nil {
			log.Error("Failed to send reset email", "email", askRecoverForm.Email, "error", err)
			return c.JSON(http.StatusInternalServerError, RecoverError400{Message: "Failed to send reset email, please contact support at support@quarto.fr"})
		}
	}

	return c.NoContent(http.StatusCreated)
}
