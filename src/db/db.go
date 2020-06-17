package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println("Panic!")
		panic(err)
	}
}

type Database struct {
	name string
	conn *sql.DB
}

type TrainingSample struct {
	Pi [7]float64 //MCTS based propablity of moves
	P  []float32  //NN propablity for moves
	V  float32    // NN value
}

// At end of game update all Z values for this iteration from ENDGAME result
func (db *Database) UpdateWinner(iteration int, playerWon string) {
	//Update Z for winner
	sql_str := fmt.Sprintf("UPDATE training SET z = 1 WHERE playerToMove = '%s' AND iteration = %d", playerWon, iteration)
	fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.conn.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	//Update Z for looser
	sql_str = fmt.Sprintf("UPDATE training SET z = -1 WHERE playerToMove != '%s' AND iteration = %d", playerWon, iteration)
	fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err = db.conn.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

func (db *Database) CreateTable() {
	dbCon, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)

	stmt, err := dbCon.Prepare("CREATE TABLE IF NOT EXISTS training(uid integer PRIMARY KEY AUTOINCREMENT, boardIndex TEXT, playerToMove TEXT, created datetime, iteration INTEGER, json TEXT, z INTEGER);")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	db.conn = dbCon
}

func (db Database) CreateSample(Pi [7]float64, P []float32, V float32) TrainingSample {
	var sample TrainingSample

	sample.Pi = Pi
	sample.P = P
	sample.V = V
	return sample
}

func (db Database) Insert(boardIndex big.Int, iteration int, sample TrainingSample, playerToMove string) {
	// insert
	now := time.Now()
	jsonString, err := json.Marshal(sample)
	checkErr(err)
	//Insert Z = 0 and other parameters
	sql_str := fmt.Sprintf("INSERT INTO training(boardIndex, playerToMove, created, iteration, json, z) values('%s', '%s', '%s', %d, '%s', 0)", boardIndex.String(), playerToMove, now.String(), iteration, jsonString)
	fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.conn.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}
