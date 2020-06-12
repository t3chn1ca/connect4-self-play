package main

import (
	"api"
	"fmt"
	"math/rand"

	"github.com/pkg/profile"
)

const MAX_MCTS_ITERATIONS = 500

/*
 * Questions: Why is the end game not changing with change in iterations or random seed?
   A: It changes when nn seed is changed, MCTS randomness does not effect outcome for large iterations as
	  the outcome is driven by state rather than initial parameters

   TODO:
   1. Integrate sql save, iteration, state, p, pi, v, z to db
   2. Load training data from python and train nn
*/
func main() {
	defer profile.Start().Stop()
	rand.Seed(int64(1234))
	//var iteration = 0
	var game = api.NewConnect4()

	game.DumpBoard()
	//fmt.Scanln()
	var selectedChild *api.Node
	var currRootNode *api.Node
	selectedChild = nil
	for {
		currRootNode = selectedChild
		selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)
		fmt.Printf("Move played by player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

		if currRootNode != nil {
			fmt.Println(currRootNode.ToString())
			fmt.Printf("Pi : %v\n", currRootNode.GetPi())
		}

		game.PlayMove(selectedChild.GetAction())

		if game.IsGameOver() {
			println("GAME OVER")
			game.DumpBoard()
			break
		}

	}

	//fmt.Println("Test1 PASS")
	//fmt.Scanln()

}
