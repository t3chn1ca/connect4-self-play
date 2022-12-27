package main

import (
	"api"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const MAX_MCTS_ITERATIONS = 400

//Do a search with traditional MCTS instead of using a backend nn
const MCTS_TREE_SEARCH_TRUE = true
const MCTS_TREE_SEARCH_FALSE = false
const NN_TREE_SEARCH_TRUE = false

//Sample moves propablistically instead of picking best move
const PROPABLISTIC_SAMPLING_FALSE = false
const PROPABLISTIC_SAMPLING_TRUE = true
const PICK_BEST_MOVE_TRUE = false

//Debugs enabled
const DEBUGS_TRUE = true
const DEBUGS_FALSE = false

var selectedChild *api.Node = nil

const DO_FIRST_MOVE = "-1"

var game *api.Connect4

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getAiThinkingStatus(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Getting thinking status of AI")
	enableCors(&w)
	fmt.Fprintf(w, "{\"aiThinkingStatus\" : %.0f}", api.MctsIterationPercent)
}

func resetBoard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Resetting board")
	enableCors(&w)
	game = api.NewConnect4()
	selectedChild = nil
}

func playMove(w http.ResponseWriter, r *http.Request) {
	moves, _ := r.URL.Query()["move"]

	enableCors(&w)
	if moves[0] == DO_FIRST_MOVE {
		fmt.Printf("Doing first move")
		start := time.Now()
		selectedChild = api.MonteCarloTreeSearch(NN_TREE_SEARCH_TRUE, PROPABLISTIC_SAMPLING_FALSE, game, MAX_MCTS_ITERATIONS, api.TRAIN_SERVER_PORT, selectedChild, DEBUGS_FALSE)
		duration := time.Since(start)
		fmt.Printf("MCTS(%d) took %f long\n", MAX_MCTS_ITERATIONS, duration.Seconds())

	} else {

		humanMove, err := strconv.Atoi(moves[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Human move: %d\n", humanMove)
		if humanMove < 0 || humanMove > 6 {
			fmt.Printf("Incorrect Move!!!\nHuman move: ")
			return
		}
		//Play humans move
		game.PlayMove(humanMove)

		if game.IsGameOver() {
			println("GAME OVER")
			api.MonteCarloCacheSyncToFile()
			return
		}
		if selectedChild != nil {
			for _, child := range selectedChild.ChildNodes {
				if child.GetAction() == humanMove {
					selectedChild = child
				}
			}
		}
		//Let AI do its move
		start := time.Now()
		selectedChild = api.MonteCarloTreeSearch(NN_TREE_SEARCH_TRUE, PROPABLISTIC_SAMPLING_FALSE, game, MAX_MCTS_ITERATIONS, api.TRAIN_SERVER_PORT, selectedChild, DEBUGS_FALSE)
		duration := time.Since(start)

		fmt.Printf("MCTS(%d) took %f long\n", MAX_MCTS_ITERATIONS, duration.Seconds())
	}
	fmt.Printf("AI move: %d\n", selectedChild.GetAction())
	game.PlayMove(selectedChild.GetAction())
	game.DumpBoard()
	if game.IsGameOver() {
		println("GAME OVER")
		api.MonteCarloCacheSyncToFile()
	}
	fmt.Fprintf(w, "{\"move\" : %d}", selectedChild.GetAction())
}

func main() {

	//defer profile.Start().Stop()
	api.MonteCarloCacheInit()
	rand.Seed(int64(api.Seed_for_rand))
	game = api.NewConnect4()
	http.HandleFunc("/", playMove)
	http.HandleFunc("/resetBoard", resetBoard)
	http.HandleFunc("/getAiThinkingStatus", getAiThinkingStatus)
	fmt.Printf("Listening on port 8888\n")
	log.Fatal(http.ListenAndServe(":8888", nil))

}
