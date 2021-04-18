package main

import (
	"api"
	"fmt"
	"math/rand"
)

const PASS = true
const FAIL = false

const MAX_MCTS_ITERATIONS = 900

var failCount = 0
var runCount = 0

func setupGameWithBoardIndex(boardIndex string) *api.Connect4 {
	var game = api.NewConnect4FromIndex(boardIndex)
	game.DumpBoard()
	return game
}

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

func testCaseWrapperBoardIndex(boardIndex string, validActions []int, testCaseName string) bool {
	fmt.Println("\n================================" + testCaseName + ": START======================")

	game := setupGameWithBoardIndex(boardIndex)
	var result bool = FAIL

	nnOut := api.NnForwardPass(game, api.TRAIN_SERVER_PORT)
	fmt.Printf("Value = %v\n", nnOut.Value)
	fmt.Printf("props %v\n", nnOut.P)

	selectedChild := api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, api.TRAIN_SERVER_PORT, nil, false, false)

	runCount++
	for _, action := range validActions {
		if action == selectedChild.GetAction() {

			result = PASS
			break
		}
	}

	if result == PASS {
		fmt.Println("\n#####" + testCaseName + " PASS")
	} else {
		failCount++
		fmt.Println("\n#####" + testCaseName + " FAIL")
	}
	fmt.Println("\n================================" + testCaseName + ": END======================")
	return result
}

func testCaseWrapper(moves []int, validActions []int, testCaseName string) bool {

	fmt.Println("\n================================" + testCaseName + ": START======================")
	var game = api.NewConnect4()
	game = setupGame(game, moves)
	var result bool = FAIL

	nnOut := api.NnForwardPass(game, api.TRAIN_SERVER_PORT)
	fmt.Printf("Value = %v\n", nnOut.Value)
	fmt.Printf("props %v\n", nnOut.P)

	selectedChild := api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, api.TRAIN_SERVER_PORT, nil, false, false)

	runCount++
	for _, action := range validActions {
		if action == selectedChild.GetAction() {

			result = PASS
			break
		}
	}

	if result == PASS {
		fmt.Println("\n#####" + testCaseName + " PASS")
	} else {
		failCount++
		fmt.Println("\n#####" + testCaseName + " FAIL")
	}
	fmt.Println("\n================================" + testCaseName + ": END======================")
	return result
}

func test1() bool {
	/*
		   	- - - - - - -
		   	- - - o - - -
		   	x - - o o - x
		   	o o x o x - x
		   	o o o x x - o
		   	o x x x o x x
		-----------------------------------
	*/

	moves := []int{1, 0, 2, 4, 3, 0, 6, 1, 3, 2, 4, 6, 2, 0, 4, 1, 6, 3, 0, 4, 5, 3, 6, 3} //, 6, 2}

	result := testCaseWrapper(moves, []int{2, 3}, "test1")
	return result
}

func test2() bool {
	/*
			- - - - - - -
			- - - - - - -
			- - - - - - -
			- - - - - - -
			- - o o - - -
			- - x x - - -
		-----------------------------------
	*/

	moves := []int{3, 3, 2, 2} //, 6, 2}
	result := testCaseWrapper(moves, []int{4, 1}, "test2")
	return result
}

func test3() bool {
	/*
		             o
		             V
			 - - - o - - -
			 - o - x - - -
			 - x o x - - -
			 - x x x - - -
			 - x o o x - x
			 o o o x x o o
			-----------------------------------
			 0 1 2 3 4 5 6
	*/

	moves := []int{3, 3, 3, 6, 3, 2, 3, 3, 6, 2, 2, 1, 1, 0, 1, 2, 1, 1, 4, 5, 4}
	result := testCaseWrapper(moves, []int{4}, "test3")
	return result

}
func test4() bool {
	/*
		- - - - - - -
		- - - - - - -
		- - - - - - -
		- - - - - - -
		- - - o - - -
		- x - x - o -
		-------------------------------
		0 1 2 3 4 5 6
	*/

	moves := []int{3, 3, 1, 5}
	result := testCaseWrapper(moves, []int{0, 2, 4}, "test4")
	return result

}

func test5() bool {
	/*
		- - - - - - -
		- - - - - - -
		- - - - - - -
		x - - - - - -
		o - - x - - -
		x - - o o - -
		-------------------------------
		0 1 2 3 4 5 6
	*/

	moves := []int{0, 0, 0, 3, 3, 4}
	result := testCaseWrapper(moves, []int{2, 5}, "test5")
	return result

}

func test6() bool {
	/*
		- - - - - - -
		- - - - - - -
		- - - - - - -
		x - - - - - -
		o - x - - - -
		x - o o - - -
		-------------------------------
		0 1 2 3 4 5 6
	*/

	moves := []int{0, 0, 0, 2, 2, 3}
	result := testCaseWrapper(moves, []int{1, 4, 5}, "test6")
	return result

}

func test7() bool {
	/*
		- - - - - - -
		- - - - - - -
		- - - - - - -
		- - - - - - -
		- - o o - - -
		x x o o x x o
		-------------------------------
		0 1 2 3 4 5 6
	*/

	moves := []int{0, 2, 1, 3, 4, 2, 5, 3}
	result := testCaseWrapper(moves, []int{1, 4}, "test7")
	return result

}

func test8() bool {
	/*
		To block 'o' from doubling
		V
		- - - o x x -
		- - x o o o -
		- - o x x x -
		- - o o x x -
		- o x o o o -
		x x x o o x -
		-----------------------------------
		0 1 2 3 4 5 6


	*/

	boardIndex := "000211000122200021110002211002122201112210"
	result := testCaseWrapperBoardIndex(boardIndex, []int{0}, "test8")
	return result

}

func test9() bool {
	/*
		- - - - - - -
		- - - - - - -
		- - - - - - -
		- - - - - - -
		- - - - - - -
		- o - o - x x
		-----------------------------------
		0 1 2 3 4 5 6
	*/

	boardIndex := "000000000000000000000000000000000000202011"
	result := testCaseWrapperBoardIndex(boardIndex, []int{2}, "test9")
	return result

}

/*
 * TODO: Fix bug with UBC calculation, visit counts of all children in select is 0
 */
func main() {
	rand.Seed(int64(1234))
	api.MonteCarloCacheInit()

	const repeatTestCount = 1

	for i := 0; i < repeatTestCount; i++ {

		test1()
		test2()
		test3()
		test4()
		test5()
		test6()
		test7()
		test8()
		test9()

	}
	fmt.Printf("Failcount= %d/%d\n", failCount, runCount)
}
