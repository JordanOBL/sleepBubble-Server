package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jordanOBL/sleepBubble/internal/server"
)
func main() {
    // Create New Server
    server, app, logger := server.NewServer()

    // Load DB of connected Expo tokens into the server's client list
    err := app.LoadSubscribedTokens("sleepbubble.csv")
    if err != nil {
        logger.Error("error loading subscribed tokens", "err", err.Error())
    }

    // Start the server
    logger.Info("Starting server...")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.Error("Server error", "err", err.Error())
        os.Exit(1)
    }
}

