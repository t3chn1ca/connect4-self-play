package main

import (
	"api"
	"fmt"
)

func main() {
	results := api.Tournament(20)

	fmt.Printf(" Results BEST PLAYER Wins = %f NEW PLAYER WINS = %f\n", 100*results.BestPlayerWins/results.TotalGames, 100*results.NewTrainedPlayerWins/results.TotalGames)
	if results.NewTrainedPlayerWins/results.TotalGames > 0.55 {
		//Newly trained player is better than the known best, replace the best
		fmt.Printf("New Trained player wins %f pc of games\n", results.NewTrainedPlayerWins*100/results.TotalGames)
		api.NnSaveTrainedModelToBest()
		fmt.Printf("Saving newly trained player to best model\n")
	} else {
		fmt.Printf("Best model beats newly trained model with win pc = %f\n", results.BestPlayerWins*100/results.TotalGames)
	}

}
