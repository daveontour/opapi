package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/daveontour/opapi/opapi/models"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

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
	DROP TABLE IF EXISTS flightgob;
	CREATE TABLE flightgob(id INTEGER PRIMARY KEY, gob BLOB);
	`
	_, err = db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

	stm, err := db.Prepare("INSERT INTO flightgob (id, gob) VALUES (?, ?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stm.Close()

	var network bytes.Buffer // Standard input
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(models.Flight{})
	// Create an encoder interface and send values
	enc := gob.NewEncoder(&network)
	for i := 1; i <= 3; i++ {
		fl := models.Flight{Action: "TEST"}
		interfaceEncode(enc, fl)
		_, err = stm.Exec(i, network.Bytes())
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
	// Create a decoder interface and receive values
	dec := gob.NewDecoder(&network)
	for i := 1; i <= 3; i++ {
		result := interfaceDecode(dec)
		fmt.Println(result.Action)
	}
}

func interfaceEncode(enc *gob.Encoder, p models.Flight) {
	// The encode will fail unless the concrete type has been
	// registered. We registered it in the calling function.
	// Pass pointer to interface so Encode sees (and hence sends) a value of
	// interface type.  If we passed p directly it would see the concrete type instead.
	// See the blog post, "The Laws of Reflection" for background.
	err := enc.Encode(&p)
	if err != nil {
		log.Fatal("encode:", err)
	}
}

// interfaceDecode decodes the value of the interface and returns
func interfaceDecode(dec *gob.Decoder) models.Flight {
	// The decode will fail unless the concrete type on the wire has been
	// registered. We registered it in the calling function.
	var p models.Flight
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal("decode:", err)
	}
	return p
}
