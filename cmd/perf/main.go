package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"quarto/models/ai"
	"quarto/models/ai/stats"
	"quarto/models/game"
)

func generateRandomQuartoGame(numMoves int) ai.GameState {
	// Créer une nouvelle partie Quarto en mémoire uniquement
	g := game.InitializeGame(1, 2)
	var state ai.GameState

	// Simuler des mouvements aléatoires
	for i := 0; i < numMoves && g.Status == game.StatusPlaying; i++ {
		// Convertir en état pour l'IA
		state = ai.ConvertGameToState(g)

		// Vérifier si le jeu est terminé
		if state.IsGameOver {
			g.Status = game.StatusFinished
			switch state.Winner {
			case 1:
				g.Winner = 1
			case -1:
				g.Winner = 2
			default:
				g.Winner = 0
			}
			break
		}

		// Obtenir les mouvements valides
		validMoves := ai.GetValidMoves(state)
		if len(validMoves) == 0 {
			break // Plus de mouvements possibles
		}

		// Choisir un mouvement aléatoire
		randomMove := validMoves[rand.Intn(len(validMoves))]

		// Appliquer le mouvement directement à l'état
		state = state.ApplyMove(randomMove)
	}

	return state
}

func runQuartoBenchmark(depth int, numGames int, numMoves int, showStats bool) {
	totalStats := make(map[string]*stats.OperationStats)
	totalTime := time.Duration(0)
	validGames := 0

	fmt.Printf("Running Quarto benchmark with %d games (%d moves each, depth %d)...\n", numGames, numMoves, depth)

	quartoAI := ai.NewEngine(depth)

	for i := 0; i < numGames; i++ {
		g := generateRandomQuartoGame(numMoves)

		// Vérifier que le jeu n'est pas terminé
		if g.IsGameOver {
			fmt.Printf("Game %d: Already finished, skipping\n", i+1)
			continue
		}

		var gameStats *stats.PerformanceStats
		if showStats {
			gameStats = stats.NewPerformanceStats()
		}

		// Mesures mémoire avant
		var memBefore runtime.MemStats
		runtime.ReadMemStats(&memBefore)

		start := time.Now()
		if showStats {
			quartoAI.SearchWithStats(g, gameStats)
		} else {
			quartoAI.Search(g)
		}
		elapsed := time.Since(start)

		validGames++

		// Mesures mémoire après
		var memAfter runtime.MemStats
		runtime.ReadMemStats(&memAfter)

		totalTime += elapsed

		// Calcul usage mémoire
		allocDiff := memAfter.Alloc - memBefore.Alloc
		totalAllocDiff := memAfter.TotalAlloc - memBefore.TotalAlloc

		fmt.Printf("Game %d: Memory: %d KB allocated, %d KB total\n",
			i+1, allocDiff/1024, totalAllocDiff/1024)

		// Accumuler les statistiques
		if showStats && gameStats != nil {
			for opName, opStats := range gameStats.Operations {
				if totalStats[opName] == nil {
					totalStats[opName] = &stats.OperationStats{
						Count: 0,
						Time:  0,
						Cache: make(map[string]int64),
					}
				}
				totalStats[opName].Count += opStats.Count
				totalStats[opName].Time += opStats.Time

				for hash, hits := range opStats.Cache {
					totalStats[opName].Cache[hash] += hits
				}
			}
			gameStats.Reset()
		}
	}

	if validGames == 0 {
		fmt.Println("No valid games processed!")
		return
	}

	fmt.Printf("\n=== AVERAGE RESULTS OVER %d VALID GAMES ===\n", validGames)
	fmt.Printf("Average time: %v\n", totalTime/time.Duration(validGames))
	fmt.Printf("Total time: %v\n", totalTime)

	if showStats {
		fmt.Printf("\n=== PERFORMANCE STATISTICS ===\n")
		for opName, opStats := range totalStats {
			fmt.Printf("\nOperation: %s\n", opName)
			fmt.Printf("  Total count: %d\n", opStats.Count)
			fmt.Printf("  Average count per game: %.1f\n", float64(opStats.Count)/float64(validGames))
			fmt.Printf("  Total time: %v\n", opStats.Time)
			fmt.Printf("  Average time per game: %v\n", opStats.Time/time.Duration(validGames))
			fmt.Printf("  Average time per operation: %v\n", opStats.Time/time.Duration(opStats.Count))

			// Analyser les hits de cache
			if len(opStats.Cache) > 0 {
				type cacheStat struct {
					Hash string
					Hits int64
				}
				var cacheStatsSlice []cacheStat
				for hash, hits := range opStats.Cache {
					cacheStatsSlice = append(cacheStatsSlice, cacheStat{Hash: hash, Hits: hits})
				}

				sort.Slice(cacheStatsSlice, func(i, j int) bool {
					return cacheStatsSlice[i].Hits > cacheStatsSlice[j].Hits
				})

				fmt.Printf("  Unique cache entries: %d\n", len(cacheStatsSlice))
				fmt.Printf("  Top cache hits:\n")
				maxShow := min(5, len(cacheStatsSlice))
				for idx := 0; idx < maxShow; idx++ {
					cs := cacheStatsSlice[idx]
					avgHits := float64(cs.Hits) / float64(validGames)
					fmt.Printf("    %s: %d hits (%.1f avg per game)\n", cs.Hash, cs.Hits, avgHits)
				}

				// Calcul du taux de cache hit pour certaines opérations
				if opName == "node_visit" {
					cacheHitRate := float64(len(opStats.Cache)) / float64(opStats.Count) * 100
					fmt.Printf("  Cache hit rate: %.2f%%\n", cacheHitRate)
				}
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	depth := flag.Int("depth", 8, "Search depth for AI")
	showStats := flag.Bool("stats", false, "Show detailed performance stats")
	numGames := flag.Int("games", 1, "Number of games to test")
	numMoves := flag.Int("moves", 10, "Number of random moves for game generation")
	flag.Parse()

	// Initialiser le générateur aléatoire
	// Note: Depuis Go 1.20, plus besoin d'appeler rand.Seed

	runQuartoBenchmark(*depth, *numGames, *numMoves, *showStats)
}
