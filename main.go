package main

import (
	"log"
	"os"

	"github.com/cpheps/internet-speed-monitor/datadog"
	"github.com/cpheps/internet-speed-monitor/speedtest"
)

func main() {
	results, err := speedtest.RunTest()
	if err != nil {
		log.Println("Encountered error:", err.Error())
		os.Exit(1)
	}

	client, err := datadog.NewClient()
	if err != nil {
		log.Println("Encountered error:", err.Error())
		os.Exit(1)
	}

	if err := client.SendTestResults(results, "wifi"); err != nil {
		log.Println("Encountered error:", err.Error())
		os.Exit(1)
	}
}
