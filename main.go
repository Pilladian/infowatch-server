package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Pilladian/logger"
)

// GLOBAL VARS
var PORT int = 8080
var PATH string = "./infowatch"
var LEN_FILE_ID int = 30

func initialize() {
	createPath(PATH)
	createPath(PATH + "/data")
	createPath(PATH + "/logs")
}

func main() {
	// initialize environment
	initialize()

	// http request handler
	http.HandleFunc("/", rootRequestHandler)
	http.HandleFunc("/api/v1/push", pushRequestHandler)

	// start web server
	server_err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	// handle web server errors
	if errors.Is(server_err, http.ErrServerClosed) {
		logger.Error("web server closed\n")
		os.Exit(1)
	} else if server_err != nil {
		logger.Error(fmt.Sprintf("error starting web server: %s\n", server_err))
		os.Exit(1)
	}
}
