package server

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

//POST "/join" adds users eastoken to the csv db file
func (app *application) JoinServerHandler(writer http.ResponseWriter, r *http.Request){

		if r.Method != "POST" {
			http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Body == nil {
			http.Error(writer, "invalid request to endpoint", http.StatusBadRequest)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(writer, "error reading body", http.StatusInternalServerError)
			return
		}

		token := strings.TrimSpace(string(bodyBytes)) // Trim any whitespace from the token

		// Check if the token length is exactly 22
		if len(token) != 41 {
			http.Error(writer, "invalid token in body", http.StatusNotAcceptable)
			return
		}

		// Check if the token already exists in the map
		if _, ok := app.clients[token]; ok {
			http.Error(writer, "Already subscribed", http.StatusOK)
			return
		}

		// To check the token is valid
		pushToken, err := expo.NewExponentPushToken(token)
		if err != nil {
			panic(err)
		}
		// Create a new Expo SDK client
		client := expo.NewPushClient(nil)
		//Add to clients map
		app.clients[string(pushToken)] = client

		// Publish message
		response, err := client.Publish(
			&expo.PushMessage{
				To:       []expo.ExponentPushToken{pushToken},
				Body:     "Lennox thanks you for letting him dream!",
				Sound:    "default",
				Title:    "Sleep Bubble",
				Priority: expo.DefaultPriority,
			},
		)

		// Check errors
		if err != nil {
			panic(err)
		}

		// Validate responses
		if response.ValidateResponse() != nil {
			fmt.Println(response.PushMessage.To, "failed")
		}

		// Prepare record to write to the CSV
		recordToWrite := []string{token}

		// Open CSV file with append and write permissions
		dbFile, err := os.OpenFile("/app/cmd/server/sleepbubble.csv", os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			http.Error(writer, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer dbFile.Close()

		// Write to the CSV file
		csvDbFile := csv.NewWriter(dbFile)
		if err := csvDbFile.Write(recordToWrite); err != nil {
			http.Error(writer, "Error writing to CSV", http.StatusInternalServerError)
			return
		}
		csvDbFile.Flush() // Ensure the data is written

		writer.WriteHeader(http.StatusOK)
		io.WriteString(writer, "Joined")
	
}