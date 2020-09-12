package main

import (
	"api"
	"fmt"
	"math/rand"
	"mcts"
)

//ref: https://github.com/jpbruneton/Alpha-Zero-algorithm-for-Connect-4-game
// Use MCTS depth as a reference to compare the NN ( Something like an ELO rating)

const MAX_MCTS_ITERATIONS = 1500

func main() {

	//defer profile.Start().Stop()

	rand.Seed(int64(api.Seed_for_rand))
	var selectedChild *mcts.Node
	selectedChild = nil

	var game = mcts.NewConnect4()

	//Run MCTS search to create child nodes for human to play
	selectedChild = mcts.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)
	selectedChild = selectedChild.GetParent()
	//var firstMoveDone bool = false
	for {

		fmt.Printf("Human move: ")
		//Let human play first
		var humanMove int
		n, _ := fmt.Scanf("%d", &humanMove)
		if humanMove < 0 || humanMove > 6 || n == 0 {
			fmt.Printf("Incorrect Move!!!\nHuman move: ")
		}
		game.PlayMove(humanMove)
		game.DumpBoard()
		if game.IsGameOver() {
			println("GAME OVER")
			break
		}
		if selectedChild != nil {
			for _, child := range selectedChild.ChildNodes {
				if child.GetAction() == humanMove {
					selectedChild = child
				}
			}
		}

		selectedChild = mcts.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)

		fmt.Printf("Move played by Player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())
		game.PlayMove(selectedChild.GetAction())
		game.DumpBoard()
		fmt.Printf("\a")
		if game.IsGameOver() {
			println("GAME OVER")
			break
		}

	}
}
