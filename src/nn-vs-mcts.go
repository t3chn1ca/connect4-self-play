package main

import (
	"api"
	"fmt"
	"math/rand"
	"mcts"
	"os"
	"strconv"
)

//ref: https://github.com/jpbruneton/Alpha-Zero-algorithm-for-Connect-4-game
// Use MCTS depth as a reference to compare the NN ( Something like an ELO rating)

const MAX_MCTS_ITERATIONS_NN = 1000
const MAX_MCTS_ITERATIONS_MCTS = 2000
const SERVER_PORT = api.TRAIN_SERVER_PORT
const QUARTER_OF_AVG_MOVES = 2

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

var MAX_TOURNAMENTS = 30

// TODO: Debug multiple instances of cache created
func main() {

	//defer profile.Start().Stop()
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		fmt.Printf("Setting Cpuct variable from arg = %s\n", argsWithoutProg[0])
		var err error
		api.C, err = strconv.ParseFloat(argsWithoutProg[0], 32)
		if err != nil {
			panic("Incorrect argument for Cpuct")
		}
	}

	rand.Seed(int64(api.Seed_for_rand))
	var selectedChildMcts *mcts.Node
	var selectedChildNn *api.Node
	var nnWinCount = 0.0
	var mctsWinCount = 0.0
	api.MonteCarloCacheInit()

	for tournament := 0; tournament < MAX_TOURNAMENTS; tournament++ {

		var gameMcts = mcts.NewConnect4()
		var gameNn = api.NewConnect4()

		//Build tree and revert selectedChild to root node
		//selectedChildMcts = mcts.MonteCarloTreeSearch(gameMcts, MAX_MCTS_ITERATIONS_MCTS, selectedChildMcts, false)
		//selectedChildMcts = selectedChildMcts.GetParent()

		//var firstMoveDone bool = false
		var move = 0
		for {

			selectedChildNn = api.MonteCarloTreeSearch(NN_TREE_SEARCH_TRUE, PICK_BEST_MOVE_TRUE, gameNn, MAX_MCTS_ITERATIONS_NN, SERVER_PORT, nil, DEBUGS_FALSE)

			if move < QUARTER_OF_AVG_MOVES*2 {
				//Let MCTS create child nodes before random selection
				//Pick child node from old parent
				selectedChildNn = selectedChildNn.GetParent().GetRandomChildNode()
				fmt.Println("RANDOM MOVE")
			}

			fmt.Printf("Move played by NN %s = %d\n", gameNn.PlayerToString(gameNn.GetPlayerToMove()), selectedChildNn.GetAction())

			gameMcts.PlayMove(selectedChildNn.GetAction())
			gameNn.PlayMove(selectedChildNn.GetAction())
			move++

			gameNn.DumpBoard()
			fmt.Println("----------------------------")
			fmt.Printf("TOURNAMENT : %d\n", tournament)
			fmt.Printf(" MCTS (depth:%d) Wins = %f/%d  NN (depth: %d) Wins = %f/%d \n", MAX_MCTS_ITERATIONS_MCTS, mctsWinCount, MAX_TOURNAMENTS, MAX_MCTS_ITERATIONS_NN, nnWinCount, MAX_TOURNAMENTS)
			fmt.Println("----------------------------")

			if gameNn.IsGameOver() {
				if gameNn.IsGameDraw() {
					nnWinCount += 0.5
					mctsWinCount += 0.5
					println("GAME DRAW")
					break
				}
				println("GAME OVER: NN Won")
				nnWinCount++
				break
			}
			/*
				//Feedback old MCTS tree to next move
				if selectedChildMcts != nil {
					for _, child := range selectedChildMcts.ChildNodes {
						if child.GetAction() == selectedChildNn.GetAction() {
							selectedChildMcts = child
						}
					}
				}*/

			selectedChildMcts = mcts.MonteCarloTreeSearch(gameMcts, MAX_MCTS_ITERATIONS_MCTS, nil, false)
			if move < QUARTER_OF_AVG_MOVES*2 {
				//Let MCTS create child nodes before random selection
				//Pick child node from old parent
				selectedChildMcts = selectedChildMcts.GetParent().GetRandomChildNode()
				fmt.Println("RANDOM MOVE")
			}

			fmt.Printf("Move played by MCTS %s = %d\n", gameMcts.PlayerToString(gameMcts.GetPlayerToMove()), selectedChildMcts.GetAction())

			gameNn.PlayMove(selectedChildMcts.GetAction())
			gameMcts.PlayMove(selectedChildMcts.GetAction())
			move++
			gameMcts.DumpBoard()

			if gameMcts.IsGameOver() {
				if gameNn.IsGameDraw() {
					nnWinCount += 0.5
					mctsWinCount += 0.5
					println("GAME DRAW")
					break
				}
				println("GAME OVER: MCTS Won")
				mctsWinCount++
				break
			}

		}
	}
	fmt.Printf(" MCTS Wins = %f/%d  NN Wins = %f/%d Cpuct = %f\n", mctsWinCount, MAX_TOURNAMENTS, nnWinCount, MAX_TOURNAMENTS, api.C)
}
