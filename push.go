package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Pilladian/logger"
)

// Response Codes:
//
//	  0 : OK
//	801 : Server could not check whether project_id exists
//	802 : Server could not read from given project
//	803 : Server could not write data to project
//	804 : Server could not create new project
func processData(project_id string, content string) (int, error) {
	path := PATH + "/data/" + project_id

	project_exists, err := exists(path)
	if err != nil {
		return 801, err
	}

	if !project_exists {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return 804, err
		}
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 802, err
	}

	unique_filename := uniqueRandomString(LEN_FILE_ID, files)
	f, err := os.Create(path + "/" + unique_filename + ".json")
	if err != nil {
		return 803, err
	}
	defer f.Close()

	f.WriteString(content)

	return 0, nil
}

func pushRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "HEAD" {
		content, _ := os.ReadFile("html/templates/api/v1/doc.html")
		fmt.Fprintf(w, string(content))
	} else if r.Method == "POST" {
		id := r.URL.Query()["id"][0]
		content, _ := ioutil.ReadAll(r.Body)
		response_code, err := processData(id, string(content))
		if err != nil {
			logger.Error(fmt.Sprintf("Server Response Code: %d - %s", response_code, err.Error()))
		}
	} else {
		fmt.Fprintf(w, "denied\n")
	}

}
