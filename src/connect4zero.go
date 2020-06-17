package main

import (
	"api"
	"db"
	"fmt"
	"math/rand"
)

const MAX_MCTS_ITERATIONS = 50

/*
 * Questions: Why is the end game not changing with change in iterations or random seed from golang ( it does change when seed is changed at nn)?
   A: It changes when nn seed is changed, MCTS randomness does not effect outcome for large iterations as
	  the outcome is driven by state rather than initial parameters


   TODO:
   1. Integrate sql save, iteration, state, p, pi, v, z to db
   2. Load training data from python and train nn
*/
func main() {
	//defer profile.Start().Stop()
	rand.Seed(int64(12345))
	var database db.Database
	database.CreateTable()
	var selectedChild *api.Node
	var currRootNode *api.Node
	var mctsRootNode *api.Node
	selectedChild = nil

	for iteration := 0; iteration < 10; iteration++ {

		var game = api.NewConnect4()
		//game.DumpBoard()
		//fmt.Scanln()
		for {
			currRootNode = selectedChild
			if selectedChild == nil {
				selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)
				//Make a copy of root node for subsequent iterations, the idea being to pass the existing MCTS back for further iterations
				mctsRootNode = selectedChild.GetParent()
			} else {
				selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)
			}
			fmt.Printf("Move played by player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

			if currRootNode != nil && currRootNode != mctsRootNode {
				fmt.Println(currRootNode.ToString())
				fmt.Printf("Pi : %v\n", currRootNode.GetPi())
				sample := database.CreateSample(currRootNode.GetPi(), currRootNode.GetP(), currRootNode.GetV())
				database.Insert(currRootNode.GetBoardIndex(), iteration, sample, game.PlayerToString(game.GetPlayerToMove()))
			}

			game.PlayMove(selectedChild.GetAction())
			game.DumpBoard()
			if game.IsGameOver() {
				database.UpdateWinner(iteration, game.PlayerToString(game.GetPlayerWhoJustMoved()))
				println("GAME OVER")
				game.DumpBoard()
				selectedChild = mctsRootNode
				break
			}

		}
	}

	//fmt.Println("Test1 PASS")
	//fmt.Scanln()

}
