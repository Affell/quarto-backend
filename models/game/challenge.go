package game

import (
	"time"
)

// Challenge représente un défi entre deux joueurs
type Challenge struct {
	ID           string    `json:"id" db:"id"`
	ChallengerID int64     `json:"challenger_id" db:"challenger_id"`
	ChallengedID int64     `json:"challenged_id" db:"challenged_id"`
	Status       string    `json:"status" db:"status"` // pending, accepted, declined, expired
	GameID       *string   `json:"game_id,omitempty" db:"game_id"`
	Message      string    `json:"message" db:"message"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
}

// Structures pour les requêtes API
type SendChallengeRequest struct {
	ChallengedUserID int64  `json:"challenged_user_id"`
	Message          string `json:"message,omitempty"`
}

type ChallengeResponse struct {
	Action string `json:"action"` // accept, decline
}

// ToWeb convertit une Challenge en format web
func (c *Challenge) ToWeb() map[string]any {
	result := map[string]any{
		"id":            c.ID,
		"challenger_id": c.ChallengerID,
		"challenged_id": c.ChallengedID,
		"status":        c.Status,
		"message":       c.Message,
		"created_at":    c.CreatedAt,
		"updated_at":    c.UpdatedAt,
		"expires_at":    c.ExpiresAt,
	}

	if c.GameID != nil {
		result["game_id"] = *c.GameID
	}

	return result
}
