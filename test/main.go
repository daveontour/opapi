package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var wg sync.WaitGroup

func main() {

	wg.Add(1)
	go main2()
	wg.Wait()
}

func main2() {

	const file string = "activities.db"
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)

	sts := `
	DROP TABLE IF EXISTS testdate;
	CREATE TABLE testdate(id INTEGER PRIMARY KEY, sto INTEGER);
	`
	_, err = db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

	stm, err := db.Prepare("INSERT INTO testdate (id, sto) VALUES (?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stm.Close()

	for i := 1; i <= 3; i++ {
		_, err = stm.Exec(i, time.Now().Unix())
		if err != nil {
			fmt.Printf(err.Error())
		}

		time.Sleep(5 * time.Second)
	}

	wg.Done()
}
