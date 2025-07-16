package game

import (
	"time"
)

// CheckWin vérifie s'il y a une victoire sur le plateau
func CheckWin(board [4][4]Piece) bool {
	// Vérifier toutes les lignes, colonnes et diagonales
	return checkLines(board) || checkColumns(board) || checkDiagonals(board)
}

// checkLines vérifie les lignes horizontales
func checkLines(board [4][4]Piece) bool {
	for i := range 4 {
		if hasCommonCharacteristic([]Piece{board[i][0], board[i][1], board[i][2], board[i][3]}) {
			return true
		}
	}
	return false
}

// checkColumns vérifie les colonnes verticales
func checkColumns(board [4][4]Piece) bool {
	for j := range 4 {
		if hasCommonCharacteristic([]Piece{board[0][j], board[1][j], board[2][j], board[3][j]}) {
			return true
		}
	}
	return false
}

// checkDiagonals vérifie les diagonales
func checkDiagonals(board [4][4]Piece) bool {
	// Diagonale principale
	if hasCommonCharacteristic([]Piece{board[0][0], board[1][1], board[2][2], board[3][3]}) {
		return true
	}
	// Diagonale secondaire
	if hasCommonCharacteristic([]Piece{board[0][3], board[1][2], board[2][1], board[3][0]}) {
		return true
	}
	return false
}

// hasCommonCharacteristic vérifie si 4 pièces ont au moins une caractéristique commune
func hasCommonCharacteristic(pieces []Piece) bool {
	if len(pieces) != 4 {
		return false
	}
	color, shape, size, fill := GetPieceCharacteristics(pieces[0])
	matchColor, matchShape, matchSize, matchFill := true, true, true, true

	for _, piece := range pieces[1:] {
		if piece == PieceEmpty {
			return false
		}
		color2, shape2, size2, fill2 := GetPieceCharacteristics(piece)
		if color2 != color {
			matchColor = false
		}
		if shape2 != shape {
			matchShape = false
		}
		if size2 != size {
			matchSize = false
		}
		if fill2 != fill {
			matchFill = false
		}
	}
	return matchColor || matchShape || matchSize || matchFill
}

func IsValidRow(row int) bool {
	return row >= 0 && row < 4
}

func IsValidCol(col int) bool {
	return col >= 0 && col < 4
}

func IsValidPiece(piece Piece) bool {
	return piece >= 0 && piece <= 15
}

func GetEmptyBoard() [4][4]Piece {
	return [4][4]Piece{
		{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
		{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
		{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
		{PieceEmpty, PieceEmpty, PieceEmpty, PieceEmpty},
	}
}

// InitializeGame initialise un nouveau jeu
func InitializeGame(player1ID, player2ID int64) Game {
	availablePieces := GetAllPieces()

	return Game{
		Player1ID:       player1ID,
		Player2ID:       player2ID,
		CurrentTurn:     player1ID, // Le joueur 1 commence
		GamePhase:       GamePhaseSelectPiece,
		Board:           GetEmptyBoard(),
		AvailablePieces: availablePieces,
		Status:          StatusPlaying,
		Winner:          0,
		History:         []Move{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// ApplyMove applique un mouvement sur un plateau et retourne le nouveau plateau
func ApplyMove(board [4][4]Piece, move Move) [4][4]Piece {
	newBoard := board // Copie du plateau
	if IsValidRow(move.Position.Row) && IsValidCol(move.Position.Col) &&
		newBoard[move.Position.Row][move.Position.Col] == PieceEmpty {
		newBoard[move.Position.Row][move.Position.Col] = move.Piece
	}
	return newBoard
}

// CanApplyMove vérifie si un mouvement peut être appliqué sur un plateau
func CanApplyMove(board [4][4]Piece, move Move) bool {
	if !IsValidRow(move.Position.Row) || !IsValidCol(move.Position.Col) {
		return false
	}
	return board[move.Position.Row][move.Position.Col] == PieceEmpty
}

func GetValidMoves(gamePhase int, board [4][4]Piece, availablePieces []Piece) (moves []Move) {
	switch gamePhase {
	case GamePhaseSelectPiece:
		for _, piece := range availablePieces {
			moves = append(moves, Move{Piece: piece})
		}
	case GamePhasePlacePiece:
		for row := range 4 {
			for col := range 4 {
				if board[row][col] == PieceEmpty {
					for _, piece := range availablePieces {
						moves = append(moves, Move{
							Piece:    piece,
							Position: Position{Row: row, Col: col},
						})
					}
				}
			}
		}
	}
	return
}
