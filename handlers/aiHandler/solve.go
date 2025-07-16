package aiHandler

import (
	"fmt"
	"quarto/models/ai"
	"quarto/models/game"

	"github.com/labstack/echo/v4"
)

type SolveRequest struct {
	History       []string   `json:"history"`
	SelectedPiece game.Piece `json:"selected_piece"`
	Depth         int        `json:"depth"`
}

type SolveResponse struct {
	BestMove       string   `json:"best_move"`
	Score          int      `json:"score"`
	SuggestedPiece int      `json:"suggested_piece"` // Pièce suggérée pour l'adversaire au coup suivant
	Continuation   []string `json:"continuation"`    // Liste des coups de la continuation
}

// solve handles the AI solve request for finding the best move in a Quarto game.
//
// @Summary Find the best move using AI
// @Description Analyzes the current game state and returns the optimal move using minimax algorithm
// @Tags AI
// @Accept json
// @Produce json
// @Param request body SolveRequest true "Solve request containing game history and search depth"
// @Success 200 {object} SolveResponse "Best move and evaluation score"
// @Failure 400 {object} map[string]string "Bad request - invalid format, depth, move history, or game state"
// @Router /ai/solve [post]
func solve(c echo.Context) error {

	var req SolveRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "Invalid request format")
	}

	// Validate depth
	if req.Depth < 1 || req.Depth > 16 {
		return echo.NewHTTPError(400, "Depth must be between 1 and 16")
	}

	// Parse the move history
	var moves []game.Move
	for _, moveStr := range req.History {
		piece, strPosition, err := game.ParseMoveNotation(moveStr)
		if err != nil {
			return echo.NewHTTPError(400, "Invalid move in history: "+moveStr)
		}

		row, col, err := game.PositionToCoords(strPosition)
		if err != nil {
			return echo.NewHTTPError(400, "Invalid position in move: "+moveStr)
		}

		fmt.Printf("Parsed move: Piece ID=%d, Position=%s (Row=%d, Col=%d)\n", piece, strPosition, row, col)

		moves = append(moves, game.Move{
			Piece:    piece,
			Position: game.Position{Row: row, Col: col},
		})
	}
	// Initialize AI engine with the specified depth
	depth := req.Depth
	if len(moves) <= 6 {
		depth = 5
	}
	engine := ai.NewEngine(depth)

	state := ai.ConvertHistoryToGameState(moves)
	found := false
	newAvailablePieces := make([]game.Piece, 0, len(state.AvailablePieces))
	for _, piece := range state.AvailablePieces {
		if piece == req.SelectedPiece {
			found = true
			continue
		}
		newAvailablePieces = append(newAvailablePieces, piece)
	}
	if !found {
		return echo.NewHTTPError(400, "Selected piece not found in available pieces")
	}
	state.SelectedPiece = req.SelectedPiece
	state.AvailablePieces = newAvailablePieces

	fmt.Printf("State: AvailablePieces=%v, SelectedPiece=%v, IsGameOver=%t, Winner=%d\n",
		state.AvailablePieces, state.SelectedPiece, state.IsGameOver, state.Winner)

	result := engine.Search(state)

	fmt.Printf("Best move found: Score=%d (%v), Depth=%d, Move=%d (%d) on %d,%d give %d\n", result.Score, len(state.AvailablePieces)%2 == 0, result.Depth, result.BestMoves[0].Move.Piece, req.SelectedPiece, result.BestMoves[0].Move.Position.Row, result.BestMoves[0].Move.Position.Col, result.BestMoves[0].SelectedPiece)

	bestMoveNotation := game.CreateMoveNotation(result.BestMoves[0].Move.Piece, game.CoordsToPosition(result.BestMoves[0].Move.Position.Row, result.BestMoves[0].Move.Position.Col))

	var strContinuation []string
	for _, move := range result.BestMoves {
		notation := game.CreateMoveNotation(move.Move.Piece, game.CoordsToPosition(move.Move.Position.Row, move.Move.Position.Col))
		strContinuation = append(strContinuation, notation)
	}

	return c.JSON(200, SolveResponse{
		BestMove:       bestMoveNotation,
		Score:          result.Score,
		SuggestedPiece: int(result.BestMoves[0].SelectedPiece),
		Continuation:   strContinuation,
	})

}
