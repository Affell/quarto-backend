package challenge

import (
	"database/sql"
	"fmt"
	"quarto/models/postgresql"
	"time"

	"github.com/jackc/pgx/v4"
)

// ScanChallenge scanne une ligne de résultat SQL en structure Challenge
func ScanChallenge(row pgx.Row) (c Challenge, err error) {
	var (
		id                                           sql.NullString
		challengerID, challengedID                   sql.NullInt64
		status, message                              sql.NullString
		gameID                                       sql.NullString
		createdAt, updatedAt, expiresAt, respondedAt sql.NullTime
	)

	err = row.Scan(
		&id,
		&challengerID,
		&challengedID,
		&status,
		&message,
		&gameID,
		&createdAt,
		&updatedAt,
		&expiresAt,
		&respondedAt,
	)

	if err != nil {
		return
	}

	c = Challenge{
		ID:           id.String,
		ChallengerID: challengerID.Int64,
		ChallengedID: challengedID.Int64,
		Status:       status.String,
		Message:      message.String,
		GameID:       gameID.String,
		CreatedAt:    createdAt.Time,
		UpdatedAt:    updatedAt.Time,
		ExpiresAt:    expiresAt.Time,
		RespondedAt:  respondedAt.Time,
	}

	return
}

// CreateChallenge insère un nouveau défi en base
func CreateChallenge(challenge Challenge) error {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		INSERT INTO challenges (id, challenger_id, challenged_id, status, message, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = sqlCo.Exec(postgresql.SQLCtx, query,
		challenge.ID, challenge.ChallengerID, challenge.ChallengedID,
		challenge.Status, challenge.Message, challenge.CreatedAt,
		challenge.UpdatedAt, challenge.ExpiresAt)

	return err
}

// GetChallengeByID récupère un défi par son ID
func GetChallengeByID(challengeID string) (*Challenge, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		SELECT id, challenger_id, challenged_id, status, message, game_id,
			created_at, updated_at, expires_at, responded_at
		FROM challenges WHERE id = $1`

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, challengeID)
	challenge, err := ScanChallenge(row)
	if err != nil {
		return nil, err
	}

	return &challenge, nil
}

// GetPendingChallengeBetween récupère un défi en attente entre deux joueurs
func GetPendingChallengeBetween(challengerID, challengedID int64) (*Challenge, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		SELECT id, challenger_id, challenged_id, status, message, game_id,
			created_at, updated_at, expires_at, responded_at
		FROM challenges 
		WHERE ((challenger_id = $1 AND challenged_id = $2) OR (challenger_id = $2 AND challenged_id = $1))
		AND status = 'pending'
		AND expires_at > NOW()`

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, challengerID, challengedID)
	challenge, err := ScanChallenge(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &challenge, nil
}

// GetUserChallenges récupère tous les défis d'un utilisateur (envoyés et reçus)
func GetUserChallenges(userID int64) ([]Challenge, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		SELECT id, challenger_id, challenged_id, status, message, game_id,
			created_at, updated_at, expires_at, responded_at
		FROM challenges 
		WHERE challenger_id = $1 OR challenged_id = $1
		ORDER BY created_at DESC`

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []Challenge
	for rows.Next() {
		challenge, err := ScanChallenge(rows)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, challenge)
	}

	return challenges, nil
}

// UpdateChallengeStatus met à jour le statut d'un défi
func UpdateChallengeStatus(challengeID, status string, gameID *string) error {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query string
	var args []interface{}

	if gameID != nil {
		query = `
			UPDATE challenges 
			SET status = $1, game_id = $2, responded_at = $3, updated_at = $4
			WHERE id = $5`
		args = []interface{}{status, *gameID, time.Now(), time.Now(), challengeID}
	} else {
		query = `
			UPDATE challenges 
			SET status = $1, responded_at = $2, updated_at = $3
			WHERE id = $4`
		args = []interface{}{status, time.Now(), time.Now(), challengeID}
	}

	_, err = sqlCo.Exec(postgresql.SQLCtx, query, args...)
	return err
}

// CleanupExpiredChallenges supprime ou marque comme expirés les défis anciens
func CleanupExpiredChallenges() error {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		UPDATE challenges 
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'pending' AND expires_at <= NOW()`

	_, err = sqlCo.Exec(postgresql.SQLCtx, query)
	return err
}
