package main

import (
	"api"
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

const MAX_MCTS_ITERATIONS = 500

var QUARTER_OF_AVG_MOVES = 2

func main() {
	//defer profile.Start().Stop()
	rand.Seed(int64(api.Seed_for_rand))
	var selectedChild *api.Node
	selectedChild = nil
	var randomWin = 0
	var connect4zeroWin = 0

	for iteration := 0; iteration < 100; iteration++ {

		var game = api.NewConnect4()
		selectedChild = nil
		//game.DumpBoard()
		//fmt.Scanln()
		for {
			//ChildNodesselectedChild = nil
			selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false, false)
			fmt.Printf("Move played by Player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

			game.PlayMove(selectedChild.GetAction())
			game.DumpBoard()
			if game.IsGameOver() {
				println("GAME OVER")
				connect4zeroWin++
				break
			}

			//Pick random move from available moves
			randomNode := selectedChild.GetRandomChildNode()
			randomMove := randomNode.GetAction()
			fmt.Printf("Random move: %d\n", randomMove)
			game.PlayMove(randomMove)
			game.DumpBoard()
			if game.IsGameOver() {
				println("GAME OVER")
				randomWin++
				break
			}
			for _, child := range selectedChild.ChildNodes {
				if child.GetAction() == randomMove {
					selectedChild = child
				}
			}

		}
	}

	fmt.Printf(" Summary Random win = %f  Connect4Zero win = %f \n", 100*float32(randomWin)/float32(randomWin+connect4zeroWin), 100*float32(connect4zeroWin)/float32(connect4zeroWin+randomWin))

}
