package main

import (
	"fmt"
	"net/http"
)

// --------------------------------- Main ---------------------------------

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello There!</h1><p>This is the main page.</p>")
}

func healthyRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "running")
}

// --------------------------------- Testing ---------------------------------

func prettify(expectation string, actual string) string {
	return fmt.Sprintf("\n\t[ WANTED ] %s\n\t[ ACTUAL ] %s\n\n", expectation, actual)
}
