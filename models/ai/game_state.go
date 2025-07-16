package ai

import (
	"fmt"
	"quarto/models/game"
)

// ConvertGameToState convertit un Game en GameState pour l'IA
func ConvertGameToState(g game.Game) GameState {

	winner := 0
	switch g.Winner {
	case g.Player1ID:
		winner = 1
	case g.Player2ID:
		winner = -1
	}

	return GameState{
		Board:           g.Board,
		AvailablePieces: g.AvailablePieces,
		SelectedPiece:   g.SelectedPiece,
		IsGameOver:      g.Status == game.StatusFinished,
		Winner:          winner,
	}
}

// ConvertHistoryToGameState reconstruit un GameState à partir d'un historique de mouvements
func ConvertHistoryToGameState(moves []game.Move) GameState {
	// Initialiser l'état du jeu au début
	state := GameState{
		Board:           game.GetEmptyBoard(),
		AvailablePieces: game.GetAllPieces(),
		SelectedPiece:   game.PieceEmpty,
		IsGameOver:      false,
	}

	// Appliquer chaque mouvement de l'historique
	for _, move := range moves {
		state.SelectedPiece = move.Piece
		state = state.ApplyMove(AIMove{
			Move:          move,
			SelectedPiece: game.PieceEmpty,
		})

		fmt.Printf("Apply move %v %d:\n", move.Position, move.Piece)
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if state.Board[i][j] == game.PieceEmpty {
					fmt.Printf("[ ] ")
				} else {
					fmt.Printf("[%d] ", state.Board[i][j])
				}
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		fmt.Printf("Available pieces: %v\n", state.AvailablePieces)
	}

	return state
}

// CheckGameOver vérifie si le jeu est terminé et retourne si il y a un gagnant
func (state GameState) CheckGameOver() (bool, bool) {
	win := game.CheckWin(state.Board)
	return win || len(state.AvailablePieces) == 0, win
}

// GetValidMoves retourne tous les mouvements valides pour l'état actuel (version classique)
func GetValidMoves(state GameState) []AIMove {
	basicMoves := game.GetValidMoves(game.GamePhasePlacePiece, state.Board, []game.Piece{state.SelectedPiece})

	var moves []AIMove
	// Add piece selection moves
	for _, move := range basicMoves {
		for _, piece := range state.AvailablePieces {
			if piece != move.Piece {
				moves = append(moves, AIMove{
					Move:          move,
					SelectedPiece: piece,
				})
			}
		}
	}
	return moves
}

// ApplyMove applique un mouvement à l'état du jeu et retourne le nouvel état (version classique)
func (state GameState) ApplyMove(move AIMove) GameState {

	state.Board[move.Move.Position.Row][move.Move.Position.Col] = move.Move.Piece
	gameOver, winner := state.CheckGameOver()
	if gameOver {
		state.IsGameOver = true
		if winner {
			state.Winner = -len(state.AvailablePieces)%2*2 + 1
		}
	}

	newAvailablePieces := make([]game.Piece, 0, len(state.AvailablePieces))
	for _, pieceID := range state.AvailablePieces {
		if pieceID == move.SelectedPiece {
			continue // Ne pas garder la pièce sélectionnée pour l'adversaire
		}
		if pieceID == move.Move.Piece {
			continue // Ne pas garder la pièce qui vient d'être placée
		}
		newAvailablePieces = append(newAvailablePieces, pieceID)
	}
	state.AvailablePieces = newAvailablePieces

	// Sélectionner la pièce pour l'adversaire
	state.SelectedPiece = move.SelectedPiece

	return state
}
