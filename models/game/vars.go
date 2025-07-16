package game

import (
	"time"

	"github.com/fatih/structs"
)

// GamePhase constants
const (
	GamePhaseSelectPiece = iota
	GamePhasePlacePiece
)

// Piece constants
const (
	PieceWhiteSquareLargeFilled Piece = iota
	PieceWhiteSquareLargeEmpty
	PieceWhiteSquareSmallFilled
	PieceWhiteSquareSmallEmpty
	PieceWhiteCircleLargeFilled
	PieceWhiteCircleLargeEmpty
	PieceWhiteCircleSmallFilled
	PieceWhiteCircleSmallEmpty
	PieceBlackSquareLargeFilled
	PieceBlackSquareLargeEmpty
	PieceBlackSquareSmallFilled
	PieceBlackSquareSmallEmpty
	PieceBlackCircleLargeFilled
	PieceBlackCircleLargeEmpty
	PieceBlackCircleSmallFilled
	PieceBlackCircleSmallEmpty

	PieceEmpty Piece = -1
)

// Status constants
const (
	StatusPlaying = iota
	StatusFinished
)

type (
	Game struct {
		ID              string      `structs:"id" json:"id"`
		Player1ID       int64       `structs:"player1_id" json:"player1_id"`
		Player2ID       int64       `structs:"player2_id" json:"player2_id"`
		CurrentTurn     int64       `structs:"current_turn" json:"current_turn"`         // ID of the player whose turn it is
		GamePhase       int         `structs:"game_phase" json:"game_phase"`             // 0 = "selectPiece", 1 = "placePiece"
		Board           [4][4]Piece `structs:"board" json:"board"`                       // 4x4 matrix of Piece
		AvailablePieces []Piece     `structs:"available_pieces" json:"available_pieces"` // List of available pieces (1-16)
		SelectedPiece   Piece       `structs:"selected_piece" json:"selected_piece"`     // Current piece to place
		Status          int         `structs:"status" json:"status"`                     // 0 = "playing", 1 = "finished"
		Winner          int64       `structs:"winner" json:"winner"`                     // ID of the winner (0 if draw)
		History         []Move      `structs:"move_history" json:"move_history"`         // List of moves made in the game
		CreatedAt       time.Time   `structs:"created_at" json:"created_at"`
		UpdatedAt       time.Time   `structs:"updated_at" json:"updated_at"`
	}

	Piece int

	// Position représente une position sur le plateau Quarto (4x4)
	Position struct {
		Row int
		Col int
	}

	// Move représente un mouvement complet dans Quarto (placement + sélection pour l'adversaire)
	Move struct {
		Piece    Piece    // ID de la pièce sélectionnée par l'adversaire (0-15)
		Position Position // Position où placer la pièce sélectionnée
	}

	// Request/Response types
	SelectPieceRequest struct {
		PieceID Piece `json:"piece_id" validate:"required"`
	}

	PlacePieceRequest struct {
		Position string `json:"position" validate:"required"`
	}

	AIResponse struct {
		SelectedPiece  *int    `json:"selected_piece"`
		Position       *string `json:"position"`
		SuggestedPiece int     `json:"suggested_piece"` // Pièce suggérée pour l'adversaire
		Confidence     float64 `json:"confidence"`
		Evaluation     float64 `json:"evaluation"`
		CacheHits      int     `json:"cache_hits"`
		CacheMisses    int     `json:"cache_misses"`
		NodesVisited   int     `json:"nodes_visited"`
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
	for i := range 16 {
		pieces = append(pieces, Piece(i))
	}
	return pieces
}
