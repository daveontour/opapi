package cmd

import (
	"runtime"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/repo"
	"github.com/daveontour/opapi/opapi/server"
)

func demo() {

	// Start the system in demo mode. Resources and flights are created as per test.json
	// Does not require Rabbit MQ to be running.
	globals.DemoMode = true

	runtime.GOMAXPROCS(runtime.NumCPU())
	globals.Wg.Add(1)
	go server.StartGinServer(true)
	go eventMonitor()

	// // Initiate the User Change Subscriptions
	globals.UserChangeSubscriptionsMutex.Lock()
	for _, up := range globals.GetUserProfiles() {
		if up.Enabled {
			ucs := up.UserChangeSubscriptions
			userKey := up.Key
			for i, _ := range ucs {
				ucs[i].UserKey = userKey
			}
			globals.UserChangeSubscriptions = append(globals.UserChangeSubscriptions, ucs...)
		}
	}
	globals.UserChangeSubscriptionsMutex.Unlock()

	repo.StartChangePushWorkerPool(globals.ConfigViper.GetInt("NumberOfChangePushWorkers"))
	repo.PerfTestInit()

	repo.SchedulePushes("APT", true)
	globals.Wg.Wait()
}
