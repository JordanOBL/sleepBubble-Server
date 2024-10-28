package server

import (
	"cmp"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

const version = "0.0.1"

type config struct {
	port int
}

type application struct {
	config config
	logger *slog.Logger
	clients map[string]*expo.PushClient
}
type Response struct {
	SleepStatus   string `json:"sleepStatus"`
	Statement string `json:"statement"`
}



func NewServer() (*http.Server, *application , *slog.Logger){
	// Create Config	
	var cfg config

	// Try to read environment variable for port (given by railway). Otherwise use default
	// Use `PORT` provided in environment or default to 3000
  	port := cmp.Or(os.Getenv("PORT"), "3000")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		intPort = 3000
	}

	// Set the port to run the API on
	cfg.port = intPort

	// create the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// create the application
	app := &application{
		config: cfg,
		logger: logger,
		clients: make(map[string]*expo.PushClient),
	}

	// create the server
	srv := &http.Server{
			Addr: fmt.Sprintf(":%d", cfg.port),
			Handler:      app.routes(),
			IdleTimeout:  45 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Using port:", "port", cfg.port)

	logger.Info("server started", "addr", srv.Addr)
	return srv, app, logger
}