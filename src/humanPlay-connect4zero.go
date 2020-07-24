package main

import (
	"api"
	"fmt"
	"math/rand"
)

const MAX_MCTS_ITERATIONS = 1000

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

func main() {

	//defer profile.Start().Stop()

	rand.Seed(int64(api.Seed_for_rand))
	var selectedChild *api.Node
	selectedChild = nil

	var game = api.NewConnect4()
	//setupGame(game, []int{3, 3, 2, 4, 1, 0, 3, 3, 2, 2, 1, 1, 1, 0, 0, 2, 3, 2, 2, 3, 0, 0})
	//fmt.Scanln()
	for {

		selectedChild = api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, selectedChild, false, false)
		fmt.Printf("Move played by Player %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), selectedChild.GetAction())

		game.PlayMove(selectedChild.GetAction())
		game.DumpBoard()
		fmt.Printf("\a")
		if game.IsGameOver() {
			println("GAME OVER")
			break
		}
		fmt.Printf("Human move: ")
		var humanMove int
		fmt.Scanf("%d", &humanMove)
		game.PlayMove(humanMove)
		game.DumpBoard()
		if game.IsGameOver() {
			println("GAME OVER")
			break
		}
		if selectedChild != nil {
			for _, child := range selectedChild.ChildNodes {
				if child.GetAction() == humanMove {
					selectedChild = child
				}
			}
		}

	}
}
