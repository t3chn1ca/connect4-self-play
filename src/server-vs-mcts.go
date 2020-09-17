package main

import (
	"api"
	"fmt"
	"log"
	"math/rand"
	"mcts"
	"net/http"
	"strconv"
)

const MAX_MCTS_ITERATIONS = 6000

var selectedChild *mcts.Node = nil

const DO_FIRST_MOVE = "-1"

var game *mcts.Connect4

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func resetBoard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Resetting board")
	enableCors(&w)
	game = mcts.NewConnect4()
	selectedChild = nil
}

func playMove(w http.ResponseWriter, r *http.Request) {
	moves, _ := r.URL.Query()["move"]

	enableCors(&w)
	if moves[0] == DO_FIRST_MOVE {
		fmt.Printf("Doing first move")
		selectedChild = mcts.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)

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
		selectedChild = mcts.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)

	}
	fmt.Printf("MCTS move: %d\n", selectedChild.GetAction())
	game.PlayMove(selectedChild.GetAction())
	game.DumpBoard()
	fmt.Fprintf(w, "{\"move\" : %d}", selectedChild.GetAction())
}

func main() {

	//defer profile.Start().Stop()
	api.MonteCarloCacheInit()
	rand.Seed(int64(api.Seed_for_rand))
	game = mcts.NewConnect4()
	http.HandleFunc("/", playMove)
	http.HandleFunc("/resetBoard", resetBoard)
	log.Fatal(http.ListenAndServe(":8888", nil))

}
