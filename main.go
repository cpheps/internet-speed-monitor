package main

import (
	"log"

	"github.com/cpheps/internet-speed-monitor/speedtest"
)

func main() {
	results, err := speedtest.RunTest()
	if err != nil {
		log.Println("Encountered error:", err.Error())
		return
	}

	log.Printf("Results: %+v", results)
}
