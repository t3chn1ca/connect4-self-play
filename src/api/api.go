package api

const START_GAME = -1

type Game int

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
