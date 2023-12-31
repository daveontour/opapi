package main

import (
	"net/http"

	cmd "github.com/daveontour/opapi/opapi/cmd"
	"github.com/daveontour/opapi/opapi/globals"
	gob "github.com/daveontour/opapi/opapi/gob"
	"github.com/daveontour/opapi/opapi/timeservice"

	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
)

// func initProcess() {
// 	var (
// 		inCgroup = false
// 		percent  = 0.80
// 	)
// 	go gctuner.NewTuner(inCgroup, percent)
// }

func main() {

	timeservice.InitTimeService()

	inService, err := svc.IsWindowsService()

	if err != nil {
		log.Fatalf("Failed to determine if we are running in service: %v", err)
	}

	if inService {
		cmd.RunService(globals.ConfigViper.GetString("ServiceName"), false)
		return
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	gob.GobStorageInit()

	//Sets up the CLI
	cmd.InitCobra()

	//Invoke the CLI
	cmd.ExecuteCobra()
}
