package cmd

import (
	"errors"

	"os"
	"strconv"

	globals "github.com/daveontour/opapi/opapi/globals"
	version "github.com/daveontour/opapi/opapi/version"

	"github.com/spf13/cobra"
)

func InitCobra() {

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Version = version.Version

	rootCmd.AddCommand(splashCmd)
	rootCmd.AddCommand(demoCmd)
	rootCmd.AddCommand(perfTestCmd)
	rootCmd.AddCommand(parseCmd)
}
func ExecuteCobra() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "apitest",
	Short: `apitest is a CLI to seed the `,
	Long:  "\nflightresourcerestapi is a CLI to GetFlightsAPI with flights when it is run in demo mode",
}

var demoCmd = &cobra.Command{
	Use:   "demo  {number of flights to create} {number of custom properties} {appendMode bool}",
	Short: `Seed GetFlightsAPI with flights`,
	Long:  "Seed GetFlightsAPI with flights\n",
	Run: func(cmds *cobra.Command, args []string) {
		globals.IsDebug = true
		splash()
		demo(args[0], args[1], args[2])
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("Number of initial flights, custom properties not specified or append mode not set")
		}
		_, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("Invalid format or invalid number of flights entered on command line")
		}
		_, err = strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Invalid format or invalid number of custom properties entered on command line")
		}
		return nil
	},
}

var perfTestCmd = &cobra.Command{
	Use:   "perfTest  {number of flights to create} {number of custom properties} {appendMode bool}",
	Short: `Seed GetFlightsAPI with flights via RabbitMQ`,
	Long:  "Seed GetFlightsAPI with flights via RabbitMQ\n",
	Run: func(cmds *cobra.Command, args []string) {
		globals.IsDebug = true
		splash()
		perfTest(args[0], args[1], args[2])
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("Number of initial flights, custom properties or continuousUpdate flag not specified")
		}
		_, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("Invalid format or invalid number of flights entered on command line")
		}
		_, err = strconv.Atoi(args[1])
		if err != nil {
			return errors.New("Invalid format or invalid number of custom properties entered on command line")
		}
		return nil
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse  {number of flights to create} {number of custom properties} {appendMode bool}",
	Short: `Parse The memory usage`,
	Run: func(cmds *cobra.Command, args []string) {
		parse(args[0])
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Number of initial flights, custom properties not specified or append mode not set")
		}
		return nil
	},
}
