package game

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

// CreateNewGame crée une nouvelle partie entre deux joueurs
func CreateNewGame(player1ID, player2ID int64) (*Game, error) {
	newGame := InitializeGame(player1ID, player2ID)
	newGame.ID = uuid.New().String()
	newGame.CreatedAt = time.Now()
	newGame.UpdatedAt = time.Now()

	err := CreateGame(*newGame)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la partie: %v", err)
	}

	return newGame, nil
}

// SelectPiece sélectionne une pièce pour le prochain coup
func SelectPiece(gameID string, userID int64, pieceID int) (*Game, error) {
	g, err := GetGameByID(gameID)
	if err != nil {
		return nil, err
	}

	// Vérifier que c'est le bon joueur
	if !isPlayerTurn(g, userID) {
		return nil, fmt.Errorf("ce n'est pas votre tour")
	}

	// Vérifier que c'est la phase de sélection
	if g.GamePhase != "selectPiece" {
		return nil, fmt.Errorf("ce n'est pas la phase de sélection de pièce")
	}

	// Vérifier que la pièce est disponible
	var availablePieces []int
	if err := json.Unmarshal([]byte(g.AvailablePieces), &availablePieces); err != nil {
		return nil, fmt.Errorf("erreur de désérialisation des pièces disponibles: %v", err)
	}

	pieceAvailable := false
	for _, available := range availablePieces {
		if available == pieceID {
			pieceAvailable = true
			break
		}
	}

	if !pieceAvailable {
		return nil, fmt.Errorf("cette pièce n'est pas disponible")
	}

	// Mettre à jour le jeu
	g.SelectedPiece = &pieceID
	g.GamePhase = "placePiece"
	g.CurrentTurn = getOtherPlayer(g.CurrentTurn)
	g.UpdatedAt = time.Now()

	err = UpdateGame(*g)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la mise à jour du jeu: %v", err)
	}

	log.Debug("Pièce sélectionnée", "game", g.ID, "piece", pieceID, "nextPlayer", g.CurrentTurn)
	return g, nil
}

// PlacePiece place une pièce sur le plateau
func PlacePiece(gameID string, userID int64, position string) (*Game, error) {
	g, err := GetGameByID(gameID)
	if err != nil {
		return nil, err
	}

	// Vérifier que c'est le bon joueur
	if !isPlayerTurn(g, userID) {
		return nil, fmt.Errorf("ce n'est pas votre tour")
	}

	// Vérifier que c'est la phase de placement
	if g.GamePhase != "placePiece" {
		return nil, fmt.Errorf("ce n'est pas la phase de placement de pièce")
	}

	// Vérifier qu'une pièce est sélectionnée
	if g.SelectedPiece == nil {
		return nil, fmt.Errorf("aucune pièce n'est sélectionnée")
	}

	// Convertir la position en coordonnées
	col, row := ParsePosition(position)
	if col == -1 || row == -1 {
		return nil, fmt.Errorf("position invalide: %s", position)
	}

	// Vérifier que la position est libre
	var board [][]interface{}
	if err := json.Unmarshal([]byte(g.Board), &board); err != nil {
		return nil, fmt.Errorf("erreur de désérialisation du plateau: %v", err)
	}

	if board[row][col] != nil {
		return nil, fmt.Errorf("cette position est déjà occupée")
	}

	// Placer la pièce
	board[row][col] = *g.SelectedPiece

	// Retirer la pièce des pièces disponibles
	var availablePieces []int
	if err := json.Unmarshal([]byte(g.AvailablePieces), &availablePieces); err != nil {
		return nil, fmt.Errorf("erreur de désérialisation des pièces disponibles: %v", err)
	}

	for i, piece := range availablePieces {
		if piece == *g.SelectedPiece {
			availablePieces = append(availablePieces[:i], availablePieces[i+1:]...)
			break
		}
	}

	// Sérialiser les données mises à jour
	boardBytes, _ := json.Marshal(board)
	availablePiecesBytes, _ := json.Marshal(availablePieces)

	// Sauvegarder la pièce placée pour les logs
	placedPiece := *g.SelectedPiece

	// Ajouter le mouvement à l'historique
	var moveHistory []string
	if g.MoveHistory != "" {
		json.Unmarshal([]byte(g.MoveHistory), &moveHistory)
	}

	notation := GenerateNotation(placedPiece, position)
	moveHistory = append(moveHistory, notation)
	moveHistoryBytes, _ := json.Marshal(moveHistory)

	// Mettre à jour le jeu
	g.Board = string(boardBytes)
	g.AvailablePieces = string(availablePiecesBytes)
	g.MoveHistory = string(moveHistoryBytes)
	g.SelectedPiece = nil
	g.UpdatedAt = time.Now()

	// Vérifier les conditions de victoire
	if CheckWin(board) {
		g.Status = "finished"
		winner := g.CurrentTurn
		g.Winner = &winner
	} else if len(availablePieces) == 0 {
		// Match nul - toutes les pièces ont été placées
		g.Status = "finished"
		draw := "draw"
		g.Winner = &draw
	} else {
		// Continuer le jeu - phase de sélection pour le prochain joueur
		g.GamePhase = "selectPiece"
		// Le tour reste au même joueur qui doit sélectionner la pièce
	}

	err = UpdateGame(*g)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la mise à jour du jeu: %v", err)
	}

	log.Debug("Pièce placée", "game", g.ID, "position", position, "piece", placedPiece)
	return g, nil
}

// GetGame récupère une partie et vérifie les droits d'accès
func GetGame(gameID string, userID int64) (*Game, error) {
	g, err := GetGameByID(gameID)
	if err != nil {
		return nil, err
	}

	// Vérifier que l'utilisateur fait partie de cette partie
	if g.Player1ID != userID && g.Player2ID != userID {
		return nil, fmt.Errorf("vous n'avez pas accès à cette partie")
	}

	return g, nil
}

// isPlayerTurn vérifie si c'est le tour du joueur
func isPlayerTurn(g *Game, userID int64) bool {
	if g.CurrentTurn == "player1" && g.Player1ID == userID {
		return true
	}
	if g.CurrentTurn == "player2" && g.Player2ID == userID {
		return true
	}
	return false
}

// getOtherPlayer retourne l'autre joueur
func getOtherPlayer(currentPlayer string) string {
	if currentPlayer == "player1" {
		return "player2"
	}
	return "player1"
}

// ForfeitGame abandonne une partie
func ForfeitGame(gameID string, userID int64) (*Game, error) {
	g, err := GetGameByID(gameID)
	if err != nil {
		return nil, err
	}

	// Vérifier que l'utilisateur fait partie de cette partie
	if g.Player1ID != userID && g.Player2ID != userID {
		return nil, fmt.Errorf("vous n'avez pas accès à cette partie")
	}

	// Vérifier que la partie est active
	if g.Status != "active" {
		return nil, fmt.Errorf("cette partie n'est plus active")
	}

	// Déterminer le gagnant (l'autre joueur)
	var winner string
	if g.Player1ID == userID {
		winner = "player2"
	} else {
		winner = "player1"
	}

	// Mettre à jour la partie
	g.Status = "finished"
	g.Winner = &winner
	g.UpdatedAt = time.Now()

	err = UpdateGame(*g)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'abandon de la partie: %v", err)
	}

	return g, nil
}
