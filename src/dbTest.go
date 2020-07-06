package main

import (
	"db"
	"fmt"
	"math/big"
)

func main() {
	var sample db.TrainingSample
	var boardIndex big.Int
	boardIndex.SetString("49954297991781395565", 10)
	//pi := [7]float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.4}
	sample.Pi = [7]float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.4}
	//	sample.Z = 1
	sample.P = []float32{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.4}
	sample.V = -0.77

	var db db.Database
	db.CreateTable()
	db.Insert(boardIndex, 123, sample, "PLAYER_1")
	fmt.Println(db.GetLastUid())
}
