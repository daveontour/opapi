package repo

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"os"

	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/models"
	"github.com/daveontour/opapi/opapi/timeservice"
	amqp "github.com/rabbitmq/amqp091-go"
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
				_, err := s.Every(sub.ReptitionHours).Hours().StartAt(startTime).Tag(token).Do(func() {
					schedulePushJobChannel <- models.SchedulePushJob{Sub: sub, UserToken: token, UserName: user, UserProfile: &u}
				})
				if err != nil {
					globals.Logger.Error(fmt.Errorf("Error scheduling hour push %s", err))
				}
				globals.Logger.Info(fmt.Sprintf("Scheduled Push for user %s, starting from %s, repeating every %v hours", u.UserName, startTimeStr, sub.ReptitionHours))
			}
			if sub.ReptitionMinutes != 0 {
				_, err := s.Every(sub.ReptitionMinutes).Minutes().StartAt(time.Now()).Tag(token).Do(func() {
					schedulePushJobChannel <- models.SchedulePushJob{Sub: sub, UserToken: token, UserName: user, UserProfile: &u}
				})
				if err != nil {
					globals.Logger.Error(fmt.Errorf("Error scheduling minute push %s", err))
				}
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
	checkForImpactedSubscription(mess, "UPDATE")
}

func HandleFlightCreate(mess models.FlightUpdateChannelMessage) {
	checkForImpactedSubscription(mess, "CREATE")
}

func HandleFlightDelete(flt models.Flight) {
	flt.Action = "DELETE"
	checkForImpactedDeleteSubscription(flt)
	publishAllUpdatesToRabbit(flt)
}

// Check if any of the registered change subscriptions are interested in this change
func checkForImpactedSubscription(mess models.FlightUpdateChannelMessage, action string) {

	flt := GetRepo(mess.AirportCode).GetFlight(mess.FlightID)
	flt.Action = action
	sto := flt.GetSTO()

	if sto.Local().After(time.Now().Local().Add(36 * time.Hour)) {
		return
	}

	publishAllUpdatesToRabbit(*flt)

	globals.UserChangeSubscriptionsMutex.Lock()
	defer globals.UserChangeSubscriptionsMutex.Unlock()
	users := globals.GetUserProfiles()

NextSub:
	for _, sub := range globals.UserChangeSubscriptions {

		profile := models.UserProfile{}
		for _, u := range users {
			if sub.UserKey == u.Key {
				profile = u
				break
			}
		}

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
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt, UserProfile: &profile}
			continue NextSub
		}
		if sub.CreateFlight && (*flt).Action == globals.CreateAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt, UserProfile: &profile}
			continue NextSub
		}
		if !sub.DeleteFlight && (*flt).Action == globals.DeleteAction {
			continue
		}

		if sub.DeleteFlight && (*flt).Action == globals.DeleteAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt, UserProfile: &profile}
			continue NextSub
		}
		if sub.CreateFlight && (*flt).Action == globals.CreateAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt, UserProfile: &profile}
			continue NextSub
		}
		if !sub.UpdateFlight && (*flt).Action == globals.UpdateAction {
			continue
		}
		// Required Parameter Field Changes
		for _, change := range (*flt).FlightChanges.Changes {

			if globals.Contains(sub.ParameterChange, change.PropertyName) {
				changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt, UserProfile: &profile}
				continue NextSub
			}

			if (change.PropertyName == "Stand" && sub.StandChange) ||
				(change.PropertyName == "Gate" && sub.GateChange) ||
				(change.PropertyName == "CheckInCounters" && sub.CheckInChange) ||
				(change.PropertyName == "Carousel" && sub.CarouselChange) ||
				(change.PropertyName == "Chute" && sub.ChuteChange) {

				changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: flt, UserProfile: &profile}
				continue NextSub
			}

		}

		change := models.ChangePushJob{Sub: sub, Flight: flt, UserProfile: &profile}

		if sub.CheckInChange && (*flt).FlightChanges.CheckinSlotsChange != nil {
			changePushJobChannel <- change
			continue
		}
		if sub.GateChange && (*flt).FlightChanges.GateSlotsChange != nil {
			changePushJobChannel <- change
			continue
		}
		if sub.StandChange && (*flt).FlightChanges.StandSlotsChange != nil {
			changePushJobChannel <- change
			continue
		}
		if sub.ChuteChange && (*flt).FlightChanges.ChuteSlotsChange != nil {
			changePushJobChannel <- change
			continue
		}
		if sub.CarouselChange && (*flt).FlightChanges.CarouselSlotsChange != nil {
			changePushJobChannel <- change
			continue
		}

		if sub.AircraftTypeOrRegoChange && (*flt).FlightChanges.AircraftTypeChange != nil {
			changePushJobChannel <- change
			continue
		}
		if sub.AircraftTypeOrRegoChange && (*flt).FlightChanges.AircraftChange != nil {
			changePushJobChannel <- change
			continue
		}
		if sub.RouteChange && (*flt).FlightChanges.RouteChange != nil {
			changePushJobChannel <- change
			continue
		}
		if sub.LinkedFlightChange && (*flt).FlightChanges.LinkedFlightChange != nil {
			changePushJobChannel <- change
			continue
		}
	}
}
func checkForImpactedDeleteSubscription(flt models.Flight) {

	sto := flt.GetSTO()

	if sto.Local().After(time.Now().Local().Add(36 * time.Hour)) {
		return
	}

	globals.UserChangeSubscriptionsMutex.Lock()
	defer globals.UserChangeSubscriptionsMutex.Unlock()

	users := globals.GetUserProfiles()
	for _, sub := range globals.UserChangeSubscriptions {

		profile := models.UserProfile{}
		for _, u := range users {
			if sub.UserKey == u.Key {
				profile = u
				break
			}
		}

		if !sub.Enabled {
			continue
		}
		if sub.Airport != flt.GetIATAAirport() {
			continue
		}

		if sub.DeleteFlight && flt.Action == globals.DeleteAction {
			changePushJobChannel <- models.ChangePushJob{Sub: sub, Flight: &flt, UserProfile: &profile}
		}
	}
}

