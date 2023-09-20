package repo

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/models"
	"github.com/daveontour/opapi/opapi/timeservice"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context, userToken string) *models.UserProfile {

	//defer globals.ExeTime("Getting User Profile")()

	//key := userToken

	if c != nil {
		keys := c.Request.Header["Token"]
		userToken = "default"

		if keys != nil {
			userToken = keys[0]
		}

	}

	for _, u := range globals.GetUserProfiles() {
		if userToken == u.Key {
			return &u
		}
	}

	return nil
}

func GetRequestedFlightsAPI(c *gin.Context) {
	//defer globals.ExeTime(fmt.Sprintf("Get Flight Processing time for %s", c.Request.RequestURI))()
	defer func() {
		// var m runtime.MemStats
		// runtime.ReadMemStats(&m)
		// fmt.Println("Initial HeapAlloc: ", m.HeapAlloc)

		// Trigger the garbage collector
		runtime.GC()

		// // Print the memory usage after the garbage collector has run
		// runtime.ReadMemStats(&m)
		// fmt.Println("After GC HeapAlloc: ", m.HeapAlloc)
	}()

	userProfile := GetUserProfile(c, "")

	if !userProfile.Enabled {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "User Access Has Been Disabled"})
		return
	}

	globals.RequestLogger.Info("User: " + userProfile.UserName + " IP: " + c.RemoteIP() + " Request:% " + c.Request.RequestURI)

	//apt := c.Param("apt")
	direction := strings.ToUpper(c.Query("direction"))
	if direction == "" {
		direction = strings.ToUpper(c.Query("d"))
	}
	airline := c.Query("al")
	flt := c.Query("flt")
	if flt == "" {
		flt = c.Query("flight")
	}

	route := strings.ToUpper(c.Query("route"))
	if route == "" {
		route = c.Query("r")
	}

	response, err := GetRequestedFlightsCommon(c.Param("apt"), direction, airline, flt, c.Query("from"), c.Query("to"), route, "", c, nil)
	if err.Err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	totalFlights := -1
	if c.Query("tf") == "true" {
		totalFlights = GetRepo(c.Param("apt")).FlightLinkedList.Len()
	}
	response.TotalFlights = totalFlights
	fileName, err2 := writeFlightResponseToFile(response, userProfile, c.Query("max"), true)
	defer func() {
		if os.Remove(fileName) != nil {
			fmt.Println(err.Error())
		}
	}()

	if err2 == nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.File(fileName)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}
}

func GetRequestedFlightsSub(sub models.UserPushSubscription, userToken string) (models.Response, models.GetFlightsError) {
	apt := sub.Airport
	direction := strings.ToUpper(sub.Direction)
	airline := sub.Airline
	from := sub.From
	to := sub.To
	route := strings.ToUpper(sub.Route)
	qf := sub.QueryableCustomFields

	return GetRequestedFlightsCommon(apt, direction, airline, "", strconv.Itoa(from), strconv.Itoa(to), route, userToken, nil, qf)

}
func GetRequestedFlightsCommon(apt, direction, airline, flt, from, to, route, userToken string, c *gin.Context, qf []models.ParameterValuePair) (models.Response, models.GetFlightsError) {

	userProfilePtr := GetUserProfile(c, userToken)

	// Create the response object so we can return early if required
	response := models.Response{
		Route:       route,
		User:        userProfilePtr.UserName,
		AirportCode: apt,
	}

	// Add the flights the response object and return nil for errors
	if direction != "" {
		if strings.HasPrefix(direction, "A") {
			response.Direction = "Arrival"
		}
		if strings.HasPrefix(direction, "D") {
			response.Direction = "Departure"
		}
	} else {
		response.Direction = "ARR/DEP"
	}

	if apt == "" {
		return response, models.GetFlightsError{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("Airport not specified"),
		}
	}

	// Check that the requested airport exists in the repository

	aptPtr := GetRepo(apt)

	if aptPtr == nil {
		return response, models.GetFlightsError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("Airport %s not found", apt),
		}
	}

	// Set Default airline if none set
	if airline == "" && userProfilePtr.DefaultAirline != "" {
		airline = userProfilePtr.DefaultAirline
		response.AddWarning("Airline set to " + airline + " by the administration configuration")
	} else {
		response.Airline = airline
	}

	if flt != "" {
		response.Flight = flt
	}

	//Check that the user is allowed to access the requested airport
	if !globals.Contains(userProfilePtr.AllowedAirports, apt) &&
		!globals.Contains(userProfilePtr.AllowedAirports, "*") {
		return response, models.GetFlightsError{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("User is not allowed to access requested airport"),
		}
	}

	// Build the request object
	request := models.Request{Direction: direction, Airline: airline, FltNum: flt, From: from, To: to, UserProfile: *userProfilePtr, Route: route}

	// Reform the request based on the user Profile and the request parameters
	request, response = processCustomFieldQueries(request, response, c, qf)

	// If the user is requesting a particular airline, check that they are allowed to access that airline
	if airline != "" && userProfilePtr.AllowedAirlines != nil {
		if !globals.Contains(userProfilePtr.AllowedAirlines, airline) &&
			!globals.Contains(userProfilePtr.AllowedAirlines, "*") {
			return response, models.GetFlightsError{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New("unavailable"),
			}
		}
	}

	var err error

	// Get the filtered and pruned flights for the request
	globals.MapMutex.Lock()
	flights := aptPtr.FlightLinkedList
	response, err = filterFlights(request, response, flights, c, aptPtr)
	globals.MapMutex.Unlock()

	if err == nil {
		return response, models.GetFlightsError{
			StatusCode: http.StatusOK,
			Err:        nil,
		}
	} else {
		return response, models.GetFlightsError{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}
}

