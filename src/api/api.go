package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const START_GAME = -1

type Game int

var connect4 = NewConnect4()

type ArgIn struct {
	X         int
	SessionId int //IF session id = -1 then start game
	Player    int
}

type ArgRet struct {
	Reward    int
	GameState int
	SessionId int
	Board     [maxY][maxX]int64
}

func (g Game) InitAllStates(r *http.Request, argIn *ArgIn, result *ArgRet) error {
	fmt.Println("==============Initializing all states======================")
	return nil
}

func (g *Game) ResetBoard(r *http.Request, argIn *ArgIn, result *ArgRet) error {
	fmt.Printf("ResetBoard: \n")
	connect4.ResetBoard()
	return nil
}

/*PlayMove API to play a move by any player */
func (g *Game) PlayMove(r *http.Request, argIn *ArgIn, result *ArgRet) error {
	//fmt.Printf("PlayMove: \n")
	//fmt.Printf("In client: playing move %d\n", argIn.X)

	result.SessionId = 1234
	status := connect4.PlayMove(argIn.X)
	result.GameState = status
	result.Reward = calculateReward(status)
	result.Board = connect4.GetBoard()

	return nil
}

/* API to read input from command line from server for Player 1 */
func (g *Game) GetMovePlayer(r *http.Request, argIn *ArgIn, result *ArgRet) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter move: ")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	s := strings.Split(text, " ")
	x, _ := strconv.Atoi(s[0])
	fmt.Printf("PLAYER_1 move: %d\n", x)
	status := connect4.PlayMove(x)
	result.GameState = status
	result.Board = connect4.GetBoard()
	result.Reward = calculateReward(status)
	result.SessionId = 1234

	return nil
}

func calculateReward(status int) int {
	//Always reward at end to make it propagate rewards backwards nicely in tictactoe RL code
	reward := 0
	switch status {
	case GAME_IN_PROGRESS:
		reward = 0
	case GAME_DRAW: //Half of wins
		reward = 8
	case GAME_WON_PLAYER_1:
		reward = -15
	case GAME_WON_PLAYER_2:
		reward = 15
	case GAME_FAULT_MOVE:
		reward = -20
	}
	return reward
}