func executeChangePushWorker(id int, jobs <-chan models.ChangePushJob) {

	for job := range jobs {
		globals.Logger.Debug(fmt.Sprintf("Push Worker: %d Executing Change Push for User ", id))

		if !job.Sub.HTTPEnabled && !job.Sub.RMQEnabled {
			globals.Logger.Debug(fmt.Sprintf("Push Worker: %d No endpoints enabled for user ", id))
			return
		}

		file, errs := os.CreateTemp("", "changeflighttemp-*.txt")
		if errs != nil {
			fmt.Println(errs)
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
			err = os.Remove(file.Name())
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Temp file deleted for Change PUSH Send")
			}
		}()

		fwb := bufio.NewWriterSize(file, 32768)
		err := job.Flight.WriteJSON(fwb, job.UserProfile, false)
		if err != nil {
			globals.Logger.Error(fmt.Errorf("Error in executeChangePushWorker %s", err))
			return
		}
		err = fwb.Flush()
		if err != nil {
			globals.Logger.Error(fmt.Errorf("Error in executeChangePushWorker %s", err))
			return
		}
		bytesdata, err := os.ReadFile(file.Name())
		if err != nil {
			globals.Logger.Error(fmt.Errorf("Error in executeChangePushWorker %s", err))
			return
		}

		//If configured, send to Rabbit Exchange
		if job.Sub.PublishChangesRabbitMQConnectionString != "" &&
			job.Sub.PublishChangesRabbitMQExchange != "" &&
			job.Sub.PublishChangesRabbitMQTopic != "" &&
			job.Sub.RMQEnabled {
			publishChangeToRabbit(file.Name(), &job)
		}

		if job.Sub.DestinationURL == "" ||
			!job.Sub.HTTPEnabled {
			return
		}

		req, err := http.NewRequest(http.MethodPost, job.Sub.DestinationURL, bytes.NewReader(bytesdata))
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
			continue
		}
		if r == nil {
			globals.Logger.Error(fmt.Sprintf("Scheduled Push Client for user: Error making http request to: %s\n", job.Sub.DestinationURL))
			continue
		}
		if r.StatusCode != 200 {
			globals.Logger.Error(fmt.Sprintf("Change Push Client. Error making HTTP request: Returned status code = %v. URL = %s", r.StatusCode, job.Sub.DestinationURL))
			continue
		}

	}
}

