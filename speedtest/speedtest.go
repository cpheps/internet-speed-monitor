// Package speedtest is a wrapper around the speedtest-cli
package speedtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"time"
)

// Results is partial output of the JSON speedtest-cli
type Results struct {
	Download float64   `json:"download"`
	Upload   float64   `json:"upload"`
	Ping     float64   `json:"ping"`
	Timstamp time.Time `json:"timestamp"`
}

// RunTest runs the speedtest-cli and returns results
func RunTest() (*Results, error) {
	cmd := exec.Command("speedtest-cli", "--json")

	var outBuff, errBuff bytes.Buffer
	cmd.Stdout, cmd.Stderr = &outBuff, &errBuff

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	stdErr := errBuff.String()

	if stdErr != "" {
		return nil, errors.New(stdErr)
	}

	return parseOutput(outBuff.Bytes())
}

func parseOutput(output []byte) (*Results, error) {
	var results Results
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	return &results, nil
}
