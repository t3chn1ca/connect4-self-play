package main

import (
	"api"
	"db"
	"fmt"
	"log"
	"math/rand"
	"os"
)

/*
 * Questions:
   In MCTS backend samples, mctsWithQVal_Atl, the 0th game is draw but AvgV for both players is -0.5 & +0.5 why? Should it not be zero?
   A:

   Why is the end game not changing with change in iterations or random seed from golang ( it does change when seed is changed at nn)?
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

//Sample moves propablistically instead of picking best move
const PROPABLISTIC_SAMPLING_FALSE = false
const PROPABLISTIC_SAMPLING_TRUE = true
const PICK_BEST_MOVE_TRUE = false

//Do a search with traditional MCTS instead of using a backend nn
const MCTS_TREE_SEARCH_TRUE = true
const MCTS_TREE_SEARCH_FALSE = false
const NN_TREE_SEARCH_TRUE = false

const TRAINING_GAMES = 1000

//Debugs enabled
const DEBUGS_TRUE = true
const DEBUGS_FALSE = false

//First moves are randomized to create a rich set of diverse games for training, else the games are repetetive due to the NN always responding the same way
var QUARTER_OF_AVG_MOVES = 1 //Disable random first moves with -1 after some training

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
	var database db.Database
	database.CreateTable("mctsWithQVal_dec_20_2022")

	var selectedChild *api.Node
	var currRootNode *api.Node
	var mctsRootNode *api.Node
	var lastUid int32 = 0

	selectedChild = nil
	for {

		//Clear cache beggining of new training
		api.MonteCarloCacheInit()
		//In past all samples were used for training, now moved it to last 50 iterations
		lastUid = database.GetLastUid()
		log.Printf("LastUid= %d\n", lastUid)
		var game = api.NewConnect4()
		//Run MCTS first time to create root node and discard results of search , ie the selected child
		selectedChild = api.MonteCarloTreeSearch(MCTS_TREE_SEARCH_TRUE, PROPABLISTIC_SAMPLING_FALSE, game, MAX_MCTS_ITERATIONS, api.TRAIN_SERVER_PORT, selectedChild, DEBUGS_FALSE)
		//Make a copy of root node for caching, the idea being to pass the existing MCTS back for further iterations
		mctsRootNode = selectedChild.GetParent()
		//Set selectedChild to mctsRootNode so that an gameIteration can start from root
		selectedChild = mctsRootNode
		//Play over TRAINING_GAMES games with one another and create move samples
		for gameIteration := 0; gameIteration < TRAINING_GAMES; gameIteration++ {
			var move = 0
			game = api.NewConnect4()

			for {
				currRootNode = selectedChild
				selectedChild = api.MonteCarloTreeSearch(MCTS_TREE_SEARCH_TRUE, PROPABLISTIC_SAMPLING_FALSE, game, MAX_MCTS_ITERATIONS, api.TRAIN_SERVER_PORT, currRootNode, DEBUGS_FALSE)
				//Check we are 1/4 through the game for both players if not pick random
				if move < QUARTER_OF_AVG_MOVES*2 {
					//Let MCTS create child nodes before random selection
					//Pick child node from old parent
					selectedChild = currRootNode.GetRandomChildNode()
				}
				//fmt.Printf(api.DumpTree(mctsRootNode, 0))

				move++
				fmt.Printf("Move played by player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

				if currRootNode != mctsRootNode {
					fmt.Println(currRootNode.ToString())
					//fmt.Printf("Pi : %v\n", currRootNode.GetPi())
					sample := database.CreateSample(currRootNode.GetPi(false), currRootNode.GetP(), currRootNode.GetV(), currRootNode.GetQ())
					database.Insert(currRootNode.GetBoardIndex(), gameIteration, sample, game.PlayerToString(game.GetPlayerToMove()))
				}

				game.PlayMove(selectedChild.GetAction())
				fmt.Printf("gameIteration: %d\n", gameIteration)
				game.DumpBoard()
				//fmt.Print("Press 'Enter' to continue...")
				//fmt.Scanln()

				if game.IsGameOver() {
					database.UpdateWinner(lastUid, gameIteration, game.PlayerToString(game.GetPlayerWhoJustMoved()))
					api.MonteCarloCacheSyncToFile()
					println("GAME OVER")
					game.DumpBoard()
					selectedChild = mctsRootNode
					break
				}

			}
		}

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

}
