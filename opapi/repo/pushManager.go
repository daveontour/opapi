package repo

import (
	"bytes"
	"crypto/tls"
	"os"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/models"
	"github.com/daveontour/opapi/opapi/timeservice"
)

// Channels for handling the push notifications
// The size of the channel is the number of elements the channel can
// buffer without blocking
var changePushJobChannel = make(chan models.ChangePushJob, 20)
var schedulePushJobChannel = make(chan models.SchedulePushJob, 20)

// func ReloadschedulePushes(airportCode string) {
// 	if _, ok := globals.SchedulerMap[airportCode]; ok {
// 		globals.SchedulerMap[airportCode].Clear()
// 		delete(globals.SchedulerMap, airportCode)
// 	}
// 	go SchedulePushes(airportCode, false)
// }

func StartChangePushWorkerPool(numWorkers int) {
	for w := 1; w <= numWorkers; w++ {
		go executeChangePushWorker(w, changePushJobChannel)
	}
}
func StartSchedulePushWorkerPool(numWorkers int) {
	for w := 1; w <= numWorkers; w++ {
		go executeScheduledPushWorker(w, schedulePushJobChannel)
	}
}

func SchedulePushes(airportCode string, demoMode bool) {

	StartChangePushWorkerPool(globals.ConfigViper.GetInt("NumberOfChangePushWorkers"))
	StartSchedulePushWorkerPool(globals.ConfigViper.GetInt("NumberOfSchedulePushWorkers"))

	today := time.Now().Format("2006-01-02")
	s := gocron.NewScheduler(time.Local)

	globals.SchedulerMap[airportCode] = s

	for _, u := range globals.GetUserProfiles() {

		if !u.Enabled {
			continue
		}

		for _, sub := range u.UserPushSubscriptions {
			if sub.Airport != airportCode || !sub.Enabled || (demoMode && !sub.EnableInDemoMode) {
				continue
			}

			startTimeStr := today + "T" + sub.Time
			startTime, _ := time.ParseInLocation("2006-01-02T15:04:05", startTimeStr, timeservice.Loc)
			user := u.UserName
			token := u.Key

			if sub.ReptitionHours != 0 {
				s.Every(sub.ReptitionHours).Hours().StartAt(startTime).Tag(token).Do(func() {
					schedulePushJobChannel <- models.SchedulePushJob{Sub: sub, UserToken: token, UserName: user, UserProfile: &u}
				})
				globals.Logger.Info(fmt.Sprintf("Scheduled Push for user %s, starting from %s, repeating every %v hours", u.UserName, startTimeStr, sub.ReptitionHours))
			}
			if sub.ReptitionMinutes != 0 {
				s.Every(sub.ReptitionMinutes).Minutes().StartAt(time.Now()).Tag(token).Do(func() {
					schedulePushJobChannel <- models.SchedulePushJob{Sub: sub, UserToken: token, UserName: user, UserProfile: &u}
				})
				globals.Logger.Info(fmt.Sprintf("Scheduled Push for user %s, starting from now, repeating every %v minutes", u.UserName, sub.ReptitionMinutes))

			}

			if sub.PushOnStartUp {
				schedulePushJobChannel <- models.SchedulePushJob{Sub: sub, UserToken: token, UserName: user, UserProfile: &u}
			}
		}
	}

	s.StartBlocking()
}

func HandleFlightUpdate(mess models.FlightUpdateChannelMessage) {
	checkForImpactedSubscription(mess)
	return
}

func HandleFlightCreate(mess models.FlightUpdateChannelMessage) {
	checkForImpactedSubscription(mess)
	return
}

func HandleFlightDelete(flt models.Flight) {
	checkForImpactedDeleteSubscription(flt)
	return
}

// Check if any of the registered change subscriptions are interested in this change
func checkForImpactedSubscription(mess models.FlightUpdateChannelMessage) {

	flt := GetRepo(mess.AirportCode).GetFlight(mess.FlightID)
	sto := flt.GetSTO()

	if sto.Local().After(time.Now().Local().Add(36 * time.Hour)) {
		return
	}

	globals.UserChangeSubscriptionsMutex.Lock()
	defer globals.UserChangeSubscriptionsMutex.Unlock()

NextSub:
	for _, sub := range globals.UserChangeSubscriptions {

		if !sub.Enabled {
			continue
		}
		if sub.Airport != (*flt).GetIATAAirport() {
			continue
		}

		if !sub.UpdateFlight && (*flt).Action == globals.UpdateAction {
			continue
		}
		if sub.All {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue NextSub
		}
		if sub.CreateFlight && (*flt).Action == globals.CreateAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue NextSub
		}
		if !sub.DeleteFlight && (*flt).Action == globals.DeleteAction {
			continue
		}

		if sub.DeleteFlight && (*flt).Action == globals.DeleteAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue NextSub
		}
		if sub.CreateFlight && (*flt).Action == globals.CreateAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue NextSub
		}
		if !sub.UpdateFlight && (*flt).Action == globals.UpdateAction {
			continue
		}
		// Required Parameter Field Changes
		for _, change := range (*flt).FlightChanges.Changes {

			if globals.Contains(sub.ParameterChange, change.PropertyName) {
				changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
				continue NextSub
			}

			if (change.PropertyName == "Stand" && sub.StandChange) ||
				(change.PropertyName == "Gate" && sub.GateChange) ||
				(change.PropertyName == "CheckInCounters" && sub.CheckInChange) ||
				(change.PropertyName == "Carousel" && sub.CarouselChange) ||
				(change.PropertyName == "Chute" && sub.ChuteChange) {

				changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
				continue NextSub
			}

		}

		if sub.CheckInChange && (*flt).FlightChanges.CheckinSlotsChange != nil {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue
		}
		if sub.GateChange && (*flt).FlightChanges.GateSlotsChange != nil {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue
		}
		if sub.StandChange && (*flt).FlightChanges.StandSlotsChange != nil {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue
		}
		if sub.ChuteChange && (*flt).FlightChanges.ChuteSlotsChange != nil {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue
		}
		if sub.CarouselChange && (*flt).FlightChanges.CarouselSlotsChange != nil {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue
		}

		if sub.AircraftTypeOrRegoChange && (*flt).FlightChanges.AircraftTypeChange != nil {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue
		}
		if sub.AircraftTypeOrRegoChange && (*flt).FlightChanges.AircraftChange != nil {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt}
			continue
		}
	}

	return
}
func checkForImpactedDeleteSubscription(flt models.Flight) {

	sto := flt.GetSTO()

	if sto.Local().After(time.Now().Local().Add(36 * time.Hour)) {
		return
	}

	globals.UserChangeSubscriptionsMutex.Lock()
	defer globals.UserChangeSubscriptionsMutex.Unlock()

	for _, sub := range globals.UserChangeSubscriptions {

		if !sub.Enabled {
			continue
		}
		if sub.Airport != flt.GetIATAAirport() {
			continue
		}

		if sub.DeleteFlight && flt.Action == globals.DeleteAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: &flt}
		}
	}

	return
}

