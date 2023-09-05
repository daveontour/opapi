package cmd

import (
	"fmt"

	"runtime"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/repo"
	"github.com/daveontour/opapi/opapi/server"
)

func runProgram() {

	numCPU := runtime.NumCPU()

	globals.Logger.Debug(fmt.Sprintf("Number of cores available = %v", numCPU))

	runtime.GOMAXPROCS(runtime.NumCPU())
	//Wait group so the program doesn't exit
	globals.Wg.Add(1)

	// The HTTP Server
	go server.StartGinServer(false)

	// Handler for the different types of messages passed by channels
	go eventMonitor()

	// Manages the population and update of the repositoiry of flights
	go repo.InitRepositories()

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
	globals.Wg.Wait()
}
