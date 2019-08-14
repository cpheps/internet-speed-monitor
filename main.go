package main

import (
	"log"

	"github.com/cpheps/internet-speed-monitor/sheets"
	"github.com/cpheps/internet-speed-monitor/speedtest"
)

func main() {
	sheetsClient, err := sheets.NewClient()
	if err != nil {
		log.Fatalln("Failed created client for Google Sheets;", err.Error())
	}

	results, err := speedtest.RunTest()
	if err != nil {
		log.Fatalln("Encountered error:", err.Error())
	}

	if err = sheetsClient.SubmitTestResults(results); err != nil {
		log.Fatalln("Failed submitting test results:", err.Error())
	}

	log.Printf("Results: %+v", results)
}
