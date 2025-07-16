package game

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"quarto/models/postgresql"
	"time"

	"github.com/jackc/pgx/v4"
)

// serializeBoardToJSON convertit le plateau [4][4]Piece en JSON
func serializeBoardToJSON(board [4][4]Piece) (string, error) {
	// Convertir le plateau en slice de slices d'int pour JSON
	jsonBoard := make([][]interface{}, 4)
	for i := 0; i < 4; i++ {
		jsonBoard[i] = make([]interface{}, 4)
		for j := 0; j < 4; j++ {
			if board[i][j] == PieceEmpty {
				jsonBoard[i][j] = nil
			} else {
				jsonBoard[i][j] = int(board[i][j])
			}
		}
	}

	jsonData, err := json.Marshal(jsonBoard)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// deserializeBoardFromJSON convertit le JSON en plateau [4][4]Piece
func deserializeBoardFromJSON(jsonStr string) ([4][4]Piece, error) {
	var board [4][4]Piece

	if jsonStr == "" || jsonStr == "null" {
		return GetEmptyBoard(), nil
	}

	var jsonBoard [][]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonBoard); err != nil {
		return board, fmt.Errorf("erreur de parsing du plateau: %v", err)
	}

	if len(jsonBoard) != 4 {
		return board, fmt.Errorf("plateau doit avoir 4 lignes, trouvé %d", len(jsonBoard))
	}

	for i := 0; i < 4; i++ {
		if len(jsonBoard[i]) != 4 {
			return board, fmt.Errorf("ligne %d doit avoir 4 colonnes, trouvé %d", i, len(jsonBoard[i]))
		}
		for j := 0; j < 4; j++ {
			if jsonBoard[i][j] == nil {
				board[i][j] = PieceEmpty
			} else {
				if val, ok := jsonBoard[i][j].(float64); ok {
					board[i][j] = Piece(int(val))
				} else {
					return board, fmt.Errorf("valeur invalide dans le plateau à [%d][%d]: %v", i, j, jsonBoard[i][j])
				}
			}
		}
	}

	return board, nil
}

