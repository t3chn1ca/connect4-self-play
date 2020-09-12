package main

import (
	"api"
	"fmt"
	"math/big"
)

const PASS = true
const FAIL = false

const MAX_MCTS_ITERATIONS = 900

var failCount = 0
var runCount = 0

func setupGameWithBoardIndex(boardIndex big.Int) *api.Connect4 {
	var game = api.NewConnect4FromIndex(boardIndex)
	game.DumpBoard()
	boardFlat := game.GetBoardFlat()
	fmt.Printf("Board Flat: %v\n", boardFlat)
	return game
}

func testCaseWrapperBoardIndex(boardIndex big.Int, validActions []int, testCaseName string) bool {
	fmt.Println("\n================================" + testCaseName + ": START======================")

	_ = setupGameWithBoardIndex(boardIndex)
	return true

}

func test3() bool {

	var boardIndex big.Int
	boardIndex.SetString("1352087347153435893", 10)
	result := testCaseWrapperBoardIndex(boardIndex, []int{2, 3}, "test3")

	return result
}

/*
 * TODO: Fix bug with UBC calculation, visit counts of all children in select is 0
 */
func main() {
	//rand.Seed(int64(1234))
	api.MonteCarloCacheInit()

	const repeatTestCount = 1

	//test1()
	//test2()
	test3()
	fmt.Printf("Failcount= %d/%d\n", failCount, runCount)
}
