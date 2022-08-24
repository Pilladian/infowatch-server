package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Pilladian/logger"
)

// Response Codes:
//
//	  0 : OK
//	601 : Could not read content of provided project
//	602 : Could not read data from specific file inside the project
func queryAllData(id string) (string, int, error) {
	path := PATH + "/data/" + id + "/"
	data := make(map[string]map[string]interface{})

	files, dir_read_err := ioutil.ReadDir(path)
	if dir_read_err != nil {
		return "", 601, dir_read_err
	}

	for i, f := range files {
		filename := f.Name()
		if filename == "schema.json" {
			continue
		}
		data_f_b, data_file_err := ioutil.ReadFile(path + filename)
		if data_file_err != nil {
			return "", 602, errors.New(fmt.Sprintf("Could not read data from file \"%s\" : %s", filename, data_file_err))
		}
		var data_f_content map[string]interface{}
		json.Unmarshal(data_f_b, &data_f_content)
		data[fmt.Sprintf("%d", i)] = data_f_content
	}

	data_b, _ := json.MarshalIndent(data, "", " ")
	data_s := string(data_b)

	return data_s, 0, nil
}

func queryRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "HEAD" {
		content, os_read_err := os.ReadFile("html/templates/api/v1/query.html")
		if os_read_err != nil {
			logger.Error(fmt.Sprintf("cannot access file query.html : %s", os_read_err.Error()))
		}
		fmt.Fprintf(w, string(content))
		return
	} else if r.Method == "GET" {
		id := r.URL.Query()["id"]
		if len(id) != 1 {
			logger.Error(fmt.Sprintf("query parameter \"id\" could not be determined correctly: http://%s%s?%s", r.Host, r.URL.Path, r.URL.RawQuery))
			content, _ := os.ReadFile("html/templates/api/v1/error.html")
			fmt.Fprintf(w, fmt.Sprintf(string(content), "InfoWatch could not process your request."))
			return
		}

		if id_err := validateID(id[0]); id_err != nil {
			logger.Error(fmt.Sprintf("ID validation failed: %s", id_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}

		data, query_err_code, query_err := queryAllData(id[0])
		if query_err != nil {
			logger.Error(fmt.Sprintf("Server response code \"%d\" : %s", query_err_code, query_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}
		fmt.Fprintf(w, data+"\n")
	}
}
