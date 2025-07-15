package user

import (
	"net/mail"
	"regexp"
	"strings"
)

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidPassword(password string) (err string) {

	containUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	containLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	containSpecialChar := strings.ContainsAny(password, " \"!#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	containNumeric := strings.ContainsAny(password, "0123456789")

	if len(password) < 8 {
		err = "Password must be at least 8 characters long"
	} else if !containLower {
		err = "Password must contain a lowercase letter"
	} else if !containSpecialChar {
		err = "Password must contain a special character"
	} else if !containNumeric {
		err = "Password must contain a number"
	} else if !containUpper {
		err = "Password must contain an uppercase letter"
	}

	return
}

func ValidUsername(username string) (err string) {
	var re = regexp.MustCompile(`^(?mi)([\w.\-]+)$`)
	if len(username) < 3 || len(username) > 20 || !re.MatchString(username) {
		return "le nom d'utilisateur ne doit contenir que des caractères alphanumériques et '-' ou '_' et doit avoir entre 3 et 20 caractères"
	}

	return
}
func ValidBase64(base64 string) bool {
	var re = regexp.MustCompile(`[^-A-Za-z0-9+/=]|=[^=]|={3,}$`)
	return re.MatchString(base64)
}
