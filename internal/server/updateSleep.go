package server

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func (app *application) UpdateSleepStatus(writer http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Body == nil {
			http.Error(writer, "method error", http.StatusBadRequest)
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(writer, "Error reading incoming Body", http.StatusInternalServerError)
			return
		}
		//Parse new Sleeping Status
		incomingSleepStatus := strings.TrimSpace(string(bodyBytes))

		//Get current sleeping Status
		contents, err := os.OpenFile("/app/cmd/server/sleepbubble.csv", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf(" '/', Error reading File: %v\n", err)
		}
		defer contents.Close()

		csvReader := csv.NewReader(contents)

		csvReader.Comment = '#'
		csvReader.FieldsPerRecord = 1

		prevSleepStatus, err := csvReader.Read()
		
if err != nil {
    http.Error(writer, "Error reading db csv", http.StatusInternalServerError)
    return
}
if len(prevSleepStatus) == 0 {
    http.Error(writer, "CSV file contains no sleep status", http.StatusInternalServerError)
    return
}
		if err != nil {
			http.Error(writer, "Error reading db csv", http.StatusInternalServerError)
		}

		if incomingSleepStatus == prevSleepStatus[0] {
			writer.WriteHeader(http.StatusNotModified)
			io.WriteString(writer, "No Update, Sleep Staus Same As Previous")
			return
		}

		contents.Close()

    var wg sync.WaitGroup // Create a WaitGroup to manage goroutines

    for token, client := range app.clients {
        wg.Add(1) // Increment WaitGroup counter for each goroutine

        // Launch each task in a goroutine
        go func(token string, client *expo.PushClient) {
            defer wg.Done() // Decrement the WaitGroup counter when done
			var responseStatement string
			
            // To check the token is valid
            pushToken, err := expo.NewExponentPushToken(token)
            if err != nil {
                fmt.Printf("Error creating push token for %s: %v",token, err)

                return
            }

			var status string

			if incomingSleepStatus == "0" {
				status = "Awake"
				responseStatement = getRandomSaying(AwakeSayings)
			} 
			if incomingSleepStatus == "1" {
				status = "Sleeping"
				responseStatement = getRandomSaying(SleepingSayings)
			}

            response, err := client.Publish(
                &expo.PushMessage{
                    To:       []expo.ExponentPushToken{pushToken},
                    Body:     responseStatement,
                    Sound:    "default",
                    Title:    status,
                    Priority: expo.DefaultPriority,
                },
            )
            if err != nil {
                fmt.Println("Error publishing message:", err)
                return
            }

            // Validate responses
            if err := response.ValidateResponse(); err != nil {
                fmt.Printf("Message to %v failed: %v\n", response.PushMessage.To, err)
            } else {
                fmt.Printf("Message to %v sent successfully\n", response.PushMessage.To)
            }
		
        }(token, client) // Pass token and client as parameters to avoid closure issues
    }

    // Wait for all goroutines to finish
    wg.Wait()
    fmt.Println("All notifications processed")


	statusHeader := []string{"#This line below represents boolean value of sleeping state: 0 = awake 1 = sleeping"}

	newStatus := []string{incomingSleepStatus}

	tokenHeader := []string{"#Every line below is a Devices EAS Teoken for FCM Push notifications"}

	// Step 4: Re-open the file for writing (this will overwrite the file)
	file, err := os.Create("/app/cmd/server/sleepbubble.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Step 5: Write all modified records back to the file
	newCsvFileWriter := csv.NewWriter(file)
	err = newCsvFileWriter.Write(statusHeader)
	if err != nil {
		fmt.Println("Error writing new status header to new csv file")
	}
	err = newCsvFileWriter.Write(newStatus)
	if err != nil {
		fmt.Println("Error writing new status to new csv file")
	}
	err = newCsvFileWriter.Write(tokenHeader)
	if err != nil {
		fmt.Println("Error writing clients tokens header to new csv file")
	}
	for token := range app.clients {
		tokenArr := []string{token}
		err = newCsvFileWriter.Write(tokenArr)
		if err != nil {
			fmt.Println("Error writing tokens to new csv file")
		}
	}
	newCsvFileWriter.Flush()
	if err != nil {
		fmt.Println("Error writing CSV file:", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
	io.WriteString(writer, newStatus[0])
	
}