package game

import (
	"testing"
)

func TestApplyMove(t *testing.T) {
	tests := []struct {
		name          string
		initialBoard  [4][4]Piece
		move          Move
		expectedBoard [4][4]Piece
		shouldChange  bool
	}{
		{
			name:         "Apply move to empty position",
			initialBoard: GetEmptyBoard(),
			move: Move{
				Piece:    PieceWhiteSquareLargeFilled,
				Position: Position{Row: 0, Col: 0},
			},
			expectedBoard: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			shouldChange: true,
		},
		{
			name: "Apply move to occupied position",
			initialBoard: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			move: Move{
				Piece:    PieceBlackCircleSmallEmpty,
				Position: Position{Row: 0, Col: 0},
			},
			expectedBoard: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			shouldChange: false,
		},
		{
			name:         "Apply move to center position",
			initialBoard: GetEmptyBoard(),
			move: Move{
				Piece:    PieceBlackCircleLargeFilled,
				Position: Position{Row: 2, Col: 2},
			},
			expectedBoard: [4][4]Piece{
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceBlackCircleLargeFilled, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			shouldChange: true,
		},
		{
			name:         "Apply move to bottom-right corner",
			initialBoard: GetEmptyBoard(),
			move: Move{
				Piece:    PieceWhiteCircleSmallEmpty,
				Position: Position{Row: 3, Col: 3},
			},
			expectedBoard: [4][4]Piece{
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceWhiteCircleSmallEmpty},
			},
			shouldChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyMove(tt.initialBoard, tt.move)

			// Vérifier que le plateau résultant correspond à l'attendu
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if result[i][j] != tt.expectedBoard[i][j] {
						t.Errorf("Position [%d][%d]: expected %d, got %d",
							i, j, tt.expectedBoard[i][j], result[i][j])
					}
				}
			}

			// Vérifier que le plateau original n'a pas été modifié
			originalChanged := false
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					if tt.initialBoard[i][j] != GetEmptyBoard()[i][j] &&
						tt.name == "Apply move to empty position" {
						originalChanged = true
						break
					}
				}
			}
			if originalChanged {
				t.Error("Original board was modified, but it should remain unchanged")
			}
		})
	}
}

