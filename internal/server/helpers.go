package server

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

var AwakeSayings = []string{
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

var SleepingSayings = []string{
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

func (a *application) LoadSubscribedTokens(fname string) error {
	if len(a.clients) != 0 {
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
		a.clients[record[0]] = client

		fmt.Println(record)
	}

	fmt.Printf("Finished loading data")

	return nil
}

// The writeJSON() method is a generic helper for writing JSON to a response
func (app *application) writeJSON(w http.ResponseWriter, sCode int, data any, headers http.Header) error {
	marshalledJson, err := json.Marshal(data)

	if err != nil {
		return err
	}

	// Valid json requires newline
	marshalledJson = append(marshalledJson, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sCode)
	w.Write(marshalledJson)

	return nil
}
