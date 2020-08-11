package api

import (
	"fmt"
	"math/big"
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
}

var RandomNumGenerator *rand.Rand

func NewConnect4FromIndex(boardIndex big.Int, nextPlayerToMove int64) *Connect4 {
	connect4 := new(Connect4)
	connect4.board = connect4.IndexToBoard(boardIndex)
	connect4.nextPlayerToMove = nextPlayerToMove
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

// Get all possible board states from current move
func (b Connect4) GetValidBoardIndexesFromMove() [MAX_CHILD_NODES]big.Int {
	var validBoardIndexes [MAX_CHILD_NODES]big.Int
	validMoves := b.GetValidMoves()

	for _, action := range validMoves {
		var boardTemp Connect4 = b
		boardTemp.PlayMove(action)
		validBoardIndexes[action] = boardTemp.GetBoardIndex()
	}
	return validBoardIndexes
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

//TODO: Fix board index
func (b Connect4) GetBoardIndex() big.Int {
	board := b.GetBoard()
	boardIndex := new(big.Int)
	posIndex := new(big.Int)
	//fmt.Printf("Board: %v\n", board)
	for y := 0; y < 6; y++ {
		for x := 0; x < maxX; x++ {
			var posVal big.Int
			//posVal = 3^posIndex
			posVal.Exp(big.NewInt(3), posIndex, big.NewInt(0)) // * float64(board[y][x])
			//fmt.Printf("posIndex: %d boardVal %d posVal: %d \n", posIndex, board[y][x], posVal)
			//posVal *= board[y][x]
			posVal.Mul(&posVal, big.NewInt(board[y][x]))
			//boardIndex += posVal
			boardIndex.Add(boardIndex, &posVal)
			//posIndex++
			posIndex.Add(posIndex, big.NewInt(1))
		}
	}
	return *boardIndex
}

func (b Connect4) IndexToBoard(boardIndex big.Int) [maxY][maxX]int64 {
	//var board [maxY * maxX]int64
	var board [maxY][maxX]int64
	const arrayLen = maxY * maxX
	temp := boardIndex
	ternary_ith_pos := new(big.Int)

	for i := 0; i < (maxY * maxX); i++ {
		//ternary_ith_pos = temp%3
		ternary_ith_pos = ternary_ith_pos.Mod(&temp, big.NewInt(3))
		//board[y][x] = ternary_ith_pos
		fmt.Printf(" i = %d\n", i)
		board[i/maxX][i%maxX] = ternary_ith_pos.Int64()
		//temp = temp/3
		temp.Div(&temp, big.NewInt(3))
		//if temp == 0
		if temp.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return board
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
