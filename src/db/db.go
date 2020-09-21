package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"shared"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
	if err == sql.ErrNoRows {
		return
	}

	if err != nil {
		fmt.Println("Panic!")
		//panic(err)
		os.Exit(-1)
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

func (db Database) Insert(boardIndex string, iteration int, sample TrainingSample, playerToMove string) {
	// insert
	now := time.Now()
	jsonString, err := json.Marshal(sample)
	checkErr(err)
	//Insert Z = 0 and other parameters
	sql_str := fmt.Sprintf("INSERT INTO training(boardIndex, playerToMove, created, iteration, json, z) values('%s', '%s', '%s', %d, '%s', 0)", boardIndex, playerToMove, now.String(), iteration, jsonString)
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
	_, err := db.connCache.Exec(sql_str)
	checkErr(err)
}

func (db *Database) CreateCacheTable(dbFile string) {

	dbCon, err := sql.Open("sqlite3", "file::memory:") //?cache=shared&_sync=0&_mutex=no") //":memory:?_sync=0&_mutex=no") //:?_mutex=no") //":memory:?_sync=0") //""./"+dbFile+".db")
	checkErr(err)
	fmt.Println("Creating cache table..")
	//stmt, err := dbCon.Prepare("CREATE TABLE IF NOT EXISTS cache(boardIndex TEXT PRIMARY KEY, json TEXT);")
	//checkErr(err)
	sql_str := "CREATE TABLE IF NOT EXISTS cache(boardIndex TEXT PRIMARY KEY, json TEXT);"
	_, err = dbCon.Exec(sql_str)
	checkErr(err)
	db.connCache = dbCon
}

func (db *Database) InsertCacheEntries(boardIndexes [8]string, nnOutArray [8]shared.NNOut) {

	for index, boardIndex := range boardIndexes {
		//retCode, _ := db.GetCacheEntry(boardIndex)
		//if retCode == false {
		db.InsertCacheEntry(boardIndex, nnOutArray[index])
		//}
		//	} else {
		//fmt.Println("Skipping boardIndex as it already exists " + boardIndex.String())
		//	}
	}
}

func (db *Database) InsertCacheEntry(boardIndex string, nnOut shared.NNOut) {
	/*//DEBUG
	return
	//DEBUG: END*/
	//if db.IsEntryPresent(boardIndex) {
	//fmt.Printf("BoardIndex %s exists in cache, return\n", boardIndex.String())
	//	return
	//}

	jsonString, err := json.Marshal(nnOut)
	checkErr(err)
	sql_str := fmt.Sprintf("INSERT INTO cache(boardIndex, json) values('%s', '%s')", boardIndex, string(jsonString))
	//fmt.Printf(" SQL : %s\n", sql_str)

	_, err = db.connCache.Exec(sql_str)
	//checkErr(err)
}

func (db Database) IsEntryPresent(boardIndex string) bool {
	/*//DEBUG
	return false, nnOut
	//DEBUG: END*/
	sql_str := fmt.Sprintf("SELECT COUNT(*) from cache where boardIndex = '%s'", boardIndex)
	//fmt.Printf(" SQL : %s\n", sql_str)

	var count int
	err := db.connCache.QueryRow(sql_str).Scan(&count)
	checkErr(err)

	if count > 0 {
		return true
	}

	return false

}

func (db *Database) GetCacheEntry(boardIndex string) (bool, shared.NNOut) {

	var nnOut shared.NNOut
	var ret bool = false

	/*//DEBUG
	return false, nnOut
	//DEBUG: END*/
	sql_str := fmt.Sprintf("SELECT json from cache where boardIndex = '%s' limit 1", boardIndex)
	//fmt.Printf(" SQL : %s\n", sql_str)
	var jsonString string
	err := db.connCache.QueryRow(sql_str).Scan(&jsonString)
	checkErr(err)

	if len(jsonString) > 0 {
		ret = true
		//fmt.Println("Json:" + jsonString)
		err = json.Unmarshal([]byte(jsonString), &nnOut)
		checkErr(err)
	}

	return ret, nnOut
}
