package server

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

//Get "/" returns status of baby sleeping
func (app *application) GetStatusHandler(writer http.ResponseWriter, r *http.Request) {
	//Request is GET
	if r.Method != "GET" {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Open database file or create new
	contents, err := os.OpenFile("sleepbubble.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf(" '/', Error reading File: %v\n", err)
	}
	defer contents.Close()

	//Read DB into csv
	csvReader := csv.NewReader(contents)
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = 1

	//Read 1st line only; gets current sleepStatus
	sleepStatus, err := csvReader.Read()
	if err != nil {
		http.Error(writer, "Error reading db csv", http.StatusInternalServerError)
	}

	if sleepStatus[0] != "0" && sleepStatus[0] != "1" {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, "invalid database value for sleep status")
	}

	// Prepare response data
	response := Response{
		SleepStatus:  sleepStatus[0],
		Statement: getRandomSaying(AwakeSayings),
	}
	if sleepStatus[0] == "1" {
		response.Statement = getRandomSaying(SleepingSayings)
	}

	// Encode the response as JSON
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}