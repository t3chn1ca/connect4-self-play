package main

import (
	"api"
	"math/big"
)

const PASS = true
const FAIL = false

const MAX_MCTS_ITERATIONS = 750

var failCount = 0
var runCount = 0

func main() {
	var boardIndex big.Int
	boardIndex.SetString("78903452661500241621", 10)
	var game = api.NewConnect4FromIndex(boardIndex, api.PLAYER_1)

	game.DumpBoard()

	//nnOut := api.NnForwardPass(game, api.TRAIN_SERVER_PORT)
	//fmt.Printf("Value = %v\n", nnOut.Value)
	//fmt.Printf("props %v\n", nnOut.P)

}
