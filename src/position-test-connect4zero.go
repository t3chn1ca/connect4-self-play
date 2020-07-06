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
   3. MCTS results and selfplay results dont tally, where it should
   4. Value returned from MCTS is not correct, First board is a sure loose but returns 0.74??

   FIXED:
   1 z for player after first move is 0, it should be the end state of that game //Not seen now
   2. Server times out on long runs, change timeout grpc socket at server , caused by socket run out. clearing socket after use fixes it

   TODO:
   1. Integrate sql save, iteration, state, p, pi, v, z to db
   2. Load training data from python and train nn
*/

//On average there are 23 moves in connect-4 (ref:reddit.com/r/math/comments/1lo4od/how_many_games_of_connect4_there_are/)
//Create a randomizer which picks random moves in the first 25% (5.75) of the moves
var QUARTER_OF_AVG_MOVES = 4

func setupGame(game *api.Connect4, moves []int) *api.Connect4 {

	for _, move := range moves {
		fmt.Println("===============================")
		fmt.Printf("Player to move %s\n", game.PlayerToString(game.GetPlayerToMove()))
		game.PlayMove(move)
		fmt.Printf("Player who just Moved %s\n", game.PlayerToString(game.GetPlayerWhoJustMoved()))
	}
	return game

}

const MAX_MCTS_ITERATIONS = 500

func main() {
	//defer profile.Start().Stop()
	rand.Seed(int64(api.Seed_for_rand))
	var database db.Database
	database.CreateTable()
	var selectedChild *api.Node
	var currRootNode *api.Node
	var mctsRootNode *api.Node
	selectedChild = nil

	for iteration := 0; iteration < 1; iteration++ {

		var game = api.NewConnect4()
		/*      V V
		   	- - - - - - -
		   	- - - o - - -
		   	x - - o o - x
		   	o o x o x - x
		   	o o o x x - o
		   	o x x x o x x
		-----------------------------------
			0 1 2 3 4 5 6
			ANS: 2,3
		*/
		moves := []int{1, 0, 2, 4, 3, 0, 6, 1, 3, 2, 4, 6, 2, 0, 4, 1, 6, 3, 0, 4, 5, 3, 6, 3} //, 6, 2}
		/*    V     V
			- - - - - - -
			- - - - - - -
			- - - - - - -
			- - - - - - -
			- - o o - - -
			- - x x - - -
		-----------------------------------
			0 1 2 3 4 5 6

			ANS: 1,4
		*/
		//moves := []int{3, 3, 2, 2} //, 6, 2}

		/*
			             o
			             V
				 - - - o - - -
				 - o - x - - -
				 - x o x - - -
				 - x x x - - -
				 - x o o x - x
				 o o o x x o o
				-----------------------------------
				 0 1 2 3 4 5 6
				 ANS: 4

		*/
		//moves := []int{3, 3, 3, 6, 3, 2, 3, 3, 6, 2, 2, 1, 1, 0, 1, 2, 1, 1, 4, 5, 4, 4}
		game = setupGame(game, moves)
		//game.DumpBoard()
		//fmt.Scanln()

		for {
			currRootNode = selectedChild
			if selectedChild == nil {

				selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)
				//Make a copy of root node for subsequent iterations, the idea being to pass the existing MCTS back for further iterations
				mctsRootNode = selectedChild.GetParent()
				fmt.Printf(api.DumpTree(mctsRootNode, 0))
				fmt.Print("Press 'Enter' to continue...")
				fmt.Scanln()
				//Since first move is always < QUARTER_OF_AVG_MOVES
				//selectedChild = mctsRootNode.GetRandomChildNode()

			} else {
				selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false)
				fmt.Printf(api.DumpTree(mctsRootNode, 0))
				fmt.Print("Press 'Enter' to continue...")
				fmt.Scanln()

			}
			fmt.Printf("Move played by player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

			if currRootNode != nil && currRootNode != mctsRootNode {
				fmt.Println(currRootNode.ToString())
				//fmt.Printf("Pi : %v\n", currRootNode.GetPi())
				sample := database.CreateSample(currRootNode.GetPi(), currRootNode.GetP(), currRootNode.GetV())
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

}
