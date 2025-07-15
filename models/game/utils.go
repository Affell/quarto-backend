package game

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParsePosition convertit une position type "a1" en coordonnées (col, row)
func ParsePosition(position string) (col, row int) {
	if len(position) != 2 {
		return -1, -1
	}

	// Convertir la colonne (a-d)
	switch position[0] {
	case 'a', 'A':
		col = 0
	case 'b', 'B':
		col = 1
	case 'c', 'C':
		col = 2
	case 'd', 'D':
		col = 3
	default:
		return -1, -1
	}

	// Convertir la ligne (1-4)
	switch position[1] {
	case '1':
		row = 0
	case '2':
		row = 1
	case '3':
		row = 2
	case '4':
		row = 3
	default:
		return -1, -1
	}

	return col, row
}

// GenerateNotation génère la notation algébrique pour un mouvement
func GenerateNotation(pieceID int, position string) string {
	// Convertir l'ID de pièce en notation BCGP
	pieces := GetAllPieces()
	if pieceID >= 0 && pieceID < len(pieces) {
		piece := pieces[pieceID]
		notation := ""

		// Couleur
		if piece.Color == "blanc" {
			notation += "B"
		} else {
			notation += "N"
		}

		// Forme
		if piece.Shape == "carré" {
			notation += "C"
		} else {
			notation += "R"
		}

		// Taille
		if piece.Size == "grand" {
			notation += "G"
		} else {
			notation += "P"
		}

		// Remplissage
		if piece.Fill == "plein" {
			notation += "P"
		} else {
			notation += "T"
		}

		return fmt.Sprintf("%s-%s", notation, strings.ToLower(position))
	}

	return fmt.Sprintf("P%d-%s", pieceID, strings.ToLower(position))
}

// CheckWin vérifie s'il y a une victoire sur le plateau
func CheckWin(board [][]interface{}) bool {
	// Convertir le board en format utilisable
	var gameBoard [4][4]*int
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] != nil {
				// Convertir l'interface{} en int
				switch v := board[i][j].(type) {
				case int:
					gameBoard[i][j] = &v
				case float64:
					intVal := int(v)
					gameBoard[i][j] = &intVal
				}
			}
		}
	}

	pieces := GetAllPieces()
	return checkWinInternal(gameBoard, pieces)
}

// checkWinInternal vérifie les conditions de victoire
func checkWinInternal(board [4][4]*int, pieces []Piece) bool {
	// Vérifier les lignes
	for row := 0; row < 4; row++ {
		if checkLine(board[row][:], pieces) {
			return true
		}
	}

	// Vérifier les colonnes
	for col := 0; col < 4; col++ {
		var line [4]*int
		for row := 0; row < 4; row++ {
			line[row] = board[row][col]
		}
		if checkLine(line[:], pieces) {
			return true
		}
	}

	// Vérifier les diagonales
	var diag1, diag2 [4]*int
	for i := 0; i < 4; i++ {
		diag1[i] = board[i][i]
		diag2[i] = board[i][3-i]
	}

	return checkLine(diag1[:], pieces) || checkLine(diag2[:], pieces)
}

// checkLine vérifie si une ligne a 4 pièces avec une caractéristique commune
func checkLine(line []*int, pieces []Piece) bool {
	// Vérifier que toutes les positions sont occupées
	for _, pos := range line {
		if pos == nil {
			return false
		}
	}

	// Récupérer les pièces de la ligne
	var linePieces []Piece
	for _, pieceID := range line {
		if *pieceID >= 0 && *pieceID < len(pieces) {
			linePieces = append(linePieces, pieces[*pieceID])
		}
	}

	if len(linePieces) != 4 {
		return false
	}

	// Vérifier chaque caractéristique
	characteristics := []string{"color", "shape", "size", "fill"}
	for _, char := range characteristics {
		if hasCommonCharacteristic(linePieces, char) {
			return true
		}
	}

	return false
}

// hasCommonCharacteristic vérifie si toutes les pièces ont la même caractéristique
func hasCommonCharacteristic(pieces []Piece, characteristic string) bool {
	if len(pieces) == 0 {
		return false
	}

	var firstValue string
	switch characteristic {
	case "color":
		firstValue = pieces[0].Color
	case "shape":
		firstValue = pieces[0].Shape
	case "size":
		firstValue = pieces[0].Size
	case "fill":
		firstValue = pieces[0].Fill
	default:
		return false
	}

	for _, piece := range pieces {
		var currentValue string
		switch characteristic {
		case "color":
			currentValue = piece.Color
		case "shape":
			currentValue = piece.Shape
		case "size":
			currentValue = piece.Size
		case "fill":
			currentValue = piece.Fill
		}

		if currentValue != firstValue {
			return false
		}
	}

	return true
}

// InitializeGame crée une nouvelle instance de jeu
func InitializeGame(player1ID, player2ID int64) *Game {
	// Plateau vide (4x4)
	board := make([][]interface{}, 4)
	for i := range board {
		board[i] = make([]interface{}, 4)
	}

	// Toutes les pièces disponibles (0-15)
	availablePieces := make([]int, 16)
	for i := 0; i < 16; i++ {
		availablePieces[i] = i
	}

	// Sérialiser en JSON
	boardBytes, _ := json.Marshal(board)
	availablePiecesBytes, _ := json.Marshal(availablePieces)
	moveHistoryBytes, _ := json.Marshal([]string{})

	return &Game{
		Player1ID:       player1ID,
		Player2ID:       player2ID,
		CurrentTurn:     "player1",     // Player1 commence
		GamePhase:       "selectPiece", // Commence par sélectionner une pièce
		Board:           string(boardBytes),
		AvailablePieces: string(availablePiecesBytes),
		SelectedPiece:   nil,
		Status:          "active",
		Winner:          nil,
		MoveHistory:     string(moveHistoryBytes),
	}
}
