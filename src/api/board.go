package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

const (
	maxY = 6
	maxX = 7
)

//Player constants
const (
	PLAYER_1 = -1
	PLAYER_2 = 1
)

const (
	GAME_IN_PROGRESS    = 0
	GAME_WON_PLAYER_1   = 1
	GAME_WON_PLAYER_2   = 2
	GAME_DRAW           = 3
	GAME_STATUS_UNKNOWN = -1
	GAME_FAULT_MOVE     = -2
)

//Connect4 Board class
type Connect4 struct {
	board             [maxY][maxX]int64
	playerMadeBadMove bool
	nextPlayerToMove  int64
	gameOver          bool
	reward            [2]int //one for each player
	boardIndex        string
}

var RandomNumGenerator *rand.Rand

func NewConnect4FromIndex(boardIndex string) *Connect4 {
	connect4 := new(Connect4)
	connect4.board, connect4.nextPlayerToMove = connect4.IndexToBoard(boardIndex)
	connect4.gameOver = false
	return connect4
}

func NewConnect4() *Connect4 {
	connect4 := new(Connect4)
	connect4.nextPlayerToMove = PLAYER_1
	connect4.gameOver = false
	source := rand.NewSource(Seed_for_rand)
	RandomNumGenerator = rand.New(source)
	fmt.Println("Initializing connect4")
	return connect4
}

func (b *Connect4) ResetBoard() {
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			b.board[y][x] = 0
		}
	}
	b.playerMadeBadMove = false
	b.nextPlayerToMove = PLAYER_1
	b.DumpBoard()
}

func (b Connect4) IsDone() bool {
	if b.GameStatus() == GAME_FAULT_MOVE || b.GameStatus() == GAME_DRAW || b.GameStatus() == GAME_WON_PLAYER_1 || b.GameStatus() == GAME_WON_PLAYER_2 {
		return true
	}

	return false

}

func (b *Connect4) IsGameDraw() bool {
	if b.GameStatus() == GAME_DRAW {
		return true
	} else {
		return false
	}
}

func (b *Connect4) GameStatus() int {

	if b.playerMadeBadMove {
		b.gameOver = true
		return GAME_FAULT_MOVE
	}
	for y_sliding := 0; y_sliding <= (maxY - 4); y_sliding++ {
		for x_sliding := 0; x_sliding <= (maxX - 4); x_sliding++ {
			for rowColSelector := 0; rowColSelector < 4; rowColSelector++ {
				sumHorizontal := int64(0)
				sumVertical := int64(0)
				sumDiagonal := int64(0)
				sumInvertedDiagonal := int64(0)
				for index := 0; index < 4; index++ {
					sumHorizontal += b.board[y_sliding+index][x_sliding+rowColSelector]
					sumVertical += b.board[y_sliding+rowColSelector][x_sliding+index]
					sumDiagonal += b.board[y_sliding+index][x_sliding+index]
					sumInvertedDiagonal += b.board[y_sliding+index][x_sliding+4-index-1]
					//fmt.Printf("board[%d][%d] SumVertical = %d \n", rowColSelector, index, sumVertical)
				}
				if sumHorizontal == -4 || sumVertical == -4 || sumDiagonal == -4 || sumInvertedDiagonal == -4 {
					b.gameOver = true
					b.reward[0] = 2
					b.reward[1] = -2
					return GAME_WON_PLAYER_1
				} else if sumHorizontal == 4 || sumVertical == 4 || sumDiagonal == 4 || sumInvertedDiagonal == 4 {
					b.reward[0] = -2
					b.reward[1] = 2
					b.gameOver = true
					return GAME_WON_PLAYER_2
				}
			}
		}
	}
	if len(b.GetValidMoves()) > 0 {
		return GAME_IN_PROGRESS
	}
	b.gameOver = true
	b.reward[0] = 1
	b.reward[1] = 1
	return GAME_DRAW

}

