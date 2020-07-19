package main

import (
	"api"
	"db"
	"fmt"
	"math/rand"
)

/*
 * Questions: Why is the end game not changing with change in iterations or random seed from golang ( it does change when seed is changed at nn)?
   A: It changes when nn seed is changed, MCTS randomness does not effect outcome for large iterations as
	  the outcome is driven by state rather than initial parameters

   BUGS:
   1. Add dritchlet noise to make more explorations, make every game in the different iterations different
   2. Draws are not captured, fix that to update both winners in case of draw


   FIXED:
   1 z for player after first move is 0, it should be the end state of that game //Not seen now
   2. Server times out on long runs, change timeout grpc socket at server , caused by socket run out. clearing socket after use fixes it
   3. Node visit count != iterations <= FIXED: node does not need to set value when created only during update
   4. MCTS results and selfplay results dont tally, where it should <= Fixed after adding node for win
   5. The NN loss is stuck at 2.945, the NN architecture is not correct to learn the data


   TODO:
   1. Integrate sql save, iteration, state, p, pi, v, z to db
   2. Load training data from python and train nn
*/

//On average there are 23 moves in connect-4 (ref:reddit.com/r/math/comments/1lo4od/how_many_games_of_connect4_there_are/)
//Create a randomizer which picks random moves in the first 25% (5.75) of the moves

const MAX_MCTS_ITERATIONS = 150

var QUARTER_OF_AVG_MOVES = 2

func main() {
	//defer profile.Start().Stop()
	rand.Seed(int64(api.Seed_for_rand))
	var database db.Database
	database.CreateTable()
	var selectedChild *api.Node
	var currRootNode *api.Node
	var mctsRootNode *api.Node
	selectedChild = nil

	lastUid := database.GetLastUid()
	for {
		for iteration := 0; iteration < 50; iteration++ {

			var game = api.NewConnect4()
			//game.DumpBoard()
			//fmt.Scanln()
			var move = 0
			for {
				currRootNode = selectedChild
				if selectedChild == nil {

					selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false, true)
					//Make a copy of root node for subsequent iterations, the idea being to pass the existing MCTS back for further iterations
					mctsRootNode = selectedChild.GetParent()

					//Since first move is always < QUARTER_OF_AVG_MOVES
					selectedChild = mctsRootNode.GetRandomChildNode()
					//fmt.Print("Press 'Enter' to continue...")
					//fmt.Scanln()
				} else {
					//Check we are 1/4 through the game for both players if not pick random
					if move <= QUARTER_OF_AVG_MOVES*2 {
						//Cache old parent
						oldParent := selectedChild
						//Let MCTS create child nodes before random selection
						selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false, true)
						//Pick child node from old parent
						selectedChild = oldParent.GetRandomChildNode()
					} else {
						selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false, true)
					}
					//fmt.Printf(api.DumpTree(mctsRootNode, 0))
					//fmt.Print("Press 'Enter' to continue...")
					//fmt.Scanln()

				}
				move++
				fmt.Printf("Move played by player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

				if currRootNode != nil && currRootNode != mctsRootNode {
					fmt.Println(currRootNode.ToString())
					//fmt.Printf("Pi : %v\n", currRootNode.GetPi())
					sample := database.CreateSample(currRootNode.GetPi(false), currRootNode.GetP(), currRootNode.GetV())
					database.Insert(currRootNode.GetBoardIndex(), iteration, sample, game.PlayerToString(game.GetPlayerToMove()))
				}

				game.PlayMove(selectedChild.GetAction())
				fmt.Printf("Iteration: %d\n", iteration)
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
		api.TrainFromLastIteration(lastUid)
	}

}
