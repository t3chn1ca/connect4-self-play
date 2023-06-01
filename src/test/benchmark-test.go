/* Benchmark and provide improvements based on known data from past tests */

package main

import (
	"api"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

// Do a search with traditional MCTS instead of using a backend nn
const MCTS_TREE_SEARCH_TRUE = true
const MCTS_TREE_SEARCH_FALSE = false
const NN_TREE_SEARCH_TRUE = false

// Sample moves propablistically instead of picking best move
const PROPABLISTIC_SAMPLING_FALSE = false
const PROPABLISTIC_SAMPLING_TRUE = true
const PICK_BEST_MOVE_TRUE = false

// Debugs enabled
const DEBUGS_TRUE = true
const DEBUGS_FALSE = false

const MAX_MCTS_ITERATIONS = 700
const SERVER_PORT = api.TRAIN_SERVER_PORT //BEST_SERVER_PORT //api.TRAIN_SERVER_PORT

func setupGame(game *api.Connect4, moves []int) *api.Connect4 {

	for _, move := range moves {
		//fmt.Println("===============================")
		//fmt.Printf("Player to move %s\n", game.PlayerToString(game.GetPlayerToMove()))
		game.PlayMove(move)
		//fmt.Printf("Player who just Moved %s\n", game.PlayerToString(game.GetPlayerWhoJustMoved()))
	}
	game.DumpBoard()
	return game

}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	//defer profile.Start().Stop()
	flag.Parse()
	if *cpuprofile != "" {
		fmt.Printf("CPU profiling started")
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	//rand.Seed(int64(api.Seed_for_rand))
	//DEBUG
	rand.Seed(int64(12132))
	var selectedChild *api.Node
	selectedChild = nil

	var game = api.NewConnect4()

	api.MonteCarloCacheInit()
	//setupGame(game, []int{3, 3, 2, 4, 1, 0, 3, 3, 2, 2, 1, 1, 1, 0, 0, 2, 3, 2, 2, 3, 0, 0})
	//fmt.Scanln()

	var humanMove = []int{3, 3, 3, 5}
	var humanMoveCount = 0
	benchmarkSecondsFromPastData := 29.821007
	measuredSeconds := 0.0
	for {

		start := time.Now()
		selectedChild = api.MonteCarloTreeSearch(NN_TREE_SEARCH_TRUE, PROPABLISTIC_SAMPLING_FALSE, game, MAX_MCTS_ITERATIONS, SERVER_PORT, selectedChild, DEBUGS_FALSE)

		fmt.Printf("Move played by Player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())
		duration := time.Since(start)

		fmt.Printf("MCTS (%d) took %f long\n", MAX_MCTS_ITERATIONS, duration.Seconds())
		measuredSeconds += duration.Seconds()
		game.PlayMove(selectedChild.GetAction())
		game.DumpBoard()
		fmt.Printf("\a")
		if game.IsGameOver() {
			println("GAME OVER")
			api.MonteCarloCacheSyncToFile()
			break
		}

		//DEBUG for profiling
		//go api.MonteCarloCacheSyncToFile()
		//break

		//DEBUG: do random human move
		//humanMove := rand.Intn(7)
		game.PlayMove(humanMove[humanMoveCount])

		game.DumpBoard()
		if game.IsGameOver() {
			println("GAME OVER")
			api.MonteCarloCacheSyncToFile()
			break
		}
		if selectedChild != nil {
			for _, child := range selectedChild.ChildNodes {
				if child.GetAction() == humanMove[humanMoveCount] {
					selectedChild = child
				}
			}
		}
		humanMoveCount++

	}
	fmt.Printf("\n\n----------------Benchmark Results------------------\n\n")
	fmt.Printf(" %% Deviation from Benchmark : \x1b[32m%f%% \x1b[0mimprovement\n", (benchmarkSecondsFromPastData-measuredSeconds)*100/benchmarkSecondsFromPastData)
	fmt.Printf("\n\n----------------*****************------------------\n")
}
