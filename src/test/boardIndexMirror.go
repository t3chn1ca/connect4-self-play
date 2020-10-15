package main

import (
	"api"
	"shared"
)

const PASS = true
const FAIL = false

const MAX_MCTS_ITERATIONS = 750

var failCount = 0
var runCount = 0

func main() {

	boardIndex := "000000000021021002201221210122211022111211"
	var game = api.NewConnect4FromIndex(boardIndex)

	game.DumpBoard()

	boardIndex = shared.GetBoardIndexMirror("000000000021021002201221210122211022111211")
	game = api.NewConnect4FromIndex(boardIndex)
	game.DumpBoard()
	//nnOut := api.NnForwardPass(game, api.TRAIN_SERVER_PORT)
	//fmt.Printf("Value = %v\n", nnOut.Value)
	//fmt.Printf("props %v\n", nnOut.P)

}
