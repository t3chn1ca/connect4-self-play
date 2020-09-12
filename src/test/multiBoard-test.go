package main

import (
	"api"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
)

const PASS = true
const FAIL = false

const MAX_MCTS_ITERATIONS = 900

var failCount = 0
var runCount = 0

func setupGameWithBoardIndex(boardIndex big.Int) *api.Connect4 {
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

func coreTest(game *api.Connect4) bool {
	var result bool = PASS

	nnOutArray, boardIndexes := api.NnForwardPassMultiBoard(game, api.TRAIN_SERVER_PORT)
	for index := 0; index < api.MAX_CHILD_NODES; index++ {
		fmt.Printf("BoardIndex: " + boardIndexes[index].String() + "\n")
		fmt.Printf("Multi Value = %v\n", nnOutArray[index].Value)
		fmt.Printf("Multi props %v\n", nnOutArray[index].P)

		if boardIndexes[index].Cmp(big.NewInt(0)) != 0 { //Check if this move exist
			var gameTemp api.Connect4 = *game
			gameTemp.PlayMove(index)
			boardIndex := gameTemp.GetBoardIndex()
			if boardIndex.Cmp(&boardIndexes[index]) != 0 { //Verify the created boardIndex matches returned Index
				fmt.Println("BoardIndexes dont match!!")
				result = FAIL
			}
			nnOut := api.NnForwardPass(&gameTemp, api.TRAIN_SERVER_PORT)
			fmt.Printf("Single Value = %v\n", nnOut.Value)
			fmt.Printf("Single props %v\n", nnOut.P)

			if nnOut.Value-nnOutArray[index].Value > 0.01 {
				fmt.Printf("nnOut.Value = %v\n", nnOut.Value)
				fmt.Printf("nnOutArray[index].Value = %v\n", nnOutArray[index].Value)
				fmt.Printf("Value does not match for index" + strconv.Itoa(index) + "\n")
				result = FAIL
			}

			randomIndex := rand.Intn(api.MAX_CHILD_NODES)
			if nnOut.P[randomIndex]-nnOutArray[index].P[randomIndex] > 0.01 {
				fmt.Printf("nnOut.P[] = %v\n", nnOut.P[randomIndex])
				fmt.Printf("nnOutArray[index].P = %v\n", nnOutArray[index].P[randomIndex])
				fmt.Printf("P does not match for board index" + strconv.Itoa(index) + " and move index: " + strconv.Itoa(randomIndex) + "\n")
				result = FAIL
			}
		}
		fmt.Println()
	}
	fmt.Printf("BoardIndex: " + boardIndexes[api.MAX_CHILD_NODES].String() + "\n")
	fmt.Printf("Multi Value = %v\n", nnOutArray[api.MAX_CHILD_NODES].Value)
	fmt.Printf("Multi props %v\n", nnOutArray[api.MAX_CHILD_NODES].P)

	nnOut := api.NnForwardPass(game, api.TRAIN_SERVER_PORT)
	fmt.Printf("Single Value = %v\n", nnOut.Value)
	fmt.Printf("Single props %v\n", nnOut.P)

	if nnOut.Value-nnOutArray[api.MAX_CHILD_NODES].Value > 0.01 {
		fmt.Printf("nnOut.Value = %v\n", nnOut.Value)
		fmt.Printf("nnOutArray[7].Value = %v\n", nnOutArray[api.MAX_CHILD_NODES].Value)
		fmt.Printf("Value does not match for index" + strconv.Itoa(api.MAX_CHILD_NODES) + "\n")
		result = FAIL
	}

	randomIndex := rand.Intn(api.MAX_CHILD_NODES)
	if nnOut.P[randomIndex]-nnOutArray[api.MAX_CHILD_NODES].P[randomIndex] > 0.01 {
		fmt.Printf("nnOut.P[] = %v\n", nnOut.P[randomIndex])
		fmt.Printf("nnOutArray[7].P = %v\n", nnOutArray[api.MAX_CHILD_NODES].P[randomIndex])
		fmt.Printf("P does not match for board index" + strconv.Itoa(api.MAX_CHILD_NODES) + " and move index: " + strconv.Itoa(randomIndex) + "\n")
		result = FAIL
	}

	return result
}
func testCaseWrapperBoardIndex(boardIndex big.Int, validActions []int, testCaseName string) bool {
	fmt.Println("\n================================" + testCaseName + ": START======================")

	game := setupGameWithBoardIndex(boardIndex)
	//_ = setupGameWithBoardIndex(boardIndex)
	//return false
	result := coreTest(game)
	return result
}

func testCaseWrapper(moves []int, validActions []int, testCaseName string) bool {
	fmt.Println("\n================================" + testCaseName + ": START======================")

	var game = api.NewConnect4()
	game = setupGame(game, moves)

	result := coreTest(game)
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
		   	- - - o - - -
		   	x - - o o - x
		   	o o x o x - x
		   	o o o x x - o
		   	o x x x o x x
		-----------------------------------
	*/

	moves := []int{}

	result := testCaseWrapper(moves, []int{2, 3}, "test2")
	return result
}

func test3() bool {
	/*
		   	- - - - - - -
		   	- - - o - - -
		   	x - - o o - x
		   	o o x o x - x
		   	o o o x x - o
		   	o x x x o x x
		-----------------------------------
	*/
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
