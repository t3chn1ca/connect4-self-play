package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"shared"
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
	name         string
	connTraining *sql.DB
	connCache    *sql.DB
}

type TrainingSample struct {
	Pi   [7]float64 //MCTS based propablity of moves
	P    []float32  //NN propablity for moves
	V    float32    // NN value
	AvgV float32    //Avg Value of all children under this node
}

func (db *Database) GetLastUid() int32 {
	sql_str := fmt.Sprintf("SELECT uid from training ORDER BY uid DESC limit 1")
	//fmt.Printf(" SQL : %s\n", sql_str)
	rows, err := db.connTraining.Query(sql_str)
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
	sql_str := fmt.Sprintf("UPDATE training SET z = -1 WHERE uid > %d AND playerToMove = '%s' AND iteration = %d", lastUid, playerWon, iteration)
	//fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.connTraining.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	//Update Z for looser for current iteration only ( not past trainings iterations )
	sql_str = fmt.Sprintf("UPDATE training SET z = 1 WHERE uid > %d AND playerToMove != '%s' AND iteration = %d", lastUid, playerWon, iteration)
	//fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err = db.connTraining.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

func (db *Database) CreateTable(dbFile string) {
	dbCon, err := sql.Open("sqlite3", "./"+dbFile+".db")
	checkErr(err)

	stmt, err := dbCon.Prepare("CREATE TABLE IF NOT EXISTS training(uid integer PRIMARY KEY AUTOINCREMENT, boardIndex TEXT, playerToMove TEXT, created datetime, iteration INTEGER, json TEXT, z INTEGER);")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	db.connTraining = dbCon
}

func (db Database) CreateSample(Pi [7]float64, P []float32, V float32, AvgV float32) TrainingSample {
	var sample TrainingSample

	sample.Pi = Pi
	sample.P = P
	sample.V = V
	sample.AvgV = AvgV
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
	stmt, err := db.connTraining.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

/*    CACHE Table */

func (db *Database) ClearCache() {
	sql_str := "DELETE from cache"
	//fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.connCache.Prepare(sql_str)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
}

func (db *Database) CreateCacheTable(dbFile string) {

	dbCon, err := sql.Open("sqlite3", ":memory:") //""./"+dbFile+".db")
	checkErr(err)

	stmt, err := dbCon.Prepare("CREATE TABLE IF NOT EXISTS cache(boardIndex TEXT PRIMARY KEY, json TEXT);")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	db.connCache = dbCon
}

func (db *Database) InsertCacheEntry(boardIndex big.Int, nnOut shared.NNOut) {

	/*//DEBUG
	return
	//DEBUG: END*/
	jsonString, err := json.Marshal(nnOut)
	checkErr(err)
	sql_str := fmt.Sprintf("INSERT INTO cache(boardIndex, json) values('%s', '%s')", boardIndex.String(), jsonString)
	//fmt.Printf(" SQL : %s\n", sql_str)
	stmt, err := db.connCache.Prepare(sql_str)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

func (db *Database) GetCacheEntry(boardIndex big.Int) (bool, shared.NNOut) {

	var nnOut shared.NNOut
	var ret bool = false

	/*//DEBUG
	return false, nnOut
	//DEBUG: END*/
	sql_str := fmt.Sprintf("SELECT json from cache where boardIndex = '%s'", boardIndex.String())
	//fmt.Printf(" SQL : %s\n", sql_str)
	rows, err := db.connCache.Query(sql_str)
	checkErr(err)

	var jsonString string = ""
	for rows.Next() {
		ret = true
		rows.Scan(&jsonString)
		//fmt.Println("Json:" + jsonString)
		err = json.Unmarshal([]byte(jsonString), &nnOut)
		checkErr(err)
	}

	return ret, nnOut
}