func executeChangePushWorker(id int, jobs <-chan models.ChangePushJob) {

	for job := range jobs {
		globals.Logger.Debug(fmt.Sprintf("Push Worker: %d Executing Change Push for User ", id))

		queryBody, _ := json.Marshal(*job.Flight)
		bodyReader := bytes.NewReader([]byte(queryBody))

		req, err := http.NewRequest(http.MethodPost, job.Sub.DestinationURL, bodyReader)
		if err != nil {
			globals.Logger.Error(fmt.Sprintf("Change Push Client: could not create change request: %s\n", err))
		}

		req.Header.Set("Content-Type", "application/json")
		for _, pair := range job.Sub.HeaderParameters {
			req.Header.Add(pair.Parameter, pair.Value)
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: job.Sub.TrustBadCertificates},
		}
		client := http.Client{
			Timeout:   20 * time.Second,
			Transport: tr,
		}
		r, sendErr := client.Do(req)
		if sendErr != nil {
			globals.Logger.Error(fmt.Sprintf("Change Push Client. Error making http request: %s", sendErr))
			return
		}
		if r == nil {
			globals.Logger.Error(fmt.Sprintf("Scheduled Push Client for user: Error making http request to: %s\n", job.Sub.DestinationURL))
			return
		}
		if r.StatusCode != 200 {
			globals.Logger.Error(fmt.Sprintf("Change Push Client. Error making HTTP request: Returned status code = %v. URL = %s", r.StatusCode, job.Sub.DestinationURL))
			return
		}
	}
}

func executeScheduledPushWorker(id int, jobs <-chan models.SchedulePushJob) {

	for job := range jobs {

		globals.Logger.Info(fmt.Sprintf("Executing Scheduled Push for User %s", job.UserName))

		if strings.ToLower(job.Sub.SubscriptionType) == "flight" {

			flightresponse, _ := GetRequestedFlightsSub(job.Sub, job.UserToken)
			fileName, _ := writeFlightResponseToFile(flightresponse, job.UserProfile)

			defer func() {
				globals.FileDeleteChannel <- fileName
			}()

			sendViaHTTPClient(fileName, &job)

		} else if strings.ToLower(job.Sub.SubscriptionType) == "resource" {
			resourceresponse, _ := GetResourceSub(job.Sub, job.UserToken)
			fileName, _ := writeResourceResponseToFile(resourceresponse, job.UserProfile)

			defer func() {
				globals.FileDeleteChannel <- fileName
			}()

			sendViaHTTPClient(fileName, &job)
		}
	}
}

func sendViaHTTPClient(fileName string, job *models.SchedulePushJob) {

	bytesdata, _ := os.ReadFile(fileName)

	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Temp file deleted for HTTP Client Send")
		}
	}()

	req, err := http.NewRequest(http.MethodPost, job.Sub.DestinationURL, bytes.NewReader(bytesdata))
	if err != nil {
		globals.Logger.Error(fmt.Sprintf("Scheduled Push Client for user %s: could not create request: %s\n", job.UserName, err))
	}

	req.Header.Set("Content-Type", "application/json")
	for _, pair := range job.Sub.HeaderParameters {
		req.Header.Add(pair.Parameter, pair.Value)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: job.Sub.TrustBadCertificates},
	}
	client := http.Client{
		Timeout:   20 * time.Second,
		Transport: tr,
	}
	r, sendErr := client.Do(req)

	if sendErr != nil {
		globals.Logger.Error(fmt.Sprintf("Scheduled Push Client for user %s: Error making http request: %s\n", job.UserName, sendErr))
		return
	}
	if r == nil {
		globals.Logger.Error(fmt.Sprintf("Scheduled Push Client for user %s: Error making http request to: %s\n", job.UserName, job.Sub.DestinationURL))
		return
	}
	if r.StatusCode != 200 {
		globals.Logger.Error(fmt.Sprintf("Scheduled Push Client. Error making HTTP request: Returned status code = %v. URL = %s", r.StatusCode, job.Sub.DestinationURL))
		return
	}
}
