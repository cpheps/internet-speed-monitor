package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cpheps/internet-speed-monitor/speedtest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	credentialEnv = "SHEETS_CRED"
	sheetsIDEnv   = "SHEETS_ID"
	sheetNameEnv  = "SHEET_NAME"
)

// Client connection to google sheets
type Client struct {
	client    *sheets.Service
	sheetID   string
	sheetName string
}

// NewClient creates a new sheets.Client base on credential file. Looks for file in SHEETS_CRED env
// if the env is not present looks locally for credentials.json.
func NewClient() (*Client, error) {
	credData, err := readCredFile()
	if err != nil {
		return nil, err
	}

	client, err := createOAuthClient(credData)
	if err != nil {
		return nil, err
	}

	srv, err := sheets.New(client)
	if err != nil {
		return nil, err
	}

	sheetID, err := getSheetID()
	if err != nil {
		return nil, err
	}

	sheetName, err := getSheetName()
	if err != nil {
		return nil, err
	}

	return &Client{
		client:    srv,
		sheetID:   *sheetID,
		sheetName: *sheetName,
	}, nil
}

// SubmitTestResults submits test results to google sheet
func (c *Client) SubmitTestResults(results *speedtest.Results) error {
	values := &sheets.ValueRange{
		Values: [][]interface{}{
			{
				results.TimestampEST(),
				results.DownloadMbps(),
				results.UploadMbps(),
				results.Ping,
			},
		},
	}
	_, err := c.client.Spreadsheets.Values.Append(c.sheetID, c.sheetName, values).
		InsertDataOption("INSERT_ROWS").
		ValueInputOption("RAW").
		Context(context.Background()).
		Do()

	return err
}

func readCredFile() ([]byte, error) {
	credFile, found := os.LookupEnv(credentialEnv)
	if !found {
		credFile = "credentials.json"
	}

	return ioutil.ReadFile(credFile)
}

func getSheetID() (*string, error) {
	sheetID, ok := os.LookupEnv(sheetsIDEnv)
	if !ok {
		return nil, fmt.Errorf("%s was not set", sheetsIDEnv)
	}

	return &sheetID, nil
}

func getSheetName() (*string, error) {
	sheetName, ok := os.LookupEnv(sheetNameEnv)
	if !ok {
		return nil, fmt.Errorf("%s was not set", sheetNameEnv)
	}

	return &sheetName, nil
}

func createOAuthClient(credData []byte) (*http.Client, error) {
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credData, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	return getClient(config)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tokFile, tok); err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Code modified from Google example seen here https://developers.google.com/sheets/api/quickstart/go

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