// Get all possible board states in the form of flat boards for NNApi
func (b Connect4) GetValidFlatBoardsFromPosition() ([]byte, [MAX_CHILD_NODES + 1]string) {
	//7 future board positions + 1 current board position
	var validFlatBoards [MAX_CHILD_NODES + 1][]int32
	var boardIndexes [MAX_CHILD_NODES + 1]string

	validMoves := b.GetValidMoves()

	for _, action := range validMoves {
		var boardTemp Connect4 = b
		//fmt.Printf("Playing Action: %d\n", action)
		boardTemp.PlayMove(action)
		validFlatBoards[action] = boardTemp.GetBoardFlatInt32()
		boardIndexes[action] = boardTemp.GetBoardIndex()
		//fmt.Printf("Board Index = %s\n", boardIndexes[action].String())
	}

	//Create flat boards for invalid moves so that the NN does not panic on null boards
	invalidMoves := b.GetInvalidMoves()
	var emptyBoard Connect4
	for _, action := range invalidMoves {
		validFlatBoards[action] = emptyBoard.GetBoardFlatInt32()
	}

	//Last entry is the current boardIndex
	validFlatBoards[MAX_CHILD_NODES] = b.GetBoardFlatInt32()
	boardIndexes[MAX_CHILD_NODES] = b.GetBoardIndex()

	jsonByte, err := json.Marshal(validFlatBoards)
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}
	return jsonByte, boardIndexes
}

// Get all possible board states from current position
func (b Connect4) GetValidBoardIndexesFromPosition() [MAX_CHILD_NODES]string {
	var validBoardIndexes [MAX_CHILD_NODES]string
	validMoves := b.GetValidMoves()

	for _, action := range validMoves {
		var boardTemp Connect4 = b
		boardTemp.PlayMove(action)
		validBoardIndexes[action] = boardTemp.GetBoardIndex()
	}
	return validBoardIndexes
}

func (b Connect4) GetInvalidMoves() []int {
	var invalidMoves []int
	for x := 0; x < maxX; x++ {
		if b.board[0][x] != 0 {
			invalidMoves = append(invalidMoves, x)
		}
	}
	return invalidMoves
}

func (b Connect4) GetValidMoves() []int {
	var validMoves []int
	for x := 0; x < maxX; x++ {
		if b.board[0][x] == 0 {
			validMoves = append(validMoves, x)
		}
	}
	return validMoves
}

func (b Connect4) GetPlayerToMove() int64 {
	return b.nextPlayerToMove
}

/* Returns 0 for PLAYER_1 and 1 for PLAYER_2*/
func (b Connect4) GetPlayerWhoJustMoved() int64 {
	if b.nextPlayerToMove == PLAYER_1 {
		return PLAYER_2
	} else {
		return PLAYER_1
	}

}

/* Player index is simply 0,1 notation for players for indexing arrays like rewards */
func GetPlayerIndex(player int64) int64 {
	if player == PLAYER_1 {
		return 0
	}
	if player == PLAYER_2 {
		return 1
	}
	fmt.Printf("Incorrect player index %d\n", player)
	panic("GetPlayerIndex():Incorrect player index ")

}

//flatten : Change a mxn array to 1, mxn array
func flatten(board [6][7]int64) []int64 {
	var flatBoard []int64
	for y := 0; y < maxY; y++ {
		flatBoard = append(flatBoard, board[y][0:cap(board[0])]...)
	}
	return flatBoard
}

func (b Connect4) IsGameOver() bool {
	return b.gameOver
}

func (b *Connect4) calculateBoardIndex() string {
	board := b.GetBoard()
	boardIndex := ""
	//fmt.Printf("Board: %v\n", board)
	for y := 0; y < 6; y++ {
		for x := 0; x < maxX; x++ {
			if board[y][x] == 0 {
				boardIndex += "0"
			}

			if board[y][x] == 1 {
				boardIndex += "1"
			}
			if board[y][x] == 2 {
				boardIndex += "2"
			}

		}
	}
	return boardIndex
}

//TODO: Fix board index
func (b Connect4) GetBoardIndex() string {
	return b.calculateBoardIndex()
}

