package main

import (
	cmd "github.com/daveontour/opapi/webhookclient"
)

func main() {
	cmd.InitCobraTestClient()
	cmd.ExecuteCobraTestClient()
}
