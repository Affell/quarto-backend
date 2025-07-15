package game

import (
	"time"

	"github.com/fatih/structs"
)

type (
	Game struct {
		ID              string    `structs:"id" json:"id"`
		Player1ID       int64     `structs:"player1_id" json:"player1_id"`
		Player2ID       int64     `structs:"player2_id" json:"player2_id"`
		CurrentTurn     string    `structs:"current_turn" json:"current_turn"`         // "player1" | "player2"
		GamePhase       string    `structs:"game_phase" json:"game_phase"`             // "selectPiece" | "placePiece"
		Board           string    `structs:"board" json:"board"`                       // JSON serialized 4x4 board
		AvailablePieces string    `structs:"available_pieces" json:"available_pieces"` // JSON array of piece IDs
		SelectedPiece   *int      `structs:"selected_piece" json:"selected_piece"`     // Current piece ID to place
		Status          string    `structs:"status" json:"status"`                     // "active" | "finished"
		Winner          *string   `structs:"winner" json:"winner"`                     // "player1" | "player2" | "draw"
		MoveHistory     string    `structs:"move_history" json:"move_history"`         // JSON array of moves in notation
		CreatedAt       time.Time `structs:"created_at" json:"created_at"`
		UpdatedAt       time.Time `structs:"updated_at" json:"updated_at"`
	}

	Piece struct {
		ID    int    `structs:"id" json:"id"`
		Color string `structs:"color" json:"color"` // "blanc" | "noir"
		Shape string `structs:"shape" json:"shape"` // "carré" | "rond"
		Size  string `structs:"size" json:"size"`   // "grand" | "petit"
		Fill  string `structs:"fill" json:"fill"`   // "plein" | "troué"
	}

	// Request/Response types
	SelectPieceRequest struct {
		PieceID int `json:"piece_id" validate:"required"`
	}

	PlacePieceRequest struct {
		Position string `json:"position" validate:"required"`
	}

	AIResponse struct {
		SelectedPiece *int    `json:"selected_piece"`
		Position      *string `json:"position"`
		Confidence    float64 `json:"confidence"`
		Evaluation    float64 `json:"evaluation"`
	}

	AnalysisResponse struct {
		Evaluation float64  `json:"evaluation"`
		BestMoves  []string `json:"best_moves"`
		Threats    []string `json:"threats"`
	}

	GameList []Game
)

func (game Game) ToWeb() map[string]any {
	return structs.Map(game)
}

func (piece Piece) ToWeb() map[string]any {
	return structs.Map(piece)
}

// Génère toutes les 16 pièces du jeu Quarto
func GetAllPieces() []Piece {
	pieces := make([]Piece, 0, 16)
	colors := []string{"blanc", "noir"}
	shapes := []string{"carré", "rond"}
	sizes := []string{"grand", "petit"}
	fills := []string{"plein", "troué"}

	id := 0
	for _, color := range colors {
		for _, shape := range shapes {
			for _, size := range sizes {
				for _, fill := range fills {
					pieces = append(pieces, Piece{
						ID:    id,
						Color: color,
						Shape: shape,
						Size:  size,
						Fill:  fill,
					})
					id++
				}
			}
		}
	}
	return pieces
}
