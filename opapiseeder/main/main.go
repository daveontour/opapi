package main

import (
	_ "net/http/pprof"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/daveontour/opapi/opapi/timeservice"
	"github.com/daveontour/opapi/opapiseeder/cmd"
)

func main() {

	// Do a bit of initialisation
	globals.InitGlobals()
	timeservice.InitTimeService()

	//Sets up the CLI
	cmd.InitCobra()

	//Invoke the CLI
	cmd.ExecuteCobra()
}
