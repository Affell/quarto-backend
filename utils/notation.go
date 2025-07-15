package utils

import (
	"encoding/json"
	"fmt"
	"quarto/models/game"
	"strconv"
	"strings"
)

// PieceToNotation convertit une pièce en notation algébrique
func PieceToNotation(piece game.Piece) string {
	color := "B"
	if piece.Color == "noir" {
		color = "N"
	}

	shape := "C"
	if piece.Shape == "rond" {
		shape = "R"
	}

	size := "G"
	if piece.Size == "petit" {
		size = "P"
	}

	fill := "P"
	if piece.Fill == "troué" {
		fill = "T"
	}

	return fmt.Sprintf("%s%s%s%s", color, shape, size, fill)
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
func CreateMoveNotation(piece game.Piece, position string) string {
	notation := PieceToNotation(piece)
	return fmt.Sprintf("%s-%s", notation, position)
}

// ParseMoveNotation parse une notation de coup
func ParseMoveNotation(notation string) (game.Piece, string, error) {
	parts := strings.Split(notation, "-")
	if len(parts) != 2 {
		return game.Piece{}, "", fmt.Errorf("notation invalide: %s", notation)
	}

	pieceNotation := parts[0]
	position := parts[1]

	if len(pieceNotation) != 4 {
		return game.Piece{}, "", fmt.Errorf("notation de pièce invalide: %s", pieceNotation)
	}

	piece := game.Piece{}

	// Couleur
	switch pieceNotation[0] {
	case 'B':
		piece.Color = "blanc"
	case 'N':
		piece.Color = "noir"
	default:
		return game.Piece{}, "", fmt.Errorf("couleur invalide: %c", pieceNotation[0])
	}

	// Forme
	switch pieceNotation[1] {
	case 'C':
		piece.Shape = "carré"
	case 'R':
		piece.Shape = "rond"
	default:
		return game.Piece{}, "", fmt.Errorf("forme invalide: %c", pieceNotation[1])
	}

	// Taille
	switch pieceNotation[2] {
	case 'G':
		piece.Size = "grand"
	case 'P':
		piece.Size = "petit"
	default:
		return game.Piece{}, "", fmt.Errorf("taille invalide: %c", pieceNotation[2])
	}

	// Remplissage
	switch pieceNotation[3] {
	case 'P':
		piece.Fill = "plein"
	case 'T':
		piece.Fill = "troué"
	default:
		return game.Piece{}, "", fmt.Errorf("remplissage invalide: %c", pieceNotation[3])
	}

	return piece, position, nil
}

// FindPieceByNotation trouve une pièce par sa notation dans la liste des pièces disponibles
func FindPieceByNotation(notation string, pieces []game.Piece) (*game.Piece, error) {
	targetPiece, _, err := ParseMoveNotation(notation + "-a1") // Position temporaire
	if err != nil {
		return nil, err
	}

	for _, piece := range pieces {
		if piece.Color == targetPiece.Color &&
			piece.Shape == targetPiece.Shape &&
			piece.Size == targetPiece.Size &&
			piece.Fill == targetPiece.Fill {
			return &piece, nil
		}
	}

	return nil, fmt.Errorf("pièce non trouvée: %s", notation)
}

// GetPieceByID retourne une pièce par son ID
func GetPieceByID(id int, pieces []game.Piece) (*game.Piece, error) {
	for _, piece := range pieces {
		if piece.ID == id {
			return &piece, nil
		}
	}
	return nil, fmt.Errorf("pièce avec ID %d non trouvée", id)
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
