package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Pilladian/go-helper"
	"github.com/Pilladian/logger"
)

// GLOBAL VARS
var PORT int = 8080
var PATH string = "./infowatch"
var DATABASE_NAME string = "infowatch.db"
var DATABASE_PATH string = PATH + "/data/" + DATABASE_NAME
var LEN_FILE_ID int = 30
var BASIC_AUTH_USER = os.Getenv("BASIC_AUTH_USER")
var BASIC_AUTH_PASS = os.Getenv("BASIC_AUTH_PASS")

func initialize() {
	helper.CreatePath(PATH)
	helper.CreatePath(PATH + "/data")
	helper.CreatePath(PATH + "/logs")
}

func main() {
	// initialize environment
	initialize()
	logger.SetLogLevel(2)
	logger.SetLogFilename("./infowatch/logs/main.log")
	logger.Info("--------------------------------- Starting InfoWatch ---------------------------------")

	// http request handler
	logger.Info("setup http request handler")
	http.HandleFunc("/", rootRequestHandler)
	http.HandleFunc("/view", viewRequestHandler)
	http.HandleFunc("/healthy", healthyRequestHandler)
	http.HandleFunc("/api/v1/push", pushRequestHandler)
	fs := http.FileServer(http.Dir("./html/main"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// start web server
	logger.Info("start http server")
	server_err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	// handle web server errors
	if errors.Is(server_err, http.ErrServerClosed) {
		logger.Fatal("web server closed\n")
	} else if server_err != nil {
		logger.Fatal(fmt.Sprintf("error starting web server: %s\n", server_err))
	}
}
