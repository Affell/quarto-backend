package game

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func GetPieceCharacteristics(piece Piece) (int, int, int, int) {
	color := int((piece) / 8)
	shape := int(piece % 8 / 4)
	size := int(piece % 4 / 2)
	fill := int(piece % 2)
	return color, shape, size, fill
}

// PieceToNotation convertit une pièce en notation algébrique
func PieceToNotation(piece Piece) string {
	color, shape, size, fill := GetPieceCharacteristics(piece)

	var colorChar, shapeChar, sizeChar, fillChar string
	if color == 0 {
		colorChar = "B" // Blanc
	} else {
		colorChar = "N" // Noir
	}
	if shape == 0 {
		shapeChar = "C" // Carré
	} else {
		shapeChar = "R" // Rond
	}
	if size == 0 {
		sizeChar = "G" // Grand
	} else {
		sizeChar = "P" // Petit
	}
	if fill == 0 {
		fillChar = "P" // Plein
	} else {
		fillChar = "T" // Troué
	}

	return fmt.Sprintf("%s%s%s%s", colorChar, shapeChar, sizeChar, fillChar)
}

// CoordsToPosition convertit les coordonnées de la grille en position algébrique
func CoordsToPosition(row, col int) string {
	file := string(rune('a' + col))
	rank := string(rune('1' + row)) // a1 en haut à gauche
	return file + rank
}

// PositionToCoords convertit une position algébrique en coordonnées
func PositionToCoords(position string) (int, int, error) {
	if len(position) != 2 {
		return 0, 0, fmt.Errorf("position invalide: %s", position)
	}

	file := position[0]
	rank := position[1]

	// Accepter les majuscules et minuscules pour les colonnes
	if file >= 'A' && file <= 'D' {
		file = file + 32 // Convertir en minuscule
	}
	if file < 'a' || file > 'd' {
		return 0, 0, fmt.Errorf("colonne invalide: %c", file)
	}

	if rank < '1' || rank > '4' {
		return 0, 0, fmt.Errorf("ligne invalide: %c", rank)
	}

	col := int(file - 'a')
	row := int(rank - '1') // a1 en haut à gauche

	return row, col, nil
}

// CreateMoveNotation crée la notation d'un coup
func CreateMoveNotation(piece Piece, position string) string {
	notation := PieceToNotation(piece)
	return fmt.Sprintf("%s-%s", notation, position)
}

// ParseMoveNotation parse une notation de coup
func ParseMoveNotation(notation string) (Piece, string, error) {
	parts := strings.Split(notation, "-")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("notation invalide: %s", notation)
	}

	pieceNotation := parts[0]
	position := parts[1]

	if len(pieceNotation) != 4 {
		return 0, "", fmt.Errorf("notation de pièce invalide: %s", pieceNotation)
	}

	piece, err := NotationToPiece(pieceNotation)

	return piece, position, err
}

func NotationToPiece(notation string) (Piece, error) {
	parts := strings.Split(notation, "-")
	pieceNotation := parts[0]

	if len(pieceNotation) != 4 {
		return 0, fmt.Errorf("notation de pièce invalide: %s", pieceNotation)
	}

	color := 0
	if pieceNotation[0] == 'N' {
		color = 1 // Noir
	}
	shape := 0
	if pieceNotation[1] == 'R' {
		shape = 1 // Rond
	}
	size := 0
	if pieceNotation[2] == 'P' {
		size = 1 // Petit
	}
	fill := 0
	if pieceNotation[3] == 'T' {
		fill = 1 // Troué
	}
	// Calculer l'ID de la pièce
	pieceID := color*8 + shape*4 + size*2 + fill
	return Piece(pieceID), nil
}

// BoardToMatrix convertit le JSON board en matrice 4x4
func BoardToMatrix(boardJSON string) ([4][4]*int, error) {
	var board [4][4]*int

	if boardJSON == "" || boardJSON == "null" {
		return board, nil
	}

	// Utiliser le parser JSON standard pour plus de robustesse
	var jsonBoard [][]*int
	if err := json.Unmarshal([]byte(boardJSON), &jsonBoard); err != nil {
		return board, fmt.Errorf("erreur de parsing du plateau: %v", err)
	}

	// Vérifier les dimensions
	if len(jsonBoard) != 4 {
		return board, fmt.Errorf("plateau doit avoir 4 lignes, trouvé %d", len(jsonBoard))
	}

	for i := 0; i < 4; i++ {
		if len(jsonBoard[i]) != 4 {
			return board, fmt.Errorf("ligne %d doit avoir 4 colonnes, trouvé %d", i, len(jsonBoard[i]))
		}
		for j := 0; j < 4; j++ {
			board[i][j] = jsonBoard[i][j]
		}
	}

	return board, nil
}

// MatrixToBoard convertit une matrice 4x4 en JSON
func MatrixToBoard(board [4][4]*int) string {
	result := "["
	for i := 0; i < 4; i++ {
		if i > 0 {
			result += ","
		}
		result += "["
		for j := 0; j < 4; j++ {
			if j > 0 {
				result += ","
			}
			if board[i][j] == nil {
				result += "null"
			} else {
				result += strconv.Itoa(*board[i][j])
			}
		}
		result += "]"
	}
	result += "]"
	return result
}

func PrintBoard(board [4][4]Piece) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == PieceEmpty {
				fmt.Print("[ ] ")
			} else {
				fmt.Printf("[%d] ", board[i][j])
			}
		}
		fmt.Println()
	}
}
