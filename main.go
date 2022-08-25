package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Pilladian/go-helper"
	"github.com/Pilladian/logger"
)

// GLOBAL VARS
var PORT int = 8080
var PATH string = "./infowatch"
var DATABASE_NAME string = "infowatch.db"
var LEN_FILE_ID int = 30

func initialize() {
	logger.SetLogLevel(2)
	helper.CreatePath(PATH)
	helper.CreatePath(PATH + "/data")
	helper.CreatePath(PATH + "/logs")
	logger.SetLogFilename("./infowatch/logs/main.log")
}

func main() {
	// initialize environment
	initialize()

	// http request handler
	http.HandleFunc("/", rootRequestHandler)
	http.HandleFunc("/healthy", healthyRequestHandler)
	http.HandleFunc("/api/v1/push", pushRequestHandler)
	http.HandleFunc("/api/v1/query", queryRequestHandler)

	// start web server
	server_err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	// handle web server errors
	if errors.Is(server_err, http.ErrServerClosed) {
		logger.Fatal("web server closed\n")
	} else if server_err != nil {
		logger.Fatal(fmt.Sprintf("error starting web server: %s\n", server_err))
	}
}
