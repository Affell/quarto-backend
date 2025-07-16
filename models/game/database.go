package game

import (
	"database/sql"
	"fmt"
	"quarto/models/postgresql"
	"time"

	"github.com/jackc/pgx/v4"
)

// ScanGame scanne une ligne de résultat SQL en structure Game
func ScanGame(row pgx.Row) (g Game, err error) {
	var (
		id                                             sql.NullString
		player1ID, player2ID                           sql.NullInt64
		currentTurn, gamePhase, board, availablePieces sql.NullString
		selectedPiece                                  sql.NullInt32
		status                                         sql.NullString
		winner                                         sql.NullString
		moveHistory                                    sql.NullString
		createdAt, updatedAt                           sql.NullTime
	)

	err = row.Scan(
		&id,
		&player1ID,
		&player2ID,
		&currentTurn,
		&gamePhase,
		&board,
		&availablePieces,
		&selectedPiece,
		&status,
		&winner,
		&moveHistory,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return
	}

	// var selectedPiecePtr *int
	// if selectedPiece.Valid {
	// 	selectedPieceInt := int(selectedPiece.Int32)
	// 	selectedPiecePtr = &selectedPieceInt
	// }

	// var winnerPtr *string
	// if winner.Valid {
	// 	winnerPtr = &winner.String
	// }

	g = Game{
		ID:        id.String,
		Player1ID: player1ID.Int64,
		Player2ID: player2ID.Int64,
		// CurrentTurn:     currentTurn.String,
		// GamePhase:       gamePhase.String,
		// Board:           board.String,
		// AvailablePieces: availablePieces.String,
		// SelectedPiece:   selectedPiecePtr,
		// Status:          status.String,
		// Winner:          winnerPtr,
		// History:         moveHistory.String,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}

	return
}

// CreateGame insère une nouvelle partie en base
func CreateGame(game Game) error {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		INSERT INTO games (id, player1_id, player2_id, current_turn, game_phase, 
			board, available_pieces, selected_piece, status, winner, move_history, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = sqlCo.Exec(postgresql.SQLCtx, query,
		game.ID, game.Player1ID, game.Player2ID,
		game.CurrentTurn, game.GamePhase, game.Board, game.AvailablePieces,
		game.SelectedPiece, game.Status, game.Winner, game.History,
		game.CreatedAt, game.UpdatedAt)

	return err
}

// GetGameByID récupère une partie par son ID
func GetGameByID(gameID string) (g Game, err error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		err = fmt.Errorf("erreur de connexion DB: %v", err)
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		SELECT id, player1_id, player2_id, current_turn, game_phase,
			board, available_pieces, selected_piece, status, winner, move_history,
			created_at, updated_at
		FROM games WHERE id = $1`

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, gameID)
	g, err = ScanGame(row)
	if err != nil {
		err = fmt.Errorf("partie non trouvée: %v", err)
		return
	}

	return
}

// UpdateGame met à jour une partie en base
func UpdateGame(game Game) error {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		UPDATE games 
		SET current_turn = $1, game_phase = $2, board = $3, available_pieces = $4,
			selected_piece = $5, status = $6, winner = $7, move_history = $8, updated_at = $9
		WHERE id = $10`

	_, err = sqlCo.Exec(postgresql.SQLCtx, query,
		game.CurrentTurn, game.GamePhase, game.Board, game.AvailablePieces,
		game.SelectedPiece, game.Status, game.Winner, game.History,
		time.Now(), game.ID)

	return err
}

// GetUserGames récupère toutes les parties d'un utilisateur
func GetUserGames(userID int64) ([]Game, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		SELECT id, player1_id, player2_id, current_turn, game_phase,
			board, available_pieces, selected_piece, status, winner, move_history,
			created_at, updated_at
		FROM games 
		WHERE player1_id = $1 OR player2_id = $1
		ORDER BY created_at DESC`

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		game, err := ScanGame(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

// GetActiveGames récupère toutes les parties actives d'un utilisateur
func GetActiveGames(userID int64) ([]Game, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `
		SELECT id, player1_id, player2_id, current_turn, game_phase,
			board, available_pieces, selected_piece, status, winner, move_history,
			created_at, updated_at
		FROM games 
		WHERE (player1_id = $1 OR player2_id = $1) AND status = 'active'
		ORDER BY updated_at DESC`

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		game, err := ScanGame(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

// DeleteGame supprime une partie (pour les tests ou le nettoyage)
func DeleteGame(gameID string) error {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return fmt.Errorf("erreur de connexion DB: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := `DELETE FROM games WHERE id = $1`
	_, err = sqlCo.Exec(postgresql.SQLCtx, query, gameID)
	return err
}
