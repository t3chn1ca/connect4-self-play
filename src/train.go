package main

import (
	"api"
	"fmt"
	"log"
	"math/rand"
	"os"
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
   1. Integrate sql save, gameIteration, state, p, pi, v, z to db
   2. Load training data from python and train nn
*/

//On average there are 23 moves in connect-4 (ref:reddit.com/r/math/comments/1lo4od/how_many_games_of_connect4_there_are/)
//Create a randomizer which picks random moves in the first 25% (5.75) of the moves

const MAX_MCTS_ITERATIONS = 2000

const END_UID_INDEX = -1

const PROPABLISTIC_SAMPLING_FALSE = false
const PROPABLISTIC_SAMPLING_TRUE = true

//First moves are randomized to create a rich set of diverse games for training, else the games are repetetive due to the NN always responding the same way
var QUARTER_OF_AVG_MOVES = 2 //Disable random first moves with -1 after some training

func initLogging() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func main() {
	initLogging()
	log.Println("Starting connect4zero..")
	//defer profile.Start().Stop()
	rand.Seed(int64(api.Seed_for_rand))


	var lastUid int32 = 0
	// END_UID_INDEX for to index means to the end
	log.Printf("Starting to train from lastUid = %d\n", lastUid)
	api.TrainFromLastIteration(lastUid, END_UID_INDEX)

	log.Printf("Training complete\nTournament start\n")
	results := api.Tournament(20)

	fmt.Printf(" Results BEST PLAYER Wins = %f NEW PLAYER WINS = %f\n", 100*results.BestPlayerWins/results.TotalGames, 100*results.NewTrainedPlayerWins/results.TotalGames)
	if results.NewTrainedPlayerWins/results.TotalGames > 0.55 {
		//Newly trained player is better than the known best, replace the best
		log.Printf("New Trained player wins %f pc of games\n", results.NewTrainedPlayerWins*100/results.TotalGames)
		api.NnSaveTrainedModelToBest()
		log.Printf("Saving newly trained player to best model\n")
	} else {
		log.Printf("Best model beats newly trained model with win pc = %f\n", results.BestPlayerWins*100/results.TotalGames)
		api.NnLoadBestModelToGpu()
		log.Printf("Loading best model to GPU to continue new iterations\n")
	}


}
