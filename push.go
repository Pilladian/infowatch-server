package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Pilladian/go-helper"
	"github.com/Pilladian/logger"
)

// Response Codes:
//
//	  0 : OK
//	801 : Server could not check whether project_id exists
//	802 : Server could not read from given project
//	803 : Server could not write data to project
//	804 : Server could not create new project
//	805 : Unrecognized ID format
//	806 : Unknown type in json data
//	807 : Provided json data could not be parsed
//	808 : Schema file could not be opened
//	809 : Provided JSON data has invalid schema
func processData(project_id string, content string) (int, error) {
	if project_id == "" {
		return 805, errors.New(fmt.Sprintf("Unrecognized ID format \"%s\"", project_id))
	}

	path := PATH + "/data/" + project_id

	project_exists, err := helper.Exists(path)
	if err != nil {
		return 801, err
	}

	if !project_exists {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return 804, err
		}
	}

	if !project_exists {
		var json_content map[string]interface{}
		json_parse_err := json.Unmarshal([]byte(content), &json_content)
		if json_content == nil {
			return 807, errors.New(fmt.Sprintf("Provided json data %s could not be parsed : %s", content, json_parse_err.Error()))
		}

		json_schema := make(map[string]interface{})
		for k, v := range json_content {
			switch v.(type) {
			case string:
				json_schema[k] = ""
			case int:
				json_schema[k] = 0
			case float32:
				json_schema[k] = 0
			case float64:
				json_schema[k] = 0
			default:
				return 806, errors.New("Unknown type for json data")
			}
		}

		file, _ := json.MarshalIndent(json_schema, "", " ")
		_ = ioutil.WriteFile(path+"/schema.json", file, 0700)

	} else {
		valid_schema_err := validateSchema(path, content)
		if valid_schema_err != nil {
			return 809, valid_schema_err
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
		content, os_read_err := os.ReadFile("html/templates/api/v1/doc.html")
		if os_read_err != nil {
			logger.Error(fmt.Sprintf("cannot access file doc.html : %s", os_read_err.Error()))
		}
		fmt.Fprintf(w, string(content))
		return
	} else if r.Method == "POST" {
		id := r.URL.Query()["id"]
		if len(id) != 1 {
			logger.Error(fmt.Sprintf("query parameter \"id\" could not be determined correctly: http://%s%s?id=SOME_ID", r.Host, r.URL.Path))
			content, _ := os.ReadFile("html/templates/api/v1/error.html")
			fmt.Fprintf(w, fmt.Sprintf(string(content), "InfoWatch could not process your request."))
			return
		}
		data, _ := ioutil.ReadAll(r.Body)

		if id_err := validateID(id[0]); id_err != nil {
			logger.Error(fmt.Sprintf("ID validation failed: %s", id_err.Error()))
			return
		}

		if data_err := validateData(string(data)); data_err != nil {
			logger.Error(fmt.Sprintf("Data validation failed: %s", data_err.Error()))
			return
		}

		response_code, err := processData(id[0], string(data))
		if err != nil {
			logger.Error(fmt.Sprintf("Server Response Code: %d - %s", response_code, err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}
	} else {
		fmt.Fprintf(w, "denied\n")
	}
}
