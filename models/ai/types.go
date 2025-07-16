package ai

import "quarto/models/game"

type AIMove struct {
	Move          game.Move  // Le coup à jouer
	SelectedPiece game.Piece // La pièce sélectionnée pour le coup suivant
}

// GameState représente l'état du jeu pour l'IA
type GameState struct {
	Board           [4][4]game.Piece // Plateau de jeu, 0 = case vide
	AvailablePieces []game.Piece     // Pièces disponibles
	SelectedPiece   game.Piece       // Pièce sélectionnée pour placement
	IsGameOver      bool             // Jeu terminé
	Winner          int              // 0 = pas de gagnant, 1 = joueur 1, -1 = joueur 2
}

// Engine représente le moteur d'IA Quarto
type Engine struct {
	MaxDepth int                 // Profondeur maximale de recherche
	TT       *TranspositionTable // Table de transposition
}

// NewEngine crée un nouveau moteur d'IA
func NewEngine(maxDepth int) *Engine {

	return &Engine{
		MaxDepth: maxDepth,
		TT:       NewTranspositionTable(),
	}
}

const (
	WIN_SCORE  = 10000
	LOSS_SCORE = -10000
	DRAW_SCORE = 0
)

// TTFlag représente le type de valeur stockée dans la table de transposition
type TTFlag int

const (
	EXACT       TTFlag = iota // Valeur exacte
	LOWER_BOUND               // Borne inférieure (beta cutoff)
	UPPER_BOUND               // Borne supérieure (alpha cutoff)
)

// TTEntry représente une entrée dans la table de transposition
type TTEntry struct {
	Hash      string   // Hash de l'état du jeu
	Score     int      // Score évalué
	Depth     int      // Profondeur de recherche
	Flag      TTFlag   // Type de valeur
	BestMoves []AIMove // Liste des coups
}

// TranspositionTable représente la table de transposition
type TranspositionTable struct {
	entries map[string]*TTEntry
	hits    int // Compteur de cache hits pour les statistiques
	misses  int // Compteur de cache misses pour les statistiques
}

// NewTranspositionTable crée une nouvelle table de transposition
func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		entries: make(map[string]*TTEntry),
		hits:    0,
		misses:  0,
	}
}

// Store stocke une entrée dans la table de transposition
func (tt *TranspositionTable) Store(hash string, score int, depth int, flag TTFlag, bestMoves []AIMove) {
	tt.entries[hash] = &TTEntry{
		Hash:      hash,
		Score:     score,
		Depth:     depth,
		Flag:      flag,
		BestMoves: bestMoves,
	}
}

// Lookup cherche une entrée dans la table de transposition
func (tt *TranspositionTable) Lookup(hash string, depth int, alpha int, beta int) (entry *TTEntry, exists bool) {
	entry, exists = tt.entries[hash]
	if !exists {
		tt.misses++
		return
	}

	// L'entrée doit avoir été calculée à une profondeur au moins égale
	if entry.Depth < depth {
		tt.misses++
		return
	}
	exists = true
	tt.hits++
	return
}

// GetStats retourne les statistiques de la table de transposition
func (tt *TranspositionTable) GetStats() (int, int) {
	return tt.hits, tt.misses
}

// Clear vide la table de transposition
func (tt *TranspositionTable) Clear() {
	tt.entries = make(map[string]*TTEntry)
	tt.hits = 0
	tt.misses = 0
}