func (b Connect4) IndexToBoard(boardIndex string) ([maxY][maxX]int64, int64) {
	//var board [maxY * maxX]int64
	var board [maxY][maxX]int64
	const arrayLen = maxY * maxX
	var countOfMoves int8
	var nextPlayerToMove int64

	for i := 0; i < (maxY * maxX); i++ {
		//fmt.Printf(" i = %d y = %d x = %d\n", i, i/maxX, i%maxX)
		var boardSlotValue int64
		if boardIndex[i] == '1' {
			boardSlotValue = -1
			countOfMoves++
		}
		if boardIndex[i] == '2' {
			boardSlotValue = 1
			countOfMoves++
		}

		if boardIndex[i] == '0' {
			boardSlotValue = 0
			countOfMoves++
		}
		board[i/maxX][i%maxX] = boardSlotValue
	}

	if countOfMoves%2 == 0 {
		//If its even any player could move next, will pick 1 for convinience
		nextPlayerToMove = PLAYER_1
	}
	if countOfMoves%2 == 1 {
		nextPlayerToMove = PLAYER_2
	}

	return board, nextPlayerToMove
}

func (b Connect4) GetBoardFlatInt32() []int32 {
	var boardInt32 []int32
	boardInt64 := b.GetBoardFlat()
	for _, cell := range boardInt64 {
		boardInt32 = append(boardInt32, int32(cell))
	}

	return boardInt32
}

/* Returns Int64 board for use with Big.Int */
func (b Connect4) GetBoardFlat() []int64 {
	//fmt.Println("%v", b.GetBoard())
	//fmt.Println("%v", flatten(b.GetBoard()))
	return flatten(b.GetBoard())
}

// Return board in format ready for easy convertion to ternary
func (b Connect4) GetBoard() [maxY][maxX]int64 {
	// Board with indexes , 0, 1, 2
	var boardForIndexFormat [maxY][maxX]int64
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			if b.board[y][x] == -1 {
				boardForIndexFormat[y][x] = 1
			} else if b.board[y][x] == 1 {
				boardForIndexFormat[y][x] = 2
			}
		}
	}
	return boardForIndexFormat // Return board in format ready for easy convertion to ternary
}

func (b Connect4) GetPlayerMadeBadMove() bool {
	return b.playerMadeBadMove
}

func (b Connect4) GetNextPlayerToMove() int64 {
	return b.nextPlayerToMove
}

func (b Connect4) GetGameOver() bool {
	return b.gameOver
}

func (b Connect4) GetReward() [2]int {
	return b.reward
}

//Get position information
func (b Connect4) getPos(x int, y int) int64 {
	return b.board[y][x]
}

func (b *Connect4) PlayMove(x int) int {
	if x < maxX && b.board[0][x] == 0 {
		//fmt.Printf("Player %s move = %d\n", b.PlayerToString(b.nextPlayerToMove), x)
		for y := maxY - 1; y >= 0; y-- {
			if b.board[y][x] == 0 {
				b.board[y][x] = b.nextPlayerToMove // Fill the slot for first vacant  opening on that column
				//fmt.Println(b.board[y][x])
				break
			}

		}

	} else {
		b.playerMadeBadMove = true
	}
	b.nextPlayerToMove *= -1

	//b.DumpBoard()
	return b.GameStatus()

}

/* DEBUGS */
func (b Connect4) PlayerToString(player int64) string {
	if player == PLAYER_1 {
		return "PLAYER_1(x)"
	} else {
		return "PLAYER_2(o)"
	}
}

func PlayerToString(player int64) string {
	if player == PLAYER_1 {
		return "PLAYER_1(x)"
	} else {
		return "PLAYER_2(o)"
	}
}

func (b Connect4) DumpBoard() {
	fmt.Println("-----------------------------------")
	fmt.Printf("Next player: %s\n", b.PlayerToString(b.nextPlayerToMove))
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			switch b.board[y][x] {
			case -1:
				fmt.Print(" x")
			case 1:
				fmt.Print(" o")
			case 0:
				fmt.Print(" -")
			}

		}
		fmt.Println()
	}
	fmt.Println("-----------------------------------")
	fmt.Println(" 0 1 2 3 4 5 6")
}
