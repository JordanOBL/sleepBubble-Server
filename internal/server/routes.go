package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	//Get "/sleepstatus" returns status of baby sleeping
	router.HandlerFunc(http.MethodGet,"/sleepstatus", app.GetStatusHandler )

	//Join the server with Expo Token to get push notifications
	router.HandlerFunc(http.MethodPost,"/join", app.JoinServerHandler )

	//Update sleep status of the baby
	router.HandlerFunc(http.MethodPost,"/updateSleep", app.UpdateSleepStatus)

	// Define the available routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return router
}