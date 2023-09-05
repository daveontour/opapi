package cmd

import (
	"fmt"

	"runtime"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/repo"
	"github.com/daveontour/opapi/opapi/server"
)

func perfTest() {

	// Start the system in performance test mode. Resources and flights are created as per test.json
	// Requires Rabbit MQ to be running. Messages are passsed via Rabbit MQ
	numCPU := runtime.NumCPU()
	globals.Logger.Debug(fmt.Sprintf("Number of cores available = %v", numCPU))
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Wait group so the program doesn't exit
	globals.Wg.Add(1)

	// The HTTP Server
	go server.StartGinServer(true)

	// Handler for the different types of messages passed by channels
	go eventMonitor()

	// Initiate the User Change Subscriptions
	globals.UserChangeSubscriptionsMutex.Lock()
	for _, up := range globals.GetUserProfiles() {
		if up.Enabled {
			ucs := up.UserChangeSubscriptions
			userKey := up.Key
			for i := range ucs {
				ucs[i].UserKey = userKey
			}
			globals.UserChangeSubscriptions = append(globals.UserChangeSubscriptions, ucs...)
		}
	}
	globals.UserChangeSubscriptionsMutex.Unlock()
	go repo.SchedulePushes("APT", false)
	repo.StartChangePushWorkerPool(globals.ConfigViper.GetInt("NumberOfChangePushWorkers"))
	repo.PerfTestInit()
	globals.Wg.Wait()
}
