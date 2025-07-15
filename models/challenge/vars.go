package challenge

import (
	"time"

	"github.com/fatih/structs"
)

type Challenge struct {
	ID           string    `json:"id" structs:"id"`
	ChallengerID int64     `json:"challenger_id" structs:"challenger_id"`
	ChallengedID int64     `json:"challenged_id" structs:"challenged_id"`
	Status       string    `json:"status" structs:"status"` // pending, accepted, declined, expired, cancelled
	Message      string    `json:"message" structs:"message"`
	GameID       string    `json:"game_id,omitempty" structs:"game_id,omitempty"`
	CreatedAt    time.Time `json:"created_at" structs:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" structs:"updated_at"`
	ExpiresAt    time.Time `json:"expires_at" structs:"expires_at"`
	RespondedAt  time.Time `json:"responded_at,omitempty" structs:"responded_at,omitempty"`
}

// Structures pour les requêtes API
type SendChallengeRequest struct {
	ChallengedID int64  `json:"challenged_id" validate:"required"`
	Message      string `json:"message"`
}

type RespondToChallengeRequest struct {
	ChallengeID string `json:"challenge_id" validate:"required"`
	Accept      bool   `json:"accept"`
}

type ChallengeResponse struct {
	Challenge *Challenge `json:"challenge"`
	Game      any        `json:"game,omitempty"`
}

type ChallengeListResponse struct {
	Sent     []Challenge `json:"sent"`
	Received []Challenge `json:"received"`
}

// ToWeb convertit le Challenge en map pour l'API
func (c Challenge) ToWeb() map[string]any {
	return structs.Map(c)
}

// ToWebList convertit une liste de challenges en format web
func ToWebList(challenges []Challenge) []map[string]any {
	result := make([]map[string]any, len(challenges))
	for i, challenge := range challenges {
		result[i] = challenge.ToWeb()
	}
	return result
}

// IsExpired vérifie si le défi a expiré
func (c Challenge) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// CanRespond vérifie si le défi peut recevoir une réponse
func (c Challenge) CanRespond() bool {
	return c.Status == "pending" && !c.IsExpired()
}