// serializeAvailablePiecesToJSON convertit []Piece en JSON
func serializeAvailablePiecesToJSON(pieces []Piece) (string, error) {
	jsonPieces := make([]int, len(pieces))
	for i, piece := range pieces {
		jsonPieces[i] = int(piece)
	}

	jsonData, err := json.Marshal(jsonPieces)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// deserializeAvailablePiecesFromJSON convertit le JSON en []Piece
func deserializeAvailablePiecesFromJSON(jsonStr string) ([]Piece, error) {
	if jsonStr == "" || jsonStr == "null" {
		return GetAllPieces(), nil
	}

	var jsonPieces []int
	if err := json.Unmarshal([]byte(jsonStr), &jsonPieces); err != nil {
		return nil, fmt.Errorf("erreur de parsing des pièces disponibles: %v", err)
	}

	pieces := make([]Piece, len(jsonPieces))
	for i, piece := range jsonPieces {
		pieces[i] = Piece(piece)
	}

	return pieces, nil
}

// serializeHistoryToJSON convertit []Move en JSON
func serializeHistoryToJSON(history []Move) (string, error) {
	// Convertir en structure sérialisable
	type SerializableMove struct {
		Piece    int `json:"piece"`
		Position struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"position"`
	}

	jsonHistory := make([]SerializableMove, len(history))
	for i, move := range history {
		jsonHistory[i] = SerializableMove{
			Piece: int(move.Piece),
		}
		jsonHistory[i].Position.Row = move.Position.Row
		jsonHistory[i].Position.Col = move.Position.Col
	}

	jsonData, err := json.Marshal(jsonHistory)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// deserializeHistoryFromJSON convertit le JSON en []Move
func deserializeHistoryFromJSON(jsonStr string) ([]Move, error) {
	if jsonStr == "" || jsonStr == "null" {
		return []Move{}, nil
	}

	type SerializableMove struct {
		Piece    int `json:"piece"`
		Position struct {
			Row int `json:"row"`
			Col int `json:"col"`
		} `json:"position"`
	}

	var jsonHistory []SerializableMove
	if err := json.Unmarshal([]byte(jsonStr), &jsonHistory); err != nil {
		return nil, fmt.Errorf("erreur de parsing de l'historique: %v", err)
	}

	history := make([]Move, len(jsonHistory))
	for i, move := range jsonHistory {
		history[i] = Move{
			Piece: Piece(move.Piece),
			Position: Position{
				Row: move.Position.Row,
				Col: move.Position.Col,
			},
		}
	}

	return history, nil
}

// ScanGame scanne une ligne de résultat SQL en structure Game
func ScanGame(row pgx.Row) (g Game, err error) {
	var (
		id                                        sql.NullString
		player1ID, player2ID, currentTurn, winner sql.NullInt64
		gamePhase, status                         sql.NullInt32
		selectedPiece                             sql.NullInt32
		board, availablePieces, moveHistory       sql.NullString
		createdAt, updatedAt                      sql.NullTime
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

	// Désérialiser le plateau
	var gameBoard [4][4]Piece
	if board.Valid {
		gameBoard, err = deserializeBoardFromJSON(board.String)
		if err != nil {
			return
		}
	} else {
		gameBoard = GetEmptyBoard()
	}

	// Désérialiser les pièces disponibles
	var pieces []Piece
	if availablePieces.Valid {
		pieces, err = deserializeAvailablePiecesFromJSON(availablePieces.String)
		if err != nil {
			return
		}
	} else {
		pieces = GetAllPieces()
	}

	// Désérialiser l'historique
	var history []Move
	if moveHistory.Valid {
		history, err = deserializeHistoryFromJSON(moveHistory.String)
		if err != nil {
			return
		}
	} else {
		history = []Move{}
	}

	g = Game{
		ID:              id.String,
		Player1ID:       player1ID.Int64,
		Player2ID:       player2ID.Int64,
		CurrentTurn:     currentTurn.Int64,
		GamePhase:       int(gamePhase.Int32),
		Board:           gameBoard,
		AvailablePieces: pieces,
		SelectedPiece:   Piece(selectedPiece.Int32),
		Status:          int(status.Int32),
		Winner:          winner.Int64,
		History:         history,
		CreatedAt:       createdAt.Time,
		UpdatedAt:       updatedAt.Time,
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

	// Sérialiser les données complexes
	boardJSON, err := serializeBoardToJSON(game.Board)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation du plateau: %v", err)
	}

	availablePiecesJSON, err := serializeAvailablePiecesToJSON(game.AvailablePieces)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation des pièces disponibles: %v", err)
	}

	historyJSON, err := serializeHistoryToJSON(game.History)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation de l'historique: %v", err)
	}

	query := `
		INSERT INTO games (id, player1_id, player2_id, current_turn, game_phase, 
			board, available_pieces, selected_piece, status, winner, move_history, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = sqlCo.Exec(postgresql.SQLCtx, query,
		game.ID, game.Player1ID, game.Player2ID,
		game.CurrentTurn, game.GamePhase, boardJSON, availablePiecesJSON,
		int(game.SelectedPiece), game.Status, game.Winner, historyJSON,
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

	// Sérialiser les données complexes
	boardJSON, err := serializeBoardToJSON(game.Board)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation du plateau: %v", err)
	}

	availablePiecesJSON, err := serializeAvailablePiecesToJSON(game.AvailablePieces)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation des pièces disponibles: %v", err)
	}

	historyJSON, err := serializeHistoryToJSON(game.History)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation de l'historique: %v", err)
	}

	query := `
		UPDATE games 
		SET current_turn = $1, game_phase = $2, board = $3, available_pieces = $4,
			selected_piece = $5, status = $6, winner = $7, move_history = $8, updated_at = $9
		WHERE id = $10`

	_, err = sqlCo.Exec(postgresql.SQLCtx, query,
		game.CurrentTurn, game.GamePhase, boardJSON, availablePiecesJSON,
		int(game.SelectedPiece), game.Status, game.Winner, historyJSON,
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
		WHERE (player1_id = $1 OR player2_id = $1) AND status = 0
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
