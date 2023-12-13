package gobstorage

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/daveontour/opapi/opapi/models"
	"github.com/daveontour/opapi/opapi/timeservice"
	_ "github.com/mattn/go-sqlite3"
)

var enc *gob.Encoder
var dec *gob.Decoder
var db *sql.DB
var mu sync.Mutex

const cacheFile string = "flightcache.db"

func GobStorageInit() {

	var err0 error

	db, err0 = sql.Open("sqlite3", cacheFile)

	if err0 != nil {
		log.Fatal(err0)
	}

	GobStorageFlightInit()
	GobStorageResourceInit()
}
func GobStorageFlightInit() {

	gob.Register(models.Flight{})

	var version string
	err := db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	sts := `
	DROP TABLE IF EXISTS flightgob;
	CREATE TABLE flightgob(id STRING PRIMARY KEY, airport STRING,  airline STRING, fltnum STRING, kind STRING, route STRING, sto INTEGER, lastUpdate INTEGER, gob BLOB);
	`
	_, err = db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

}
func GobStorageResourceInit() {

	gob.Register(models.Flight{})

	sts := `
	DROP TABLE IF EXISTS resourcegob;
	CREATE TABLE resourcegob( airport STRING, allocfrom INTEGER, allocTo INTEGER, flightID STRING, direction STRING, route STRING, acType STRING, acRego STRING, lastUpdate INTEGER, resourceType STRING, name STRING, area STRING);
	`
	_, err := db.Exec(sts)

	if err != nil {
		log.Fatal(err)
	}

}

func StoreResourceAllocation(aI models.AllocationItem, resourceType, apt string) {

	mu.Lock()
	defer mu.Unlock()

	stm, err := db.Prepare("INSERT INTO resourcegob (airport, allocfrom, allocTo, flightID, direction, route, acType, acRego, lastUpdate, resourceType, name, area) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stm.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stm.Exec(apt, aI.From.Unix(), aI.To.Unix(), aI.FlightID, aI.Direction, aI.Route, aI.AircraftType, aI.AircraftRegistration, time.Now().Unix(), resourceType, aI.ResourceID, aI.ResourceArea)

	if err != nil {
		log.Fatal(err)
	}
}

func DeleteFlightResourceAllocation(flt models.Flight, airportCode string) {

	mu.Lock()
	defer mu.Unlock()

	id, _, _, _, _, _, _, _ := flt.GetGobParameters()

	stm1, err := db.Prepare("DELETE FROM resourcegob WHERE flightId = ? AND airport = ?")
	defer stm1.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stm1.Exec(id, airportCode)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func StoreFlight(flt models.Flight) {

	mu.Lock()
	defer mu.Unlock()

	id, airport, airline, fltnum, kind, route, _, stoUnix := flt.GetGobParameters()

	stm0, err := db.Prepare("DELETE FROM flightgob WHERE id = ?")
	defer stm0.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stm0.Exec(id)
	if err != nil {
		fmt.Printf(err.Error())
	}
	stm0.Close()

	stm, err := db.Prepare("INSERT INTO flightgob (id, airport, airline, fltnum, kind, route, sto, lastUpdate, gob) VALUES (?,?,?,?,?,?,?,?,?)")
	defer stm.Close()
	if err != nil {
		log.Fatal(err)
	}

	var bBuf bytes.Buffer // Standard input
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(models.Flight{})
	// Create an encoder interface and send values

	enc = gob.NewEncoder(&bBuf)
	err = enc.Encode(&flt)
	_, err = stm.Exec(id, airport, airline, fltnum, kind, route, stoUnix, time.Now().Unix(), bBuf.Bytes())
	fmt.Println("Record written to data store")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func DeleteFlight(flt models.Flight) {
	mu.Lock()
	defer mu.Unlock()

	id, _, _, _, _, _, _, _ := flt.GetGobParameters()

	stm1, err0 := db.Prepare("DELETE FROM resourcegob WHERE flightId = ?")
	defer stm1.Close()

	_, err := stm1.Exec(id)
	if err != nil {
		fmt.Printf(err.Error())
	}
	stm1.Close()

	stm0, err0 := db.Prepare("DELETE FROM flightgob WHERE id = ?")
	defer stm0.Close()
	if err0 != nil {
		log.Fatal(err)
	}
	_, err = stm0.Exec(id)
	if err != nil {
		fmt.Printf(err.Error())
	}

}

func GetFlight(flightID, apt string) (fltptr *models.Flight) {

	mu.Lock()
	defer mu.Unlock()

	sql := "SELECT id, gob FROM flightgob WHERE airport = ? and id = ?"
	stm, err := db.Prepare(sql)
	defer stm.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stm.Query(apt, flightID)
	stm.Close()

	if err != nil {
		log.Fatal(err)
	}

	var data []byte
	var flt models.Flight

	for rows.Next() {
		err := rows.Scan(&flightID, &data)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		dec := gob.NewDecoder(reader)

		err = dec.Decode(&flt)
		if err != nil {
			log.Fatal("decode:", err)
		}

		fltptr = &flt
	}

	return
}

func GetFlights(request models.Request, allowedAllAirline bool, to, from time.Time, airport string) (FlightLinkedList models.FlightLinkedList) {

	mu.Lock()
	defer mu.Unlock()

	al := request.Airline

	if allowedAllAirline && request.Airline == "" {
		al = "%"
	}

	kind := "%"
	if strings.HasPrefix(request.Direction, "D") {
		kind = "D%"
	}
	if strings.HasPrefix(request.Direction, "A") {
		kind = "A%"
	}

	fltNum := "%"
	if request.FltNum != "" {
		fltNum = "%" + request.FltNum + "%"
	}

	route := "%"
	if request.Route != "" {
		route = "%" + route + "%"
	}
	updatedSince := 0
	if request.UpdatedSince != "" {
		t, err := time.ParseInLocation("2006-01-02T15:04:05", request.UpdatedSince, timeservice.Loc)
		if err != nil {
			log.Fatal(err)
		} else {
			updatedSince = int(t.Unix())
		}
	}

	sql := "SELECT id, gob FROM flightgob WHERE airport = ? AND airline LIKE ? AND id LIKE ? AND kind LIKE ? AND route LIKE ? AND sto >= ? AND sto <= ? AND lastUpdate > ? ORDER BY sto ASC"
	stm, err := db.Prepare(sql)
	defer stm.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stm.Query(airport, al, fltNum, kind, route, strconv.FormatInt(from.Unix(), 10), strconv.FormatInt(to.Unix(), 10), updatedSince)
	stm.Close()

	if err != nil {
		log.Fatal(err)
	}

	var flightID string
	var data []byte

	for rows.Next() {
		err := rows.Scan(&flightID, &data)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		dec := gob.NewDecoder(reader)
		var flt models.Flight
		err = dec.Decode(&flt)
		if err != nil {
			log.Fatal("decode:", err)
			log.Println(flightID)
		} else {
			FlightLinkedList.AddNode(flt)
		}
	}

	return
}

func CleanRepository(apt string, from time.Time) (removed int64, postCount int, maxSto int, minSto int) {

	mu.Lock()
	defer mu.Unlock()

	sql := "DELETE FROM flightgob WHERE airport = ?  AND sto <= ? "
	stm, err := db.Prepare(sql)
	defer stm.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	res, err := stm.Exec(apt, strconv.FormatInt(from.Unix(), 10))
	removed, _ = res.RowsAffected()
	stm.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err := db.Query("SELECT COUNT(*) FROM flightgob")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&postCount); err != nil {
			log.Fatal(err)
		}
	}

	rows, err = db.Query("SELECT MAX(sto)as maxSto, MIN(sto)as minSto FROM flightgob")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&maxSto, &minSto); err != nil {
			maxSto = -1
			minSto = -1
		}
	}
	return
}

