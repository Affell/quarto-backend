package user

import (
	"time"

	"github.com/fatih/structs"
)

const TOKEN_EXPIRATION = time.Hour * 24 * 30

type (
	User struct {
		ID           int64  `structs:"id"`
		Email        string `structs:"email"`
		Username     string `structs:"username"`
		Password     string `structs:"-"`
		RecoverToken string `structs:"-"`
		Enable       bool   `structs:"-"`
		Admin        bool   `structs:"-"`
	}

	UserList []User

	// Structure pour la pagination
	UserPaginationResponse struct {
		Users      []UserPublic `json:"users"`
		Page       int          `json:"page"`
		PageSize   int          `json:"page_size"`
		Total      int64        `json:"total"`
		TotalPages int          `json:"total_pages"`
	}

	// Structure publique pour les autres utilisateurs
	UserPublic struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}
)

func (user User) ToSelfWebDetail() map[string]any {
	return structs.Map(user)
}

func (user User) ToWeb() map[string]any {
	return map[string]any{
		"id":       user.ID,
		"username": user.Username,
	}
}

func (user User) ToPublic() UserPublic {
	return UserPublic{
		ID:       user.ID,
		Username: user.Username,
	}
}

func (userList UserList) ToWeb() []map[string]any {
	m := make([]map[string]any, 0)
	for _, u := range userList {
		m = append(m, u.ToWeb())
	}
	return m
}
