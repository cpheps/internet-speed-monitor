package main

import (
	"log"
	"os"

	"github.com/cpheps/internet-speed-monitor/datadog"
	"github.com/cpheps/internet-speed-monitor/speedtest"
)

const classifierEnv = "CLASSIFIER"

func main() {
	results, err := speedtest.RunTest()
	if err != nil {
		log.Println("Encountered error:", err.Error())
		os.Exit(1)
	}

	classifier := os.Getenv(classifierEnv)

	client, err := datadog.NewClient()
	if err != nil {
		log.Println("Encountered error:", err.Error())
		os.Exit(1)
	}

	if err := client.SendTestResults(results, classifier); err != nil {
		log.Println("Encountered error:", err.Error())
		os.Exit(1)
	}
}
