package game

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateNewGame crée une nouvelle partie entre deux joueurs
func CreateNewGame(player1ID, player2ID int64) (g Game, err error) {
	g = InitializeGame(player1ID, player2ID)
	g.ID = uuid.New().String()
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()

	err = CreateGame(g)
	return
}

// SelectPiece sélectionne une pièce pour le prochain coup
func (g *Game) SelectPiece(piece Piece) error {

	// Vérifier que c'est la phase de sélection
	if g.GamePhase != GamePhaseSelectPiece {
		return fmt.Errorf("ce n'est pas la phase de sélection de pièce")
	}

	pieceAvailable := false
	for _, available := range g.AvailablePieces {
		if available == piece {
			pieceAvailable = true
			break
		}
	}

	if !pieceAvailable {
		return fmt.Errorf("cette pièce n'est pas disponible")
	}

	// Mettre à jour le jeu
	g.SelectedPiece = piece
	g.GamePhase = GamePhasePlacePiece

	// Retirer la pièce sélectionnée des pièces disponibles
	var newAvailablePieces []Piece
	for _, availablePiece := range g.AvailablePieces {
		if availablePiece != piece {
			newAvailablePieces = append(newAvailablePieces, availablePiece)
		}
	}
	g.AvailablePieces = newAvailablePieces

	g.switchTurn()
	g.UpdatedAt = time.Now()

	err := UpdateGame(*g)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du jeu: %v", err)
	}

	return nil
}

// PlacePiece place une pièce sur le plateau
func (g *Game) PlacePiece(position Position) (err error) {

	// Vérifier que c'est la phase de placement
	if g.GamePhase != GamePhasePlacePiece {
		return fmt.Errorf("ce n'est pas la phase de placement de pièce")
	}

	// Vérifier qu'une pièce est sélectionnée
	if !IsValidPiece(g.SelectedPiece) || g.SelectedPiece == PieceEmpty {
		return fmt.Errorf("aucune pièce n'est sélectionnée")
	}

	if g.Board[position.Row][position.Col] != PieceEmpty {
		return fmt.Errorf("cette position est déjà occupée")
	}

	// Placer la pièce
	g.Board[position.Row][position.Col] = g.SelectedPiece

	// Mettre à jour l'historique des mouvements
	g.History = append(g.History, Move{
		Piece:    g.SelectedPiece,
		Position: position,
	})

	g.SelectedPiece = PieceEmpty // Réinitialiser la pièce sélectionnée
	g.UpdatedAt = time.Now()

	// Vérifier les conditions de victoire
	if CheckWin(g.Board) {
		g.Status = StatusFinished
		g.Winner = g.CurrentTurn
	} else if len(g.AvailablePieces) == 0 {
		// Match nul - toutes les pièces ont été placées
		g.Status = StatusFinished
		g.Winner = 0
	} else {
		// Continuer le jeu - phase de sélection pour le prochain joueur
		g.GamePhase = GamePhaseSelectPiece
	}

	err = UpdateGame(*g)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du jeu: %v", err)
	}

	return
}

// GetGame récupère une partie et vérifie les droits d'accès
func GetGame(gameID string, userID int64) (g Game, err error) {
	g, err = GetGameByID(gameID)
	if err != nil {
		return
	}

	// Vérifier que l'utilisateur fait partie de cette partie
	if g.Player1ID != userID && g.Player2ID != userID {
		err = fmt.Errorf("vous n'avez pas accès à cette partie")
		return
	}

	return g, nil
}

// switchTurn
func (g *Game) switchTurn() {
	switch g.CurrentTurn {
	case g.Player1ID:
		g.CurrentTurn = g.Player2ID
	case g.Player2ID:
		g.CurrentTurn = g.Player1ID
	}
}

// ForfeitGame abandonne une partie
func (g *Game) ForfeitGame(userID int64) (err error) {

	// Vérifier que la partie est active
	if g.Status != StatusPlaying {
		return fmt.Errorf("cette partie n'est plus active")
	}

	// Déterminer le gagnant (l'autre joueur)
	var winner int64
	if g.Player1ID == userID {
		winner = g.Player2ID
	} else {
		winner = g.Player1ID
	}

	// Mettre à jour la partie
	g.Status = StatusFinished
	g.Winner = winner
	g.UpdatedAt = time.Now()

	err = UpdateGame(*g)
	return
}