func TestCanApplyMove(t *testing.T) {
	tests := []struct {
		name     string
		board    [4][4]Piece
		move     Move
		expected bool
	}{
		{
			name:  "Valid move on empty board",
			board: GetEmptyBoard(),
			move: Move{
				Piece:    PieceWhiteSquareLargeFilled,
				Position: Position{Row: 0, Col: 0},
			},
			expected: true,
		},
		{
			name: "Invalid move on occupied position",
			board: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			move: Move{
				Piece:    PieceBlackCircleSmallEmpty,
				Position: Position{Row: 0, Col: 0},
			},
			expected: false,
		},
		{
			name:  "Invalid move - row out of bounds (negative)",
			board: GetEmptyBoard(),
			move: Move{
				Piece:    PieceWhiteSquareLargeFilled,
				Position: Position{Row: -1, Col: 0},
			},
			expected: false,
		},
		{
			name:  "Invalid move - row out of bounds (too high)",
			board: GetEmptyBoard(),
			move: Move{
				Piece:    PieceWhiteSquareLargeFilled,
				Position: Position{Row: 4, Col: 0},
			},
			expected: false,
		},
		{
			name:  "Invalid move - column out of bounds (negative)",
			board: GetEmptyBoard(),
			move: Move{
				Piece:    PieceWhiteSquareLargeFilled,
				Position: Position{Row: 0, Col: -1},
			},
			expected: false,
		},
		{
			name:  "Invalid move - column out of bounds (too high)",
			board: GetEmptyBoard(),
			move: Move{
				Piece:    PieceWhiteSquareLargeFilled,
				Position: Position{Row: 0, Col: 4},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanApplyMove(tt.board, tt.move)
			if result != tt.expected {
				t.Errorf("CanApplyMove() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCheckWin(t *testing.T) {
	tests := []struct {
		name     string
		board    [4][4]Piece
		expected bool
		winType  string
	}{
		{
			name:     "No win - empty board",
			board:    GetEmptyBoard(),
			expected: false,
			winType:  "none",
		},
		{
			name: "Win - horizontal line (row 0) - all white pieces",
			board: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceWhiteCircleLargeFilled, PieceWhiteCircleLargeEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			expected: true,
			winType:  "horizontal",
		},
		{
			name: "Win - horizontal line (row 1) - all large pieces",
			board: [4][4]Piece{
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceBlackSquareLargeFilled, PieceBlackSquareLargeEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			expected: true,
			winType:  "horizontal",
		},
		{
			name: "Win - vertical line (col 0) - all square pieces",
			board: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceWhiteSquareLargeEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceBlackSquareLargeFilled, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceBlackSquareLargeEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			expected: true,
			winType:  "vertical",
		},
		{
			name: "Win - vertical line (col 2) - all filled pieces",
			board: [4][4]Piece{
				{PieceEmpty, PieceEmpty, PieceWhiteSquareLargeFilled, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceWhiteSquareSmallFilled, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceBlackSquareLargeFilled, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceBlackSquareSmallFilled, PieceEmpty},
			},
			expected: true,
			winType:  "vertical",
		},
		{
			name: "Win - main diagonal - all circle pieces",
			board: [4][4]Piece{
				{PieceWhiteCircleLargeFilled, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceWhiteCircleLargeEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceBlackCircleLargeFilled, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceBlackCircleLargeEmpty},
			},
			expected: true,
			winType:  "diagonal",
		},
		{
			name: "Win - anti-diagonal - all small pieces",
			board: [4][4]Piece{
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceWhiteSquareSmallFilled},
				{PieceEmpty, PieceEmpty, PieceWhiteSquareSmallEmpty, PieceEmpty},
				{PieceEmpty, PieceBlackCircleSmallFilled, PieceEmpty, PieceEmpty},
				{PieceBlackCircleSmallEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			expected: true,
			winType:  "anti-diagonal",
		},
		{
			name: "No win - different characteristics",
			board: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceBlackCircleSmallEmpty, PieceWhiteCircleSmallFilled, PieceBlackSquareLargeEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			expected: false,
			winType:  "none",
		},
		{
			name: "No win - incomplete line",
			board: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceWhiteCircleLargeFilled, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			expected: false,
			winType:  "none",
		},
		{
			name: "Win - horizontal line (row 3) - all empty pieces",
			board: [4][4]Piece{
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
				{PieceWhiteSquareLargeEmpty, PieceWhiteCircleLargeEmpty, PieceBlackSquareLargeEmpty, PieceBlackCircleLargeEmpty},
			},
			expected: true,
			winType:  "horizontal",
		},
		{
			name: "Win - complex board with multiple pieces but win on diagonal",
			board: [4][4]Piece{
				{PieceWhiteSquareLargeFilled, PieceBlackCircleSmallEmpty, PieceWhiteCircleSmallFilled, PieceBlackSquareLargeEmpty},
				{PieceBlackSquareSmallFilled, PieceWhiteSquareSmallFilled, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceWhiteSquareSmallEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceWhiteSquareSmallEmpty},
			},
			expected: true,
			winType:  "diagonal",
		},
		{
			name: "No win - mixed characteristics in each line",
			board: [4][4]Piece{
				{0, 2, 12, 5},
				{3, 1, PieceEmpty, PieceEmpty},
				{PieceEmpty, PieceEmpty, 4, PieceEmpty},
				{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
			},
			expected: false,
			winType:  "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckWin(tt.board)
			if result != tt.expected {
				t.Errorf("CheckWin() = %v, expected %v for %s win", result, tt.expected, tt.winType)

				// Debug: print the board for failed tests
				if result != tt.expected {
					t.Logf("Board state:")
					for i := 0; i < 4; i++ {
						t.Logf("Row %d: %v", i, tt.board[i])
					}
				}
			}
		})
	}
}

func TestHasCommonCharacteristic(t *testing.T) {
	tests := []struct {
		name     string
		pieces   []Piece
		expected bool
		reason   string
	}{
		{
			name:     "All white pieces",
			pieces:   []Piece{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceWhiteCircleLargeFilled, PieceWhiteCircleLargeEmpty},
			expected: true,
			reason:   "same color (white)",
		},
		{
			name:     "All black pieces",
			pieces:   []Piece{PieceBlackSquareLargeFilled, PieceBlackSquareLargeEmpty, PieceBlackCircleLargeFilled, PieceBlackCircleLargeEmpty},
			expected: true,
			reason:   "same color (black)",
		},
		{
			name:     "All large pieces",
			pieces:   []Piece{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceBlackSquareLargeFilled, PieceBlackSquareLargeEmpty},
			expected: true,
			reason:   "same size (large)",
		},
		{
			name:     "All small pieces",
			pieces:   []Piece{PieceWhiteSquareSmallFilled, PieceWhiteSquareSmallEmpty, PieceBlackSquareSmallFilled, PieceBlackSquareSmallEmpty},
			expected: true,
			reason:   "same size (small)",
		},
		{
			name:     "All square pieces",
			pieces:   []Piece{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceBlackSquareLargeFilled, PieceBlackSquareLargeEmpty},
			expected: true,
			reason:   "same shape (square)",
		},
		{
			name:     "All circle pieces",
			pieces:   []Piece{PieceWhiteCircleLargeFilled, PieceWhiteCircleLargeEmpty, PieceBlackCircleLargeFilled, PieceBlackCircleLargeEmpty},
			expected: true,
			reason:   "same shape (circle)",
		},
		{
			name:     "All filled pieces",
			pieces:   []Piece{PieceWhiteSquareLargeFilled, PieceWhiteSquareSmallFilled, PieceBlackSquareLargeFilled, PieceBlackSquareSmallFilled},
			expected: true,
			reason:   "same fill (filled)",
		},
		{
			name:     "All empty pieces",
			pieces:   []Piece{PieceWhiteSquareLargeEmpty, PieceWhiteSquareSmallEmpty, PieceBlackSquareLargeEmpty, PieceBlackSquareSmallEmpty},
			expected: true,
			reason:   "same fill (empty)",
		},
		{
			name:     "No common characteristics",
			pieces:   []Piece{PieceWhiteSquareLargeFilled, PieceBlackCircleSmallEmpty, PieceWhiteCircleSmallFilled, PieceBlackSquareLargeEmpty},
			expected: false,
			reason:   "no common characteristics",
		},
		{
			name:     "Wrong number of pieces (too few)",
			pieces:   []Piece{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceWhiteCircleLargeFilled},
			expected: false,
			reason:   "only 3 pieces",
		},
		{
			name:     "Wrong number of pieces (too many)",
			pieces:   []Piece{PieceWhiteSquareLargeFilled, PieceWhiteSquareLargeEmpty, PieceWhiteCircleLargeFilled, PieceWhiteCircleLargeEmpty, PieceBlackSquareLargeFilled},
			expected: false,
			reason:   "5 pieces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasCommonCharacteristic(tt.pieces)
			if result != tt.expected {
				t.Errorf("hasCommonCharacteristic() = %v, expected %v (%s)", result, tt.expected, tt.reason)

				// Debug: print piece details for failed tests
				if result != tt.expected && len(tt.pieces) == 4 {
					t.Logf("Piece details:")
					for i, piece := range tt.pieces {
						color, shape, size, fill := GetPieceCharacteristics(piece)
						t.Logf("Piece %d (%d): color=%d, shape=%d, size=%d, fill=%d",
							i, piece, color, shape, size, fill)
					}
				}
			}
		})
	}
}

func TestIsValidPosition(t *testing.T) {
	tests := []struct {
		name        string
		row         int
		col         int
		expectedRow bool
		expectedCol bool
	}{
		{"Valid position (0,0)", 0, 0, true, true},
		{"Valid position (3,3)", 3, 3, true, true},
		{"Valid position (1,2)", 1, 2, true, true},
		{"Invalid row (negative)", -1, 0, false, true},
		{"Invalid row (too high)", 4, 0, false, true},
		{"Invalid col (negative)", 0, -1, true, false},
		{"Invalid col (too high)", 0, 4, true, false},
		{"Both invalid", -1, 4, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultRow := IsValidRow(tt.row)
			resultCol := IsValidCol(tt.col)

			if resultRow != tt.expectedRow {
				t.Errorf("IsValidRow(%d) = %v, expected %v", tt.row, resultRow, tt.expectedRow)
			}
			if resultCol != tt.expectedCol {
				t.Errorf("IsValidCol(%d) = %v, expected %v", tt.col, resultCol, tt.expectedCol)
			}
		})
	}
}

func TestIsValidPiece(t *testing.T) {
	tests := []struct {
		name     string
		piece    Piece
		expected bool
	}{
		{"Valid piece (min)", Piece(0), true},
		{"Valid piece (max)", Piece(15), true},
		{"Valid piece (middle)", Piece(8), true},
		{"Invalid piece (negative)", Piece(-1), false},
		{"Invalid piece (too high)", Piece(17), false},
		{"Invalid piece (way too high)", Piece(100), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPiece(tt.piece)
			if result != tt.expected {
				t.Errorf("IsValidPiece(%d) = %v, expected %v", tt.piece, result, tt.expected)
			}
		})
	}
}

func TestGetEmptyBoard(t *testing.T) {
	board := GetEmptyBoard()

	// Vérifier que toutes les cases sont vides
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] != PieceEmpty {
				t.Errorf("Position [%d][%d] should be empty, got %d", i, j, board[i][j])
			}
		}
	}
}

func TestGetAllPieces(t *testing.T) {
	pieces := GetAllPieces()

	// Vérifier qu'on a bien 16 pièces
	if len(pieces) != 16 {
		t.Errorf("Expected 16 pieces, got %d", len(pieces))
	}

	// Vérifier que toutes les pièces sont uniques et valides
	seen := make(map[Piece]bool)
	for _, piece := range pieces {
		if seen[piece] {
			t.Errorf("Duplicate piece found: %d", piece)
		}
		seen[piece] = true

		if !IsValidPiece(piece) {
			t.Errorf("Invalid piece found: %d", piece)
		}
	}

	// Vérifier que toutes les pièces de 0 à 15 sont présentes
	for i := 0; i < 16; i++ {
		expected := Piece(i)
		if !seen[expected] {
			t.Errorf("Missing piece: %d", expected)
		}
	}
}
