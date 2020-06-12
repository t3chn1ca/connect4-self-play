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
	Z  float64    //Real value of the state

	P [7]float64 //NN propablity for moves
	V float64    // NN value
}

func (db *Database) CreateTable() {
	dbCon, err := sql.Open("sqlite3", "./test.db")
	checkErr(err)

	stmt, err := dbCon.Prepare("create table IF NOT EXISTS training(uid integer PRIMARY KEY AUTOINCREMENT, boardIndex TEXT, playerToPlay TEXT, created datetime, iteration INTEGER, json TEXT);")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	db.conn = dbCon
}

func (db Database) Insert(boardIndex big.Int, iteration int, sample TrainingSample, playerToPlay string) {
	// insert
	now := time.Now()
	jsonString, err := json.Marshal(sample)
	checkErr(err)

	sql_str := fmt.Sprintf("INSERT INTO training(boardIndex, playerToPlay, created, iteration, json) values('%s', '%s', '%s', %d, '%s')", boardIndex.String(), playerToPlay, now.String(), iteration, jsonString)
	fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.conn.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}
