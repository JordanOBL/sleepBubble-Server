package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"math/rand"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

type Response struct {
	SleepStatus   string `json:"sleepStatus"`
	Statement string `json:"statement"`
}

var awakeSayings = []string{
	"Lennox is up and ready to rock ‘n roll—brace yourself for cuteness overload!",
	"The king has awoken! Time to serve your tiny ruler. 👑",
	"Someone’s awake and demanding snacks pronto! 🍪",
	"Get your party shoes on—Lennox is up and ready to boogie!",
	"Cue the lights! Lennox is making his grand entrance. 🕺✨",
	"Ready or not, here comes Lennox! Hide your coffee. ☕️",
	"The nap is over, and the fun (and mess) is about to begin!",
	"Watch out, world—Lennox is on the move again! 🔥",
	"Lennox is up and ready to conquer the living room. You in?",
	"And we’re back in action! Lennox, party of one, has arrived.",
	"Ryan, you better keep up—Lennox is ready to test Dad’s endurance today!",
	"Look out, Joni! Lennox is awake and about to teach grandma some dance moves. 💃",
	"Jordan, prepare yourself: Lennox is up and wants to know why you haven’t built him a castle yet.",
	"Dennis, your nap time privileges are revoked until Lennox has his snacks.",
	"India, brace yourself! Godson Lennox demands more funny faces ASAP. 🤪",
	"Robyn, the official Lennox tickling session is about to commence—bring your A-game!",
	"Ki, don’t even think about resting—Lennox is awake and needs Mom's full attention!",
	"Jasper, time to bring out the toy army! Lennox is ready for some serious playtime.",
	"Quick, someone get Lennox his snacks before he unleashes his 'Lennox Level Chaos.' 🍿",
	"Attention, family! Lennox is awake and ready for his loyal subjects to line up.",
}
var sleepingSayings = []string{
	"Don’t you dare wake the little monster… You’ve been warned! 😈",
	"Shhh… Lennox is recharging his cuteness. Disturb at your own risk!",
	"Lennox is down for the count. Go quietly—no sudden moves!",
	"The beast is resting. Use this time wisely! 😴",
	"Do not disturb: Lennox’s peaceful slumber. 🎶",
	"Silent mode activated: Lennox is sleeping (for now…)",
	"Tiptoe! Lennox is in sleep mode. Let’s keep it that way!",
	"Lennox is off dreaming…probably about world domination. Shhh.",
	"Nap time alert: Lennox is temporarily out of service!",
	"Please hold…Lennox is buffering. Do not disturb the process.",
	"Ryan, if you wake him, you’re on snack duty for the rest of the week!",
	"Jordan, wake him and you'll be the official 'Lennox Diaper Changer' for life.",
	"Joni, he’s finally down—time to sneak in that book you’ve wanted to read!",
	"Dennis, shhh! Lennox’s napping. Any noise will be met with grandpa duty.",
	"India, let him sleep! He’s recharging to show you who’s boss later.",
	"Robyn, wake him and you’ll be officially crowned the 'Lennox Entertainer'.",
	"Ki, it’s finally quiet…take a breath and pretend you’re on vacation. 😌",
	"Jasper, any noise and you’ll be sentenced to endless rounds of hide and seek.",
	"Everyone hold your breath! Lennox’s sleep is in progress—failure is not an option.",
	"Warning: Lennox is asleep. Disturbing him will activate 'The Wrath of the Family'.",
}
func getRandomSaying(sayings []string) string {
	
	return sayings[rand.Intn(len(sayings))]
}

func LoadSubscribedTokens(t map[string]*expo.PushClient, fname string) error {
	if len(t) != 0 {
		return nil
	}

	//Open db file or create if not there
	dbFile, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0777)

	if err != nil {
		return fmt.Errorf("error opening db file: %v", err)
	}

	defer dbFile.Close()
	//make dbfile csv parseable
	r := csv.NewReader(dbFile)
	//set field length for each value to one
	r.FieldsPerRecord = 1
	r.Comment = '#'

	for {
		//read each record line by line
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[0] == "0" || record[0] == "1"{
			continue
		}

		//Creare new client
		client := expo.NewPushClient(nil)
		//add record to servers token array
		t[record[0]] = client

		fmt.Println(record)
	}

	fmt.Printf("Finished loading data")

	return nil
}

func main() {

	var clients = make(map[string]*expo.PushClient) // Ensure tokens is persistent across requests

	err := LoadSubscribedTokens(clients, "../database/sleepbubble.csv")

	if err != nil {
		fmt.Printf("error loading subscribed Tokens: %v", err)
	}

	//Get "/" returns status of baby sleeping
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		contents, err := os.OpenFile("../database/sleepbubble.csv", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf(" '/', Error reading File: %v\n", err)
		}
		defer contents.Close()

		csvReader := csv.NewReader(contents)

		csvReader.Comment = '#'
		csvReader.FieldsPerRecord = 1

		sleepStatus, err := csvReader.Read()
		if err != nil {
			http.Error(writer, "Error reading db csv", http.StatusInternalServerError)
		}

		if sleepStatus[0] != "0" && sleepStatus[0] != "1" {
			writer.WriteHeader(http.StatusInternalServerError)
			io.WriteString(writer, "invalid database value for sleep status")
		}

			// Decide whether Lennox is awake or asleep (you could use request data for this)

	// Prepare response data
	response := Response{
		SleepStatus:  sleepStatus[0],
		Statement: getRandomSaying(awakeSayings),
	}
	if sleepStatus[0] == "1" {
		response.Statement = getRandomSaying(sleepingSayings)
	}

	// Encode the response as JSON
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	})

	//POST "/join" adds users eastoken to the csv db file

	http.HandleFunc("/join", func(writer http.ResponseWriter, r *http.Request) {
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
		if _, ok := clients[token]; ok {
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
		clients[string(pushToken)] = client

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
		dbFile, err := os.OpenFile("../database/sleepbubble.csv", os.O_APPEND|os.O_WRONLY, 0777)
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
	})

	http.HandleFunc("/updateSleep", func(writer http.ResponseWriter, r *http.Request) {
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
		contents, err := os.OpenFile("../database/sleepbubble.csv", os.O_RDWR|os.O_CREATE, 0644)
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
		}

		if incomingSleepStatus == prevSleepStatus[0] {
			writer.WriteHeader(http.StatusNotModified)
			io.WriteString(writer, "No Update, Sleep Staus Same As Previous")
			return
		}

		contents.Close()

    var wg sync.WaitGroup // Create a WaitGroup to manage goroutines

    for token, client := range clients {
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
				responseStatement = getRandomSaying(awakeSayings)
			} 
			if incomingSleepStatus == "1" {
				status = "Sleeping"
				responseStatement = getRandomSaying(sleepingSayings)
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
	file, err := os.Create("../database/sleepbubble.csv")
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
	for token := range clients {
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
	})

	http.ListenAndServe(":3000", nil)
	select {}
}
