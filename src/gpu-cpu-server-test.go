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

const MAX_MCTS_ITERATIONS = 20

//const TRAIN_SERVER_PORT = 50051
//const BEST_SERVER_PORT = 50052

//First moves are randomized to create a rich set of diverse games for training, else the games are repetetive due to the NN always responding in same way
var QUARTER_OF_AVG_MOVES = 2 //Disable random first moves with -1 after some training

func main() {
	//defer profile.Start().Stop()
	rand.Seed(int64(api.Seed_for_rand))
	//var selectedChild *api.Node
	//var game = api.NewConnect4()
	/*
		fmt.Print("Testing GPU server")
		selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, TRAIN_SERVER_PORT, selectedChild, false, true)
		fmt.Print("GPU server test PASS\n")
		fmt.Print("Testing CPU server")
		selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, BEST_SERVER_PORT, selectedChild, false, true)
		fmt.Print("CPU server test PASS\n")
	*/
	results := Tournament(1)
	fmt.Printf("%v \n", results)
	api.NnSaveTrainedModelToBest()
}
