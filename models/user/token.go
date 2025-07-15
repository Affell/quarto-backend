package user

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	tokens map[string]UserToken = make(map[string]UserToken)
	mutex  sync.Mutex
)

type UserToken struct {
	TokenID   string    `json:"token_id" structs:"-"`
	User      User      `json:"user" structs:"-"`
	CreatedAt time.Time `json:"created_at" structs:"-"`
}

func (token UserToken) IsNil() bool {
	return token.TokenID == ""
}

func (userToken *UserToken) Store() (tokenID string) {

	if userToken.TokenID == "" {
		tokenID = uuid.New().String()
		userToken.TokenID = tokenID
	} else {
		tokenID = userToken.TokenID
	}

	mutex.Lock()
	tokens[userToken.TokenID] = *userToken
	mutex.Unlock()

	return
}

func GetUserToken(tokenID string) (userToken UserToken, err error) {

	mutex.Lock()
	userToken, ok := tokens[tokenID]
	mutex.Unlock()
	if !ok {
		err = errors.New("incorrect token")
		return
	}

	if time.Since(userToken.CreatedAt) > TOKEN_EXPIRATION {
		RevokeUserToken(tokenID)
		userToken = UserToken{}
		err = errors.New("token expired")
	}

	return
}

func RevokeUserToken(tokenID string) {
	mutex.Lock()
	delete(tokens, tokenID)
	mutex.Unlock()
}

func GetTokenFromRequest(c echo.Context) (userToken UserToken, err error) {

	if t := c.Get("userToken"); t != nil {
		userToken = t.(UserToken)
	} else {
		err = errors.New("invalid session")
		return
	}

	if time.Since(userToken.CreatedAt) > TOKEN_EXPIRATION {
		RevokeUserToken(userToken.TokenID)
		userToken = UserToken{}
		err = errors.New("token expired")
	}

	return
}