func executeScheduledPushWorker(id int, jobs <-chan models.SchedulePushJob) {

	for job := range jobs {

		globals.Logger.Info(fmt.Sprintf("Executing Scheduled Push for User %s", job.UserName))

		if !job.Sub.HTTPEnabled && !job.Sub.RMQEnabled {
			globals.Logger.Debug(fmt.Sprintf("Push Worker: %d No endpoints enabled for user ", id))
			return
		}

		var fileName string
		var err error

		if strings.ToLower(job.Sub.SubscriptionType) == "flight" {
			flightresponse, serr := GetRequestedFlightsSub(job.Sub, job.UserToken)
			if serr.Err == nil {
				fileName, err = writeFlightResponseToFile(flightresponse, job.UserProfile, "-1", true)
			} else {
				err = serr.Err
			}

		} else if strings.ToLower(job.Sub.SubscriptionType) == "resource" {
			resourceresponse, serr := GetResourceSub(job.Sub, job.UserToken)
			if serr.Err != nil {
				fileName, err = writeResourceResponseToFile(resourceresponse, job.UserProfile)
			} else {
				err = serr.Err
			}

		}

		defer func() {
			err := os.Remove(fileName)
			if err != nil {
				globals.Logger.Error(fmt.Errorf("Error in executeScheduledPushWorker deleting temp file %s", err))
			}
		}()

		if err != nil {
			globals.Logger.Error(fmt.Errorf("Error in executeScheduledPushWorker %s", err))
			return
		}

		//Send to RabbitMQ if configured and enabled
		if job.Sub.PublishStatusRabbitMQConnectionString != "" &&
			job.Sub.PublishStatusRabbitMQExchange != "" &&
			job.Sub.PublishStatusRabbitMQTopic != "" &&
			job.Sub.RMQEnabled {
			publishPushToRabbit(fileName, &job)
		}
		//Send to HTTP if configured and enabled
		if job.Sub.HTTPEnabled && job.Sub.DestinationURL != "" {
			sendViaHTTPClient(fileName, &job)
		}

	}
}

func sendViaHTTPClient(fileName string, job *models.SchedulePushJob) {

	bytesdata, _ := os.ReadFile(fileName)

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

func publishChangeToRabbit(fileName string, job *models.ChangePushJob) {

	connString := job.Sub.PublishChangesRabbitMQConnectionString
	exchange := job.Sub.PublishChangesRabbitMQExchange
	routingKey := job.Sub.PublishChangesRabbitMQTopic

	if connString == "" || exchange == "" || routingKey == "" {
		return
	}

	bytesdata, _ := os.ReadFile(fileName)

	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytesdata,
		})
	failOnError(err, "Failed to publish a message")
}

func publishPushToRabbit(fileName string, job *models.SchedulePushJob) {

	connString := job.Sub.PublishStatusRabbitMQConnectionString
	exchange := job.Sub.PublishStatusRabbitMQExchange
	routingKey := job.Sub.PublishStatusRabbitMQTopic

	if connString == "" || exchange == "" || routingKey == "" {
		return
	}

	bytesdata, _ := os.ReadFile(fileName)

	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytesdata,
		})
	failOnError(err, "Failed to publish a message")
}

func publishAllUpdatesToRabbit(flt models.Flight) {

	airportCode := flt.GetIATAAirport()
	repo := GetRepo(airportCode)

	connString := repo.PublishChangesRabbitMQConnectionString
	exchange := repo.PublishChangesRabbitMQExchange
	routingKey := repo.PublishChangesRabbitMQTopic

	if connString == "" || exchange == "" || routingKey == "" {
		return
	}

	file, errs := os.CreateTemp("", "rabbitchangeflighttemp-*.json")
	if errs != nil {
		fmt.Println(errs)
		return
	}

	profile := models.UserProfile{}
	for _, u := range globals.GetUserProfiles() {
		if "default" == u.Key {
			profile = u
			break
		}
	}

	fwb := bufio.NewWriterSize(file, 32768)
	err := flt.WriteJSON(fwb, &profile, false)
	if err != nil {
		globals.Logger.Error(fmt.Errorf("Error in publishAllUpdatesToRabbit %s", err))
		return
	}
	err = fwb.Flush()
	if err != nil {
		globals.Logger.Error(fmt.Errorf("Error in publishAllUpdatesToRabbit %s", err))
		return
	}
	bytesdata, _ := os.ReadFile(file.Name())

	defer func() {
		file.Close()

		err := os.Remove(file.Name())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Temp file deleted for Rabbit Change PUSH Send")
		}
	}()

	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytesdata,
		})
	failOnError(err, "Failed to publish a message")

}
