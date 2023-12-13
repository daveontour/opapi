package repo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"time"

	"github.com/daveontour/opapi/opapi/globals"
	gobstorage "github.com/daveontour/opapi/opapi/gob"
	"github.com/daveontour/opapi/opapi/models"
)

func UpdateFlightEntry(message string, append bool, notify bool) {

	var envel models.FlightUpdatedNotificationEnvelope
	e := xml.Unmarshal([]byte(message), &envel)
	if e != nil {
		globals.Logger.Error(fmt.Errorf("Error decoding create flight message %s", e))
		return
	}

	flight := envel.Content.FlightUpdatedNotification.Flight

	airportCode := flight.GetIATAAirport()
	repo := GetRepo(airportCode)

	if repo == nil {
		globals.Logger.Warn(fmt.Sprintf("Message for unmanaged airport %s received", airportCode))
		return
	}

	sdot := flight.GetSDO()

	if sdot.Before(time.Now().AddDate(0, 0, repo.FlightSDOWindowMinimumInDaysFromNow-2)) {
		globals.Logger.Debugf("Update for Flight Before Window. Flight ID: %s", flight.GetFlightID())
		return
	}
	if sdot.After(time.Now().AddDate(0, 0, repo.FlightSDOWindowMaximumInDaysFromNow+2)) {
		globals.Logger.Debugf("Update for Flight After Window. Flight ID: %s", flight.GetFlightID())
		return
	}

	flight.LastUpdate = time.Now()
	flight.Action = globals.UpdateAction

	// Allocation changes are expensive, so only do them if necessary

	resourceChange := flight.FlightChanges.CheckinSlotsChange != nil ||
		flight.FlightChanges.GateSlotsChange != nil ||
		flight.FlightChanges.StandSlotsChange != nil ||
		flight.FlightChanges.CarouselSlotsChange != nil ||
		flight.FlightChanges.ChuteSlotsChange != nil ||
		flight.FlightChanges.AircraftChange != nil ||
		flight.FlightChanges.AircraftTypeChange != nil ||
		flight.FlightChanges.RouteChange != nil

	if globals.UseGobStorage {
		gobstorage.StoreFlight(flight)
		if resourceChange {
			upadateAllocation(flight, airportCode, false)
		}
	} else {
		globals.MapMutex.Lock()

		if append {
			repo.FlightLinkedList.AddNode(flight)
			if resourceChange {
				upadateAllocation(flight, airportCode, true)
			}

		} else {
			repo.FlightLinkedList.ReplaceOrAddNode(flight)
			if resourceChange {
				upadateAllocation(flight, airportCode, false)
			}

		}
		globals.MapMutex.Unlock()
	}

	if notify {
		globals.FlightUpdatedChannel <- models.FlightUpdateChannelMessage{FlightID: flight.GetFlightID(), AirportCode: airportCode}
	}
}
func createFlightEntry(message string, notify bool) {

	var e error

	var envel models.FlightCreatedNotificationEnvelope
	e = xml.Unmarshal([]byte(message), &envel)
	if e != nil {
		globals.Logger.Error(fmt.Errorf("Error decoding create flight message %s", e))
		return
	}

	flight := envel.Content.FlightCreatedNotification.Flight
	gobstorage.StoreFlight(flight)
	flight.LastUpdate = time.Now()
	flight.Action = globals.CreateAction

	airportCode := flight.GetIATAAirport()
	repo := GetRepo(airportCode)
	sdot := flight.GetSDO()

	if sdot.Before(time.Now().AddDate(0, 0, GetRepo(airportCode).FlightSDOWindowMinimumInDaysFromNow-2)) {
		log.Println("Create for Flight Before Window")
		return
	}
	if sdot.After(time.Now().AddDate(0, 0, GetRepo(airportCode).FlightSDOWindowMaximumInDaysFromNow+2)) {
		log.Println("Create for Flight After Window")
		return
	}

	if globals.UseGobStorage {
		gobstorage.StoreFlight(flight)
	} else {
		repo.FlightLinkedList.ReplaceOrAddNode(flight)
	}
	upadateAllocation(flight, airportCode, false)

	if notify {
		globals.FlightCreatedChannel <- models.FlightUpdateChannelMessage{FlightID: flight.GetFlightID(), AirportCode: airportCode}
	}
}
func deleteFlightEntry(message string, notify bool) {

	var envel models.FlightDeletedNotificationEnvelope
	e := xml.Unmarshal([]byte(message), &envel)
	if e != nil {
		globals.Logger.Error(fmt.Errorf("Error decoding delte flight message %s", e))
		return
	}

	flight := envel.Content.FlightDeletedNotification.Flight
	flight.Action = globals.DeleteAction

	airportCode := flight.GetIATAAirport()
	repo := GetRepo(airportCode)

	flightCopy := flight
	if globals.UseGobStorage {
		gobstorage.DeleteFlight(flight)
	} else {
		repo.FlightLinkedList.RemoveNode(flight)
		repo.RemoveFlightAllocation(flight.GetFlightID())
	}

	if notify {
		globals.FlightDeletedChannel <- flightCopy
	}
}
func unhandledNotificationMessage(message string, notify bool) {

	if notify {
		globals.UnhandledNotificationChannel <- message
	}
}
func getFlights(airportCode string, values ...int) []byte {

	repo := GetRepo(airportCode)
	from := time.Now().AddDate(0, 0, repo.FlightSDOWindowMinimumInDaysFromNow).Format("2006-01-02")
	to := time.Now().AddDate(0, 0, repo.FlightSDOWindowMaximumInDaysFromNow+1).Format("2006-01-02")

	// Change the window based on optional inout parameters
	if len(values) >= 1 {
		from = time.Now().AddDate(0, 0, values[0]).Format("2006-01-02")
	}

	// Add in a sneaky extra day
	if len(values) >= 2 {
		to = time.Now().AddDate(0, 0, values[1]+1).Format("2006-01-02")
	}

	globals.Logger.Debug(fmt.Sprintf("Getting flight from %s to %s", from, to))
	fmt.Printf("Getting flights from %s to %s\n", from, to)

	queryBody := fmt.Sprintf(xmlGetFlightsTemplateBody, repo.AMSToken, from, to, repo.AMSAirport)
	bodyReader := bytes.NewReader([]byte(queryBody))

	req, err := http.NewRequest(http.MethodPost, repo.AMSSOAPServiceURL, bodyReader)
	if err != nil {
		globals.Logger.Error(fmt.Sprintf("client: could not create request: %s\n", err))
	}

	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Add("SOAPAction", "http://www.sita.aero/ams6-xml-api-webservice/IAMSIntegrationService/GetFlights")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		globals.Logger.Error(fmt.Sprintf("client: error making http request: %s\n", err))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		globals.Logger.Error(fmt.Sprintf("client: could not read response body: %s\n", err))
	}

	fmt.Printf("Got flights from %s to %s\n", from, to)
	return resBody
}
func upadateAllocation(flight models.Flight, airportCode string, bypassDelete bool) {

	//defer exeTime(fmt.Sprintf("Updated allocations for Flight %s", flight.GetFlightID()))()
	// Testing with 3000 flights showed unmeasurable time to 500 micro seconds, so no worries mate

	if globals.UseGobStorage {
		gobstorage.DeleteFlightResourceAllocation(flight, airportCode)
	}
	repo := GetRepo(airportCode)

	// It's too messy to do CRUD operations, so just delete all the allocations and then create them again from the current message
	//bypass delete is used for init population for perfTest or demo mode
	if !bypassDelete {
		repo.RemoveFlightAllocation(flight.GetFlightID())
	}
	flightId := flight.GetFlightID()
	direction := flight.GetFlightDirection()
	route := flight.GetFlightRoute()
	aircraftType := flight.GetAircraftType()
	aircraftRegistration := flight.GetAircraftRegistration()
	now := time.Now().Local()

	for _, checkInSlot := range flight.FlightState.CheckInSlots.CheckInSlot {
		checkInID, start, end := checkInSlot.GetResourceID()
		resourceArea := "Area"
		_, resourceArea, err := repo.GetResourceDetail("CHECKIN", checkInID)
		if err != nil {
			resourceArea = "Not Found"
		}

		if checkInID == "" {
			checkInID = "-"
		}

		allocation := models.AllocationItem{
			ResourceID:           checkInID,
			ResourceArea:         resourceArea,
			From:                 start,
			To:                   end,
			FlightID:             flightId,
			AirportCode:          airportCode,
			Direction:            direction,
			Route:                route,
			AircraftType:         aircraftType,
			AircraftRegistration: aircraftRegistration,
			LastUpdate:           now}

		if globals.UseGobStorage {
			gobstorage.StoreResourceAllocation(allocation, "CHECKIN", airportCode)
		} else {
			(*repo).CheckInList.AddAllocation(allocation)
		}
	}

	for _, gateSlot := range flight.FlightState.GateSlots.GateSlot {
		gateID, start, end := gateSlot.GetResourceID()
		resourceArea := "Area"
		_, resourceArea, err := repo.GetResourceDetail("GATE", gateID)
		if err != nil {
			resourceArea = "Not Found"
		}

		if gateID == "" {
			gateID = "-"
		}

		allocation := models.AllocationItem{
			ResourceID:           gateID,
			ResourceArea:         resourceArea,
			From:                 start,
			To:                   end,
			FlightID:             flightId,
			AirportCode:          airportCode,
			Direction:            direction,
			Route:                route,
			AircraftType:         aircraftType,
			AircraftRegistration: aircraftRegistration,
			LastUpdate:           now}

		if globals.UseGobStorage {
			gobstorage.StoreResourceAllocation(allocation, "GATE", airportCode)
		} else {
			(*repo).GateList.AddAllocation(allocation)
		}
	}

	for _, standSlot := range flight.FlightState.StandSlots.StandSlot {
		standID, start, end := standSlot.GetResourceID()
		resourceArea := "Area"
		_, resourceArea, err := repo.GetResourceDetail("STAND", standID)
		if err != nil {
			resourceArea = "Not Found"
		}

		if standID == "" {
			standID = "-"
		}

		allocation := models.AllocationItem{
			ResourceID:           standID,
			ResourceArea:         resourceArea,
			From:                 start,
			To:                   end,
			FlightID:             flightId,
			AirportCode:          airportCode,
			Direction:            direction,
			Route:                route,
			AircraftType:         aircraftType,
			AircraftRegistration: aircraftRegistration,
			LastUpdate:           now}

		if globals.UseGobStorage {
			gobstorage.StoreResourceAllocation(allocation, "STAND", airportCode)
		} else {
			(*repo).StandList.AddAllocation(allocation)
		}
	}

	for _, carouselSlot := range flight.FlightState.CarouselSlots.CarouselSlot {
		carouselID, start, end := carouselSlot.GetResourceID()
		resourceArea := "Area"
		_, resourceArea, err := repo.GetResourceDetail("CAROUSEL", carouselID)
		if err != nil {
			resourceArea = "Not Found"
		}

		if carouselID == "" {
			carouselID = "-"
		}

		allocation := models.AllocationItem{
			ResourceID:           carouselID,
			ResourceArea:         resourceArea,
			From:                 start,
			To:                   end,
			FlightID:             flightId,
			AirportCode:          airportCode,
			Direction:            direction,
			Route:                route,
			AircraftType:         aircraftType,
			AircraftRegistration: aircraftRegistration,
			LastUpdate:           now}

		if globals.UseGobStorage {
			gobstorage.StoreResourceAllocation(allocation, "CAROUSEL", airportCode)
		} else {
			(*repo).CarouselList.AddAllocation(allocation)
		}

	}

	for _, chuteSlot := range flight.FlightState.ChuteSlots.ChuteSlot {
		chuteID, start, end := chuteSlot.GetResourceID()
		resourceArea := "Area"
		_, resourceArea, err := repo.GetResourceDetail("CHUTE", chuteID)
		if err != nil {
			resourceArea = "Not Found"
		}

		if chuteID == "" {
			chuteID = "-"
		}
		allocation := models.AllocationItem{
			ResourceID:           chuteID,
			ResourceArea:         resourceArea,
			From:                 start,
			To:                   end,
			FlightID:             flightId,
			AirportCode:          airportCode,
			Direction:            direction,
			Route:                route,
			AircraftType:         aircraftType,
			AircraftRegistration: aircraftRegistration,
			LastUpdate:           now}

		if globals.UseGobStorage {
			gobstorage.StoreResourceAllocation(allocation, "CHUTE", airportCode)
		} else {
			(*repo).ChuteList.AddAllocation(allocation)
		}
	}
}
