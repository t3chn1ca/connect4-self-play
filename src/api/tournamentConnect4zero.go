package api

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

//On average there are 23 moves in connect-4 (ref:reddit.com/r/math/comments/1lo4od/how_many_games_of_connect4_there_are/)
//Create a randomizer which picks random moves in the first 25% (5.75) of the moves

const MAX_MCTS_ITERATIONS_TOURNAMENT = 1000

const TRAIN_SERVER_PORT = 50051
const BEST_SERVER_PORT = 50052
const RANDOMIZE_MOVE_COUNT = 2

type Results struct {
	TotalGames           float32
	BestPlayerWins       float32
	NewTrainedPlayerWins float32
}

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

/*
 TODO: Complete tournament code which returns results

  BUGS: The cpu and gpu servers are not able to coexist?
*/
func Tournament(maxTournaments int) Results {
	var selectedChild *Node
	var bestPlayerWins float32
	var newTrainedPlayerWins float32
	var server int
	var move = 0
	var playerStr string

	log.Println("Starting tournament()")

	fmt.Printf("Loading BestNN to cpu\n")
	NnLoadCpuModel()
	//wait here for the server to load
	fmt.Printf("Waiting for model to load to cpu\n")
	time.Sleep(5000 * time.Millisecond)
	fmt.Printf("Ready to play\n")

	fmt.Println("Starting Tournament..")
	for iteration := 0; iteration < maxTournaments; iteration++ {
		var bestPlayerMove bool
		//Random bool generator
		bestPlayerMove = (bool)((rand.Intn(99) % 2) == 1)
		//bestPlayerMove = false //DEBUG
		game := NewConnect4()

		for {

			if bestPlayerMove == true {
				server = BEST_SERVER_PORT
				playerStr = "BEST Player"
			} else {
				server = TRAIN_SERVER_PORT
				playerStr = "NEWLY trained Player"
			}

			selectedChild = nil // Start with new tree for every move to make search similar for both parties
			selectedChild = MonteCarloTreeSearch(NN_TREE_SEARCH_TRUE, PROPABLISTIC_SAMPLING_FALSE, game, MAX_MCTS_ITERATIONS_TOURNAMENT, server, nil, DEBUGS_FALSE)
			move++
			if move < RANDOMIZE_MOVE_COUNT*2 { // Randomize the first n moves for variations in games
				//Let MCTS create child nodes before random selection
				//Pick child node from old parent
				selectedChild = selectedChild.GetParentNode().GetRandomChildNode()
			}
			fmt.Printf("Move played by %s player %s = %d\n", playerStr, game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

			game.PlayMove(selectedChild.GetAction())
			fmt.Printf("Tournament : %d\n", iteration)
			game.DumpBoard()
			//fmt.Print("Press 'Enter' to continue...")
			//fmt.Scanln()

			if game.IsGameOver() {
				if game.IsGameDraw() != true {
					if bestPlayerMove == true {
						bestPlayerWins++
					} else {
						newTrainedPlayerWins++
					}
				} else { //DRAW
					bestPlayerWins += 0.5
					newTrainedPlayerWins += 0.5
				}

				println("GAME OVER")
				game.DumpBoard()
				break
			}
			bestPlayerMove = !bestPlayerMove
		}
	}

	var results Results

	results.BestPlayerWins = bestPlayerWins
	results.NewTrainedPlayerWins = newTrainedPlayerWins
	results.TotalGames = bestPlayerWins + newTrainedPlayerWins
	NnStopCpuModel()
	return results

}