func GetResourceAllocation(apt string, resource string, airline string, flightID string, direction string, route string, fromTime time.Time, toTime time.Time, updatedSince string, resourceType string) (alloc []models.AllocationResponseItem) {

	//	sql := fmt.Sprintf("SELECT * FROM resourcegob WHERE airport = '%s' AND resourceType = '%s' AND flightID LIKE ? AND ( (allocFrom > ? AND allocFrom < ?) OR (allocTo > ? AND allocTo < ?))", apt, strings.ToUpper(resourceType))
	sql := fmt.Sprintf("SELECT flightID, direction, route, acType, acRego, resourceType, name, area, allocfrom, allocTo, lastUpdate FROM resourcegob WHERE airport = '%s'", apt)

	if airline != "" {
		sql = sql + " AND flightID LIKE '" + airline + "%'"
	}

	if flightID != "" {
		sql = sql + " AND flightID LIKE '%" + flightID + "%'"
	}
	if resource != "" {
		sql = sql + " AND name LIKE '%" + resource + "%'"
	}
	if resourceType != "" {
		sql = sql + " AND resourceType = '" + strings.ToUpper(resourceType) + "'"
	}

	sql = sql + fmt.Sprintf(" AND ( (allocfrom > %v AND allocfrom < %v) OR (allocTo > %v AND allocTo < %v))", fromTime.Unix(), toTime.Unix(), fromTime.Unix(), toTime.Unix())

	if updatedSince != "" {
		updatedSinceTime, updatedSinceErr := time.ParseInLocation("2006-01-02T15:04:05", updatedSince, timeservice.Loc)
		if updatedSinceErr == nil {
			sql = sql + fmt.Sprintf(" AND lastUpdate >= %v", updatedSinceTime.Unix())
		}

	}

	stm, err := db.Prepare(sql)
	defer stm.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	rows, err := stm.Query()
	if err != nil {
		log.Fatal(err)
	}
	stm.Close()

	var acType string
	var acRego string
	var name string
	var area string
	var allocFrom int
	var allocTo int
	var lastUpdate int

	for rows.Next() {
		err := rows.Scan(&flightID, &direction, &route, &acType, &acRego, &resourceType, &name, &area, &allocFrom, &allocTo, &lastUpdate)
		if err != nil {
			log.Fatal(err)
		}

		tmf := time.Unix(int64(allocFrom), 0)
		tmt := time.Unix(int64(allocTo), 0)
		tml := time.Unix(int64(lastUpdate), 0)

		n := models.AllocationResponseItem{
			AllocationItem: models.AllocationItem{
				From:                 tmf,
				To:                   tmt,
				FlightID:             flightID,
				Direction:            direction,
				Route:                route,
				AircraftType:         acType,
				AircraftRegistration: acRego,
				LastUpdate:           tml},
			ResourceType: resourceType,
			Name:         name,
			Area:         area,
		}
		alloc = append(alloc, n)

	}

	return
}
