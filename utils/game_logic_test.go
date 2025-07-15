package utils

import (
	"quarto/models/game"
	"testing"
)

func TestPieceToNotation(t *testing.T) {
	tests := []struct {
		piece    game.Piece
		expected string
	}{
		{
			piece:    game.Piece{Color: "blanc", Shape: "carré", Size: "grand", Fill: "plein"},
			expected: "BCGP",
		},
		{
			piece:    game.Piece{Color: "noir", Shape: "rond", Size: "petit", Fill: "troué"},
			expected: "NRPT",
		},
		{
			piece:    game.Piece{Color: "blanc", Shape: "rond", Size: "petit", Fill: "plein"},
			expected: "BRPP",
		},
	}

	for _, test := range tests {
		result := PieceToNotation(test.piece)
		if result != test.expected {
			t.Errorf("PieceToNotation(%+v) = %s; want %s", test.piece, result, test.expected)
		}
	}
}

func TestCoordsToPosition(t *testing.T) {
	tests := []struct {
		row, col int
		expected string
	}{
		{0, 0, "a4"},
		{3, 3, "d1"},
		{1, 2, "c3"},
		{2, 1, "b2"},
	}

	for _, test := range tests {
		result := CoordsToPosition(test.row, test.col)
		if result != test.expected {
			t.Errorf("CoordsToPosition(%d, %d) = %s; want %s", test.row, test.col, result, test.expected)
		}
	}
}

func TestPositionToCoords(t *testing.T) {
	tests := []struct {
		position    string
		expectedRow int
		expectedCol int
		expectError bool
	}{
		{"a4", 0, 0, false},
		{"d1", 3, 3, false},
		{"c3", 1, 2, false},
		{"b2", 2, 1, false},
		{"e1", 0, 0, true}, // Position invalide
		{"a5", 0, 0, true}, // Position invalide
	}

	for _, test := range tests {
		row, col, err := PositionToCoords(test.position)
		if test.expectError {
			if err == nil {
				t.Errorf("PositionToCoords(%s) expected error but got none", test.position)
			}
		} else {
			if err != nil {
				t.Errorf("PositionToCoords(%s) unexpected error: %v", test.position, err)
			}
			if row != test.expectedRow || col != test.expectedCol {
				t.Errorf("PositionToCoords(%s) = (%d, %d); want (%d, %d)", test.position, row, col, test.expectedRow, test.expectedCol)
			}
		}
	}
}

func TestCheckWin(t *testing.T) {
	pieces := game.GetAllPieces()

	// Test avec une ligne gagnante (toutes de même couleur - blanc)
	board := [4][4]*int{}
	// Placer 4 pièces blanches en ligne
	p0, p1, p2, p3 := 0, 1, 2, 3 // IDs des 4 premières pièces blanches
	board[0][0] = &p0
	board[0][1] = &p1
	board[0][2] = &p2
	board[0][3] = &p3

	if !CheckWin(board, pieces) {
		t.Error("CheckWin should return true for winning line of same color pieces")
	}

	// Test sans victoire
	board2 := [4][4]*int{}
	p4, p5, p6, p7 := 4, 5, 6, 7 // Mélange de couleurs
	board2[0][0] = &p4
	board2[0][1] = &p5
	board2[0][2] = &p6
	board2[1][0] = &p7 // Pas en ligne

	if CheckWin(board2, pieces) {
		t.Error("CheckWin should return false for non-winning configuration")
	}
}

func TestInitializeGame(t *testing.T) {
	game := InitializeGame(1, 2)

	if game.Player1ID != 1 || game.Player2ID != 2 {
		t.Error("InitializeGame should set correct player IDs")
	}

	if game.CurrentTurn != "player1" {
		t.Error("InitializeGame should start with player1")
	}

	if game.GamePhase != "selectPiece" {
		t.Error("InitializeGame should start with selectPiece phase")
	}

	if game.Status != "playing" {
		t.Error("InitializeGame should start with playing status")
	}

	// Vérifier que toutes les pièces sont disponibles
	if !contains(game.AvailablePieces, "0") || !contains(game.AvailablePieces, "15") {
		t.Error("InitializeGame should have all pieces available")
	}
}

func TestBoardToMatrix(t *testing.T) {
	boardJSON := "[[0,1,null,2],[null,null,3,4],[5,null,null,null],[null,null,null,null]]"

	board, err := BoardToMatrix(boardJSON)
	if err != nil {
		t.Errorf("BoardToMatrix failed: %v", err)
	}

	// Vérifier quelques valeurs
	if board[0][0] == nil || *board[0][0] != 0 {
		t.Error("BoardToMatrix should parse first cell correctly")
	}

	if board[0][2] != nil {
		t.Error("BoardToMatrix should parse null cells correctly")
	}

	if board[1][3] == nil || *board[1][3] != 4 {
		t.Error("BoardToMatrix should parse numbered cells correctly")
	}
}

func TestMatrixToBoard(t *testing.T) {
	board := [4][4]*int{}
	val0, val1 := 0, 1
	board[0][0] = &val0
	board[0][1] = &val1
	// Le reste reste null

	result := MatrixToBoard(board)
	expected := "[[0,1,null,null],[null,null,null,null],[null,null,null,null],[null,null,null,null]]"

	if result != expected {
		t.Errorf("MatrixToBoard() = %s; want %s", result, expected)
	}
}

func TestCreateMoveNotation(t *testing.T) {
	piece := game.Piece{Color: "blanc", Shape: "carré", Size: "grand", Fill: "plein"}
	position := "c3"

	result := CreateMoveNotation(piece, position)
	expected := "BCGP-c3"

	if result != expected {
		t.Errorf("CreateMoveNotation(%+v, %s) = %s; want %s", piece, position, result, expected)
	}
}

// Fonction utilitaire pour les tests
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || (len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
