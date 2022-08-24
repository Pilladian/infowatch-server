package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/Pilladian/go-helper"
)

// --------------------------------- Main ---------------------------------

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello There!</h1><p>This is the main page.</p>")
}

func healthyRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "running")
}

func uniqueRandomString(id_len int, files []fs.FileInfo) string {
	cond := true
	filename := helper.RandomString(id_len)
	for cond {
		unique := true
		for _, f := range files {
			if filename == f.Name() {
				unique = false
			}
		}
		if !unique {
			filename = helper.RandomString(id_len)
		} else {
			cond = false
		}
	}
	return filename
}

// --------------------------------- Testing ---------------------------------

func prettify(expectation string, actual string) string {
	return fmt.Sprintf("\n\t[ WANTED ] %s\n\t[ ACTUAL ] %s\n\n", expectation, actual)
}
