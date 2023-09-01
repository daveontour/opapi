package cmd

import (
	"fmt"

	version "github.com/daveontour/opapi/opapi/version"
	"github.com/spf13/cobra"
)

var splashCmd = &cobra.Command{
	Use:   "splash",
	Short: `Shows the Splash text`,
	Long:  `Shows the Splash text`,
	Run: func(cmd *cobra.Command, args []string) {
		splash()
	},
}

func splash() {
	fmt.Println()
	fmt.Println("*******************************************************")
	fmt.Println("*                                                     *")
	fmt.Println("*  AMS Flights and Resources Rest API (" + version.Version + ")         * ")
	fmt.Println("*                                                     *")
	fmt.Println("*  (This is NOT official SITA Software)               *")
	fmt.Println("*  (Community Contributed Software)                   *")
	fmt.Println("*                                                     *")
	fmt.Println("*  Seeds the API with Demo flights when the API       *")
	fmt.Println("*  is running in Demo Mode                            *")
	fmt.Println("*                                                     *")
	fmt.Println("*******************************************************")
	fmt.Println()
}
