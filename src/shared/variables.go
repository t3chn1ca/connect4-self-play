package shared

import "sync"

type NNOut struct {
	Value float32
	P     []float32
}

const maxX = 7
const maxY = 6

func GetBoardIndexMirror(boardIndex string) string {
	boardIndexRune := []rune(boardIndex)
	for y := 0; y < maxY; y++ {
		for x := 0; x <= (maxX/2)-1; x++ { // only check from 0-2
			temp := boardIndexRune[y*maxX+x]
			//fmt.Printf(" boardIndexRune[%d] = boardIndexRune[%d] \n", y*maxX+x, ((y+1)*maxX)-x-1)
			boardIndexRune[y*maxX+x] = boardIndexRune[((y+1)*maxX)-x-1]
			boardIndexRune[((y+1)*maxX)-x-1] = temp
		}
	}
	return string(boardIndexRune)
}

var Wg sync.WaitGroup
