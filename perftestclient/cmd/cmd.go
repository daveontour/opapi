package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/daveontour/opapi/opapi/globals"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "perftestclient.exe {url} {fixed interval true or flase} {interval} {CSV logfile name}",
	Short: `Various set of utilities`,
	Run: func(cmds *cobra.Command, args []string) {
		execute(args)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 4 {
			return errors.New("Need to specify the parameters. See usage")
		}
		return nil
	},
}

func InitCobra() {

	rootCmd.CompletionOptions.DisableDefaultCmd = false
}
func ExecuteCobra() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func execute(args []string) {

	requestURL := args[0]
	fixed, err := strconv.ParseBool(args[1])
	if err != nil {
		fixed = true
	}
	interval, err := strconv.Atoi(args[2])
	if err != nil {
		interval = 10
	}
	fileName := args[3]

	fmt.Println(fmt.Sprintf("url:%s \nfixed:%t \ninterval:%d \nFilename:%s", requestURL, fixed, interval, fileName))
	initCSVHeader(fileName)

	for {
		waitInterval := interval
		if !fixed {
			waitInterval = int(rand.Intn(interval))
		}

		time.Sleep(time.Duration(waitInterval * int(time.Second)))

		body, duration, _ := makeRequest(requestURL)
		_, processID, memUsage := getmemusage("opapi")

		tf := -1
		var mess interface{}
		fn := -1
		err := json.Unmarshal(body, &mess)
		if err != nil {
			fn = -1
		} else {
			m, ok := mess.(map[string]interface{})
			if !ok {
				_ = fmt.Errorf("want type map[string]interface{};  got %T", mess)
			}
			for k, v := range m {
				if k == "NumberOfFlights" {
					fn, _ = strconv.Atoi(fmt.Sprint(v))
					break
				}
			}
			for k, v := range m {
				if k == "TotalFlights" {
					tf, _ = strconv.Atoi(fmt.Sprint(v))
					break
				}
			}
		}

		logusage(duration, fn, tf, memUsage, processID, fileName)

	}

}

func getmemusage(process string) (processName, processID, memUsage string) {
	cmd := exec.Command("tasklist.exe", "/fo", "csv", "/nh")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return
	}
	reader := bytes.NewReader(out)
	csvReader := csv.NewReader(reader)
	data, _ := csvReader.ReadAll()

	for _, line := range data {
		processName = line[0]
		processID = line[1]
		mem := line[4]

		if strings.Contains(processName, process) || strings.Contains(processID, process) {
			memUsage = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(mem, ",", ""), " ", ""), "K", "")
			return
		}
	}
	return
}

func makeRequest(requestURL string) (body []byte, ts time.Duration, err error) {

	t1 := time.Now()
	resp, err := http.Get(requestURL)
	t2 := time.Now()
	ts = t2.Sub(t1)

	if err != nil {
		body = nil
		return
	} else {
		defer resp.Body.Close()

		// body, err = httputil.DumpResponse(resp, true)
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			globals.Logger.Error(fmt.Sprintf("client: could not read response body: %s\n", err))
		}

		return
	}
}

func initCSVHeader(fileName string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%s,%s,%s,%s, %s,%s\n", "Time", "API Execution Time", "Num Flights", "Total Flights", "Memory Usage", "Process ID"))
	f.Close()
}

func logusage(ts time.Duration, numFlights int, tf int, memUsage string, processID string, fileName string) {

	now := time.Now()

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	ms := int64(ts / time.Millisecond)

	fmt.Printf("%s,%v,%v,%v,%s,%s\n", now.Format(time.RFC3339), ms, numFlights, tf, memUsage, processID)
	f.WriteString(fmt.Sprintf("%s,%v,%v,%v,%s,%s\n", now.Format(time.RFC3339), ms, numFlights, tf, memUsage, processID))
	f.Close()
}