func processCustomFieldQueries(request models.Request, response models.Response, c *gin.Context, qf []models.ParameterValuePair) (models.Request, models.Response) {

	customFieldQureyMap := make(map[string]string)

	if c != nil {
		// Find the potential customField queries in the request
		queryMap := c.Request.URL.Query()
		for k, v := range queryMap {
			if !globals.Contains([]string{"airport", "airline", "al", "from", "to", "direction", "d", "route", "r", "sort", "flt", "flight"}, k) {
				customFieldQureyMap[k] = v[0]
			}
		}
	} else {
		for _, pvPair := range qf {
			parameter := pvPair.Parameter
			value := pvPair.Value
			customFieldQureyMap[parameter] = value
		}
	}
	// (even if there are rubbish values still in the request, the GetPropoerty function will handle it

	// Put in new default values
	if request.UserProfile.DefaultQueryableCustomFields != nil {
		for _, pair := range request.UserProfile.DefaultQueryableCustomFields {
			if v, ok := customFieldQureyMap[pair.Parameter]; ok {
				if v != pair.Value {
					response.AddWarning(fmt.Sprintf("Setting query against %s to default value %s", pair.Parameter, pair.Value))
				}
			}
			customFieldQureyMap[pair.Parameter] = pair.Value
		}
	}

	// Remove any queries against unauthorised fields
	remove := []string{}
	for k := range customFieldQureyMap {
		if !globals.Contains(request.UserProfile.AllowedCustomFields, k) && !globals.Contains(request.UserProfile.AllowedCustomFields, "*") {
			remove = append(remove, k)
		}
	}

	for _, k := range remove {
		delete(customFieldQureyMap, k)
		response.AddWarning(fmt.Sprintf("Ignoring unauthorised query against custom field: %s", k))
	}

	presentQueryableParameters := []models.ParameterValuePair{}

	for k, v := range customFieldQureyMap {
		presentQueryableParameters = append(presentQueryableParameters, models.ParameterValuePair{Parameter: k, Value: v})
	}

	request.PresentQueryableParameters = presentQueryableParameters

	return request, response
}
func filterFlights(request models.Request, response models.Response, flightsLinkedList models.FlightLinkedList, c *gin.Context, repo *models.Repository) (models.Response, error) {

	//defer exeTime("Filter, Prune and Sort Flights")()
	//returnFlights := []models.Flight{}

	// var from time.Time
	// var to time.Time
	var updatedSinceTime time.Time

	fromOffset, fromErr := strconv.Atoi(request.From)
	if fromErr != nil {
		fromOffset = -12
	}

	from := time.Now().Add(time.Hour * time.Duration(fromOffset))

	if from.Before(repo.CurrentLowerLimit) {
		from = repo.CurrentLowerLimit
		response.AddWarning(fmt.Sprintf("Requested lower time limit outside cache boundaries. Set to %s", from.Format("2006-01-02T15:04:05")))
	}
	response.From = from.Format("2006-01-02T15:04:05")

	toOffset, toErr := strconv.Atoi(request.To)
	if toErr != nil {
		toOffset = 24
	}

	to := time.Now().Add(time.Hour * time.Duration(toOffset))
	if to.After(repo.CurrentUpperLimit) {
		to = repo.CurrentUpperLimit
		response.AddWarning(fmt.Sprintf("Requested upper time limit outside cache boundaries. Set to %s", to.Format("2006-01-02T15:04:05")))
	}

	response.To = to.Format("2006-01-02T15:04:05")

	if request.UpdatedSince != "" {
		t, err := time.ParseInLocation("2006-01-02T15:04:05", request.UpdatedSince, timeservice.Loc)
		if err != nil {
			return models.Response{}, err
		} else {
			updatedSinceTime = t
			response.To = t.String()
		}
	}

	allowedAllAirline := false
	if request.UserProfile.AllowedAirlines != nil {
		if globals.Contains(request.UserProfile.AllowedAirlines, "*") {
			allowedAllAirline = true
		}
	}

	//filterStart := time.Now()

	currentFlight := flightsLinkedList.Head

NextFlight:
	for currentFlight != nil {

		if currentFlight.GetSTO().Before(from) {
			currentFlight = currentFlight.NextNode
			continue
		}

		if currentFlight.GetSTO().After(to) {
			currentFlight = currentFlight.NextNode
			continue
		}

		for _, queryableParameter := range request.PresentQueryableParameters {
			queryValue := queryableParameter.Value
			flightValue := currentFlight.GetProperty(queryableParameter.Parameter)

			if flightValue == "" || queryValue != flightValue {
				currentFlight = currentFlight.NextNode
				continue NextFlight
			}
		}

		// Flight direction filter
		if strings.HasPrefix(request.Direction, "D") && currentFlight.IsArrival() {
			currentFlight = currentFlight.NextNode
			continue
		}
		if strings.HasPrefix(request.Direction, "A") && !currentFlight.IsArrival() {
			currentFlight = currentFlight.NextNode
			continue
		}

		// Requested Airline Code filter
		if request.Airline != "" && currentFlight.GetIATAAirline() != request.Airline {
			currentFlight = currentFlight.NextNode
			continue
		}

		// RequestedRoute filter
		if request.Route != "" && !strings.Contains(currentFlight.GetFlightRoute(), request.Route) {
			currentFlight = currentFlight.NextNode
			continue
		}

		if request.FltNum != "" && !strings.Contains(currentFlight.GetFlightID(), request.FltNum) {
			currentFlight = currentFlight.NextNode
			continue
		}

		if request.UpdatedSince != "" {
			if currentFlight.LastUpdate.Before(updatedSinceTime) {
				currentFlight = currentFlight.NextNode
				continue
			}
		}

		// Filter out airlines that the user is not allowed to see
		// "*" entry in AllowedAirlines allows all.
		if !allowedAllAirline {
			if request.UserProfile.AllowedAirlines != nil {
				if !globals.Contains(request.UserProfile.AllowedAirlines, currentFlight.GetIATAAirline()) {
					currentFlight = currentFlight.NextNode
					continue
				}
			}
		}

		// Made it here without being filtered out, so add it to the flights to be returned.
		currentFlight.Action = globals.StatusAction
		//	response.ResponseFlights.AddNode(models.FlightResponseItem{FlightPtr: currentFlight})
		response.ResponseFlights = append(response.ResponseFlights, models.FlightResponseItem{FlightPtr: currentFlight, STO: currentFlight.GetSTO()})

		currentFlight = currentFlight.NextNode
	}

	//	globals.MetricsLogger.Info(fmt.Sprintf("Filter Flights execution time: %s", time.Since(filterStart)))

	//***Important
	//Pruning is now done at the output to avoid creating additional copies of the data structure

	response.NumberOfFlights = len(response.ResponseFlights)

	//defer globals.ExeTime(fmt.Sprintf("Sorting %v Filtered Flights", response.NumberOfFlights))()
	sort.Slice(response.ResponseFlights, func(i, j int) bool {
		return response.ResponseFlights[i].STO.Before(response.ResponseFlights[j].STO)
	})

	response.CustomFieldQuery = request.PresentQueryableParameters

	return response, nil
}

