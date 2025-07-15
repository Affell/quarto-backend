package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"quarto/models/game"
)

// CheckWin vérifie s'il y a une victoire sur le plateau
func CheckWin(board [4][4]*int, pieces []game.Piece) bool {
	// Vérifier toutes les lignes, colonnes et diagonales
	return checkLines(board, pieces) || checkColumns(board, pieces) || checkDiagonals(board, pieces)
}

// checkLines vérifie les lignes horizontales
func checkLines(board [4][4]*int, pieces []game.Piece) bool {
	for i := 0; i < 4; i++ {
		if checkSequence([]*int{board[i][0], board[i][1], board[i][2], board[i][3]}, pieces) {
			return true
		}
	}
	return false
}

// checkColumns vérifie les colonnes verticales
func checkColumns(board [4][4]*int, pieces []game.Piece) bool {
	for j := 0; j < 4; j++ {
		if checkSequence([]*int{board[0][j], board[1][j], board[2][j], board[3][j]}, pieces) {
			return true
		}
	}
	return false
}

// checkDiagonals vérifie les diagonales
func checkDiagonals(board [4][4]*int, pieces []game.Piece) bool {
	// Diagonale principale
	if checkSequence([]*int{board[0][0], board[1][1], board[2][2], board[3][3]}, pieces) {
		return true
	}
	// Diagonale secondaire
	if checkSequence([]*int{board[0][3], board[1][2], board[2][1], board[3][0]}, pieces) {
		return true
	}
	return false
}

// checkSequence vérifie si 4 pièces ont au moins une caractéristique commune
func checkSequence(sequence []*int, pieces []game.Piece) bool {
	// Vérifier que toutes les positions sont occupées
	for _, pieceID := range sequence {
		if pieceID == nil {
			return false
		}
	}

	// Récupérer les pièces correspondantes
	sequencePieces := make([]game.Piece, 4)
	for i, pieceID := range sequence {
		piece, err := GetPieceByID(*pieceID, pieces)
		if err != nil {
			return false
		}
		sequencePieces[i] = *piece
	}

	// Vérifier si au moins une caractéristique est commune
	return hasCommonCharacteristic(sequencePieces)
}

// hasCommonCharacteristic vérifie si 4 pièces ont au moins une caractéristique commune
func hasCommonCharacteristic(pieces []game.Piece) bool {
	if len(pieces) != 4 {
		return false
	}

	// Vérifier la couleur
	if pieces[0].Color == pieces[1].Color && pieces[1].Color == pieces[2].Color && pieces[2].Color == pieces[3].Color {
		return true
	}

	// Vérifier la forme
	if pieces[0].Shape == pieces[1].Shape && pieces[1].Shape == pieces[2].Shape && pieces[2].Shape == pieces[3].Shape {
		return true
	}

	// Vérifier la taille
	if pieces[0].Size == pieces[1].Size && pieces[1].Size == pieces[2].Size && pieces[2].Size == pieces[3].Size {
		return true
	}

	// Vérifier le remplissage
	if pieces[0].Fill == pieces[1].Fill && pieces[1].Fill == pieces[2].Fill && pieces[2].Fill == pieces[3].Fill {
		return true
	}

	return false
}

// IsValidMove vérifie si un coup est valide
func IsValidMove(g *game.Game, pieceID int, position string) error {
	// Vérifier que c'est la bonne phase
	if g.GamePhase == "selectPiece" && position != "" {
		return errors.New("impossible de placer une pièce pendant la phase de sélection")
	}
	if g.GamePhase == "placePiece" && position == "" {
		return errors.New("position requise pendant la phase de placement")
	}

	// Vérifier les pièces disponibles
	var availablePieces []int
	if err := json.Unmarshal([]byte(g.AvailablePieces), &availablePieces); err != nil {
		return fmt.Errorf("erreur de parsing des pièces disponibles: %v", err)
	}

	// Pour la sélection de pièce
	if g.GamePhase == "selectPiece" {
		found := false
		for _, id := range availablePieces {
			if id == pieceID {
				found = true
				break
			}
		}
		if !found {
			return errors.New("pièce non disponible")
		}
	}

	// Pour le placement de pièce
	if g.GamePhase == "placePiece" {
		if g.SelectedPiece == nil || *g.SelectedPiece != pieceID {
			return errors.New("cette pièce n'est pas sélectionnée")
		}

		// Vérifier que la position est valide
		row, col, err := PositionToCoords(position)
		if err != nil {
			return fmt.Errorf("position invalide: %v", err)
		}

		// Vérifier que la position est libre
		board, err := BoardToMatrix(g.Board)
		if err != nil {
			return fmt.Errorf("erreur de parsing du plateau: %v", err)
		}

		if board[row][col] != nil {
			return errors.New("position déjà occupée")
		}
	}

	return nil
}

// GetAvailablePositions retourne toutes les positions libres sur le plateau
func GetAvailablePositions(boardJSON string) ([]string, error) {
	board, err := BoardToMatrix(boardJSON)
	if err != nil {
		return nil, err
	}

	var positions []string
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == nil {
				positions = append(positions, CoordsToPosition(i, j))
			}
		}
	}

	return positions, nil
}

// InitializeGame initialise un nouveau jeu
func InitializeGame(player1ID, player2ID int64) *game.Game {
	// Toutes les pièces sont disponibles au début
	availablePieces := make([]int, 16)
	for i := 0; i < 16; i++ {
		availablePieces[i] = i
	}

	availablePiecesJSON, _ := json.Marshal(availablePieces)

	// Plateau vide (4x4 avec des nulls)
	board := "[[null,null,null,null],[null,null,null,null],[null,null,null,null],[null,null,null,null]]"

	// Historique vide
	moveHistory := "[]"

	return &game.Game{
		Player1ID:       player1ID,
		Player2ID:       player2ID,
		CurrentTurn:     "player1",
		GamePhase:       "selectPiece",
		Board:           board,
		AvailablePieces: string(availablePiecesJSON),
		SelectedPiece:   nil,
		Status:          "playing",
		Winner:          nil,
		MoveHistory:     moveHistory,
	}
}

// IsGameFull vérifie si le plateau est plein
func IsGameFull(boardJSON string) bool {
	positions, err := GetAvailablePositions(boardJSON)
	if err != nil {
		return false
	}
	return len(positions) == 0
}

// RemovePieceFromAvailable retire une pièce de la liste des pièces disponibles
func RemovePieceFromAvailable(availablePiecesJSON string, pieceID int) (string, error) {
	var pieces []int
	if err := json.Unmarshal([]byte(availablePiecesJSON), &pieces); err != nil {
		return "", err
	}

	for i, id := range pieces {
		if id == pieceID {
			pieces = append(pieces[:i], pieces[i+1:]...)
			break
		}
	}

	result, err := json.Marshal(pieces)
	return string(result), err
}

// AddMoveToHistory ajoute un coup à l'historique
func AddMoveToHistory(historyJSON string, moveNotation string) (string, error) {
	var history []string
	if historyJSON != "" && historyJSON != "[]" {
		if err := json.Unmarshal([]byte(historyJSON), &history); err != nil {
			return "", err
		}
	}

	history = append(history, moveNotation)
	result, err := json.Marshal(history)
	return string(result), err
}
