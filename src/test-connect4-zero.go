package main

import (
	"api"
	"fmt"
	"math/rand"
)

const PASS = true
const FAIL = false

const MAX_MCTS_ITERATIONS = 750

var failCount = 0
var runCount = 0

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

func testCaseWrapper(moves []int, validActions []int) bool {
	var game = api.NewConnect4()
	game = setupGame(game, moves)
	selectedChild := api.MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, nil, false, false)

	runCount++
	for _, action := range validActions {
		if action == selectedChild.GetAction() {
			return PASS
		}
	}
	failCount++
	return FAIL
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

	result := testCaseWrapper(moves, []int{2, 3})
	if result {
		fmt.Println("\n#####Test1 PASS")
	} else {
		fmt.Println("\n#####Test1 FAIL: Incorrect move")
	}
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
	result := testCaseWrapper(moves, []int{4, 1})
	if result {
		fmt.Println("\n#####Test2 PASS")
	} else {
		fmt.Println("\n#####Test2 FAIL: Incorrect move")
	}
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
	result := testCaseWrapper(moves, []int{4})
	if result {
		fmt.Println("\n#####Test3 PASS")
	} else {
		fmt.Println("\n#####Test3 FAIL: Incorrect move")
	}
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
	result := testCaseWrapper(moves, []int{0, 2, 4})
	if result {
		fmt.Println("\n#####Test4 PASS")
	} else {
		fmt.Println("\n#####Test4 FAIL: Incorrect move")
	}
	return result

}

func test5() bool {
	/*
		- - - - - - -
		- - - - - - -
		x - - - - - -
		x - - - - - -
		o - - - - - -
		x - - o o - -
		-------------------------------
		0 1 2 3 4 5 6
	*/

	moves := []int{0, 0, 0, 3, 0, 4}
	result := testCaseWrapper(moves, []int{2, 5})
	if result {
		fmt.Println("\n#####Test5 PASS")
	} else {
		fmt.Println("\n#####Test5 FAIL: Incorrect move")
	}
	return result

}

func test6() bool {
	/*
		- - - - - - -
		- - - - - - -
		x - - - - - -
		x - - - - - -
		o - - - - - -
		x - - o o - -
		-------------------------------
		0 1 2 3 4 5 6
	*/

	moves := []int{0, 0, 0, 3, 0, 4}
	result := testCaseWrapper(moves, []int{2, 5})
	if result {
		fmt.Println("#####Test6 PASS")
	} else {
		fmt.Println("#####Test6 FAIL: Incorrect move")
	}
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
	result := testCaseWrapper(moves, []int{1, 4})
	if result {
		fmt.Println("#####Test7 PASS")
	} else {
		fmt.Println("#####Test7 FAIL: Incorrect move")
	}
	return result

}

/*
 * TODO: Fix bug with UBC calculation, visit counts of all children in select is 0
 */
func main() {
	rand.Seed(int64(1234))

	const repeatTestCount = 1

	for i := 0; i < repeatTestCount; i++ {

		test1()
		test2()
		test3()
		test4()
		test5()
		test6()
		test7()
	}
	fmt.Printf("Failcount= %d/%d\n", failCount, runCount)
}
