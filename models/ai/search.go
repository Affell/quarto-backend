package ai

import (
	"fmt"
	"math"
	"quarto/models/ai/stats"
	"quarto/models/game"
	"time"
)

// SearchResult contient le résultat d'une recherche
type SearchResult struct {
	BestMoves []AIMove
	Score     int
	Depth     int
	Stats     *stats.PerformanceStats
}

// Search effectue une recherche minimax avec élagage alpha-beta
func (e *Engine) Search(state GameState) SearchResult {
	return e.SearchWithStats(state, nil)
}

// SearchWithStats effectue une recherche avec suivi des statistiques
func (e *Engine) SearchWithStats(state GameState, perfStats *stats.PerformanceStats) SearchResult {

	result := SearchResult{
		Score:     0,
		Depth:     0,
		Stats:     perfStats,
		BestMoves: []AIMove{},
	}

	searchStart := time.Now()
	if len(state.AvailablePieces) == 16 {
		result.BestMoves = []AIMove{
			{
				SelectedPiece: state.AvailablePieces[0], // Sélectionne la première pièce disponible
			},
		}
		return result
	}

	if state.SelectedPiece == game.PieceEmpty { // Si aucune pièce n'est sélectionnée, on ne peut pas jouer
		return result
	}

	// Vider la table de transposition pour chaque nouvelle recherche
	e.TT.Clear()

	isMaximizing := len(state.AvailablePieces)%2 == 0 // Maximiser si c'est le tour du joueur 1
	fmt.Printf("Starting search: Maximizing=%v, AvailablePieces=%d, SelectedPiece=%d\n", isMaximizing, len(state.AvailablePieces), state.SelectedPiece)

	// Recherche avec élagage alpha-beta et table de transposition
	score, bestMoves := e.minimax(state, e.MaxDepth, LOSS_SCORE-1, WIN_SCORE+1, isMaximizing, perfStats)

	if perfStats != nil {
		searchDuration := time.Since(searchStart)
		stateHash := e.hashGameState(state)
		perfStats.RecordOperation("total_search", searchDuration, stateHash)
	}

	result.Score = score
	result.BestMoves = bestMoves
	result.Depth = e.MaxDepth

	return result
}