func writeFlightResponseToFile(response models.Response, userProfile *models.UserProfile, max string, statusOnly bool) (fileName string, e error) {

	file, e := os.CreateTemp("", "getflighttemp-*.json")
	if e != nil {
		fmt.Println(e)
		return
	}

	maxFlights, err := strconv.Atoi(max)
	if err != nil {
		maxFlights = -1
	}

	if maxFlights != -1 {
		response.AddWarning(fmt.Sprintf("Maximum number of returned flights limited to: %d", maxFlights))
	}
	fwb := bufio.NewWriterSize(file, 32768)
	defer file.Close()

	e = fwb.WriteByte('{')
	if e != nil {
		return
	}
	_, e = fwb.WriteString("\"Airport\":\"" + response.AirportCode + "\"," +
		"\"Direction\":\"" + response.Direction + "\"," +
		"\"ScheduleFlightsFrom\":\"" + response.From + "\"," +
		"\"ScheduleFlightsTo\":\"" + response.To + "\"," +
		"\"NumberOfFlights\":\"" + strconv.Itoa(response.NumberOfFlights) + "\",")
	if e != nil {
		return
	}

	if response.TotalFlights > 0 {
		_, e = fwb.WriteString("\"TotalFlights\":\"" + strconv.Itoa(response.TotalFlights) + "\",")
	}

	if response.Airline != "" {
		_, e = fwb.WriteString("\"Airline\":\"" + response.Airline + "\",")
	} else {
		_, e = fwb.WriteString("\"Airline\":\"*\",")
	}
	if e != nil {
		return
	}
	if response.Flight != "" {
		_, e = fwb.WriteString("\"Flight\":\"" + response.Flight + "\",")
	} else {
		_, e = fwb.WriteString("\"Flight\":\"*\",")
	}
	if e != nil {
		return
	}
	if response.Route != "" {
		_, e = fwb.WriteString("\"Route\":\"" + response.Route + "\",")
	} else {
		_, e = fwb.WriteString("\"Route\":\"*\",")
	}
	if e != nil {
		return
	}
	_, e = fwb.WriteString("\"CustomFieldQuery\":[")
	if e != nil {
		return
	}
	for idx, w := range response.CustomFieldQuery {
		if idx > 0 {
			_, e = fwb.WriteString(",")
			if e != nil {
				return
			}
		}
		_, e = fwb.WriteString("{\"Parameter\":\"" + w.Parameter + "\",\"Value\":\"" + w.Value + "\"}")
		if e != nil {
			return
		}
	}
	_, e = fwb.WriteString("],")
	if e != nil {
		return
	}

	_, e = fwb.WriteString("\"Warnings\":[")
	if e != nil {
		return
	}
	for idx, w := range response.Warnings {
		if idx > 0 {
			_, e = fwb.WriteString(",")
			if e != nil {
				return
			}
		}
		_, e = fwb.WriteString("\"" + w + "\"")
		if e != nil {
			return
		}
	}
	_, e = fwb.WriteString("],")
	if e != nil {
		return
	}

	_, e = fwb.WriteString("\"Errors\":[")
	if e != nil {
		return
	}
	for idx, w := range response.Errors {
		if idx > 0 {
			_, e = fwb.WriteString(",")
			if e != nil {
				return
			}
		}
		_, e = fwb.WriteString("\"" + w + "\"")
		if e != nil {
			return
		}
	}
	_, e = fwb.WriteString("],")
	if e != nil {
		return
	}

	e = models.WriteFlightsInJSON(fwb, response.ResponseFlights, userProfile, statusOnly, maxFlights)

	response.ResponseFlights = nil

	if e != nil {
		return
	}
	e = fwb.WriteByte('}')
	if e != nil {
		return
	}
	e = fwb.Flush()
	if e != nil {
		return
	}

	return file.Name(), nil
}
