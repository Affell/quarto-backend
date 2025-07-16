package ai

import (
	"quarto/models/game"
	"testing"
)

func TestEvaluation(t *testing.T) {

	engine := NewEngine(3)

	// Test case 1: Horizontal win on first row
	state1 := GameState{
		Board: [4][4]game.Piece{
			{0, 2, 13, game.PieceEmpty},
			{5, 1, 14, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
			{4, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
		},
		AvailablePieces: []game.Piece{3, 7, 8, 9, 10, 11, 12, 15},
		SelectedPiece:   6,
	}

	state1 = state1.ApplyMove(AIMove{
		Move: game.Move{
			Piece:    6,
			Position: game.Position{Row: 2, Col: 0},
		},
	})

	score := engine.evaluatePosition(state1)
	if score != 10000 {
		t.Errorf("Test 1 : Expected score 10000 for horizontal win, got %d", score)
	}

	// Test case 2: Horizontal win on first row
	state2 := GameState{
		Board: [4][4]game.Piece{
			{0, 2, 12, game.PieceEmpty},
			{game.PieceEmpty, 1, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, 3, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
		},
		AvailablePieces: []game.Piece{5, 6, 7, 8, 9, 10, 11, 13, 14, 15},
		SelectedPiece:   4,
	}

	state2 = state2.ApplyMove(AIMove{
		Move: game.Move{
			Piece:    4,
			Position: game.Position{Row: 3, Col: 3},
		},
	})

	score = engine.evaluatePosition(state2)
	if score != 10000 {
		t.Errorf("Test 2 :Expected score 10000 for horizontal win, got %d", score)
	}

	// Test case 3: Diagonal win
	state3 := GameState{
		Board: [4][4]game.Piece{
			{0, 2, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, 1, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, 4, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
		},
		AvailablePieces: []game.Piece{3, 5, 6, 7, 8, 9, 10, 11, 13, 14, 15},
		SelectedPiece:   12,
	}

	state3 = state3.ApplyMove(AIMove{
		Move: game.Move{
			Piece:    12,
			Position: game.Position{Row: 3, Col: 3},
		},
	})

	score = engine.evaluatePosition(state3)
	if score != -10000 {
		t.Errorf("Test 3: Expected score -10000 for diagonal win, got %d", score)
	}

	// Test case 4: Anti-diagonal win
	state4 := GameState{
		Board: [4][4]game.Piece{
			{0, 2, game.PieceEmpty, 12},
			{game.PieceEmpty, 1, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, 4, game.PieceEmpty},
			{9, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
		},
		AvailablePieces: []game.Piece{3, 5, 6, 7, 8, 10, 11, 13, 14},
		SelectedPiece:   15,
	}

	state4 = state4.ApplyMove(AIMove{
		Move: game.Move{
			Piece:    15,
			Position: game.Position{Row: 0, Col: 2},
		},
	})

	score = engine.evaluatePosition(state4)
	if score != 0 {
		t.Errorf("Test 4: Expected score 0, got %d", score)
	}

	// Test case 5: No win condition
	state5 := GameState{
		Board: [4][4]game.Piece{
			{0, 2, 12, game.PieceEmpty},
			{11, 1, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, 4, game.PieceEmpty},
			{15, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
		},
		AvailablePieces: []game.Piece{5, 6, 7, 8, 9, 10, 13, 14},
		SelectedPiece:   3,
	}

	score = engine.evaluatePosition(state5)
	if score != 0 {
		t.Errorf("Test 5: Expected score 0 for no win condition, got %d", score)
	}

}

func TestEngineSearch(t *testing.T) {
	engine := NewEngine(2)

	// Test case 1
	state1 := GameState{
		Board: [4][4]game.Piece{
			{game.PieceEmpty, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
			{game.PieceEmpty, game.PieceEmpty, game.PieceEmpty, game.PieceEmpty},
		},
		AvailablePieces: []game.Piece{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		SelectedPiece:   0,
	}

	newState := state1.ApplyMove(AIMove{
		Move: game.Move{
			Piece:    3,
			Position: game.Position{Row: 0, Col: 0},
		},
	})

	t.Errorf("New state hash: %s", engine.hashGameState(newState))
	result := engine.Search(state1)
	t.Errorf("result : %v", result)
	var strContinuation []string
	for _, move := range result.BestMoves {
		notation := game.CreateMoveNotation(move.Move.Piece, game.CoordsToPosition(move.Move.Position.Row, move.Move.Position.Col))
		strContinuation = append(strContinuation, notation)
	}
	t.Logf("Best move found: Score=%d, Depth=%d, Move=%s on %s give %s ; Continuation=%v\n",
		result.Score, result.Depth, game.PieceToNotation(result.BestMoves[0].Move.Piece), game.CoordsToPosition(result.BestMoves[0].Move.Position.Row, result.BestMoves[0].Move.Position.Col), game.PieceToNotation(result.BestMoves[0].SelectedPiece), strContinuation)
}