// minimax implémente l'algorithme minimax avec élagage alpha-beta, table de transposition et retourne la continuation
func (e *Engine) minimax(node GameState, depth int, alpha, beta int, isMaximizing bool, perfStats *stats.PerformanceStats) (int, []AIMove) {
	nodeStart := time.Now()
	stateHash := e.hashGameState(node)

	// Enregistrer la visite du nœud
	if perfStats != nil {
		perfStats.RecordOperation("hash", time.Since(nodeStart), stateHash)
	}

	// Vérifier la table de transposition
	if entry, exists := e.TT.Lookup(stateHash, depth, alpha, beta); exists {
		if perfStats != nil {
			perfStats.RecordOperation("tt_hit", 0, stateHash)
		}
		return entry.Score, entry.BestMoves
	}

	// Condition d'arrêt : jeu terminé ou profondeur maximale atteinte
	if node.IsGameOver || depth == 0 {
		evalStart := time.Now()
		score := e.evaluatePosition(node)

		if perfStats != nil {
			evalDuration := time.Since(evalStart)
			if depth == 0 {
				perfStats.RecordOperation("leaf_eval", evalDuration, stateHash)
			} else {
				perfStats.RecordOperation("terminal_eval", evalDuration, stateHash)
			}
		}
		e.TT.Store(stateHash, score, depth, EXACT, []AIMove{})
		return score, []AIMove{}
	}

	movesStart := time.Now()
	validMoves := GetValidMoves(node)
	if perfStats != nil {
		perfStats.RecordOperation("generate_moves", time.Since(movesStart), stateHash)
	}

	// Si aucun mouvement valide, c'est un match nul
	if len(validMoves) == 0 {
		e.TT.Store(stateHash, DRAW_SCORE, depth, EXACT, []AIMove{})
		return DRAW_SCORE, []AIMove{}
	}

	var bestMoves []AIMove
	originalAlpha := alpha

	// Initialiser le meilleur score selon le type de joueur
	var bestScore int
	if isMaximizing {
		bestScore = math.MinInt32
	} else {
		bestScore = math.MaxInt32
	}
	for _, move := range validMoves {
		moveStart := time.Now()
		newState := node.ApplyMove(move)
		if perfStats != nil {
			moveHash := e.hashMove(move, stateHash)
			perfStats.RecordOperation("apply_move", time.Since(moveStart), moveHash+"-"+stateHash)
		}

		score, moves := e.minimax(newState, depth-1, alpha, beta, !isMaximizing, perfStats)

		// Vérifier si ce mouvement améliore le score ou si c'est un meilleur chemin vers la victoire
		isBetter := (isMaximizing && score > bestScore) || (!isMaximizing && score < bestScore)

		// Prioriser les chemins plus courts vers la victoire à score égal
		isSameScoreButShorterPath := false
		if score == bestScore && score != 0 { // Score non nul (victoire ou défaite)
			currentPathLength := len(moves) + 1 // +1 pour le coup actuel
			bestPathLength := len(bestMoves)

			// Si c'est une victoire immédiate (continuation vide), c'est toujours prioritaire
			if len(moves) == 0 && newState.IsGameOver && newState.Winner != 0 {
				isSameScoreButShorterPath = true
			} else if currentPathLength < bestPathLength && bestPathLength > 1 {
				// Chemin plus court vers la même conclusion
				isSameScoreButShorterPath = true
			}
		}

		if isBetter || isSameScoreButShorterPath {
			bestScore = score
			bestMoves = []AIMove{move}
			if len(moves) > 0 {
				bestMoves = append(bestMoves, moves...)
			}
		}

		// Mise à jour des bornes alpha-beta
		if isMaximizing {
			alpha = max(alpha, score)
		} else {
			beta = min(beta, score)
		}

		// Élagage alpha-beta
		if beta <= alpha {
			if perfStats != nil {
				perfStats.RecordOperation("alpha_beta_prune", 0, stateHash)
			}
			break
		}
	}

	// Déterminer le flag pour la table de transposition
	var flag TTFlag
	if bestScore <= originalAlpha {
		flag = UPPER_BOUND
	} else if bestScore >= beta {
		flag = LOWER_BOUND
	} else {
		flag = EXACT
	}

	e.TT.Store(stateHash, bestScore, depth, flag, bestMoves)
	return bestScore, bestMoves
}

// hashGameState génère un hash pour un état de jeu
func (e *Engine) hashGameState(state GameState) string {
	var hash uint64 = 0

	// Hacher le plateau de jeu (4x4 board)
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			piece := state.Board[row][col]
			if piece != game.PieceEmpty {
				// Position sur 4 bits (0-15)
				position := uint64(row*4 + col)
				// Utiliser directement la valeur de la pièce (plus efficace que d'extraire les caractéristiques)
				hash ^= uint64(piece) << (position * 4)
			}
		}
	}

	// Optimiser la concaténation en évitant sprintf pour la pièce sélectionnée
	return fmt.Sprintf("%016x-%d", hash, state.SelectedPiece)
}

// hashMove génère un hash pour un mouvement
func (e *Engine) hashMove(move AIMove, stateHash string) string {
	return fmt.Sprintf("%d-%d-%d-%s", move.Move.Piece, move.Move.Position.Row, move.Move.Position.Col, stateHash)
}

// evaluatePosition évalue une position de jeu
func (e *Engine) evaluatePosition(state GameState) int {
	if state.IsGameOver {
		switch state.Winner {
		case 1:
			return WIN_SCORE
		case -1:
			return LOSS_SCORE
		default:
			return DRAW_SCORE
		}
	}

	return 0
}

// Fonctions utilitaires pour min/max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
