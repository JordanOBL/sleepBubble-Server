package main

import (
	"fmt"
	"os"

	"github.com/jordanOBL/sleepBubble/internal/server"
)

func main() {
	//Create New Server
	server, app, logger := server.NewServer()

	//Load Db of connected ExpotTokens 
	//into servers Client List
	err := app.LoadSubscribedTokens("../database/sleepbubble.csv")
	if err != nil {
		fmt.Printf("error loading subscribed Tokens: %v", err)
	}

	// Start the server
	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
	select {}
}
