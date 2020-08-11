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

func (db *Database) GetLastUid() int32 {
	sql_str := fmt.Sprintf("SELECT uid from training ORDER BY uid DESC limit 1")
	//fmt.Printf(" SQL : %s\n", sql_str)
	rows, err := db.conn.Query(sql_str)
	checkErr(err)

	var uid int32
	for rows.Next() {
		rows.Scan(&uid)
	}
	return uid
}

// At end of game update all Z values for this iteration from ENDGAME result
func (db *Database) UpdateWinner(lastUid int32, iteration int, playerWon string) {
	//Update Z for winner for current iteration only ( not past trainings iterations )
	sql_str := fmt.Sprintf("UPDATE training SET z = 1 WHERE uid > %d AND playerToMove = '%s' AND iteration = %d", lastUid, playerWon, iteration)
	//fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.conn.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	//Update Z for looser for current iteration only ( not past trainings iterations )
	sql_str = fmt.Sprintf("UPDATE training SET z = -1 WHERE uid > %d AND playerToMove != '%s' AND iteration = %d", lastUid, playerWon, iteration)
	//fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err = db.conn.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

func (db *Database) CreateTable(tableName string) {
	dbCon, err := sql.Open("sqlite3", "./"+tableName+".db")
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
	//fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.conn.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}
