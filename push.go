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
//	801 : Server could not check whether project_id exists
//	802 : Server could not read from given project
//	803 : Server could not write data to project
//	804 : Server could not create new project
//	805 : Unrecognized ID format
//	806 : Unknown type in json data
//	807 : Provided json data could not be parsed
func processData(project_id string, content string) (int, error) {
	if project_id == "" {
		return 805, errors.New(fmt.Sprintf("Unrecognized ID format \"%s\"", project_id))
	}

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

	if !project_exists {
		var json_content map[string]interface{}
		json.Unmarshal([]byte(content), &json_content)
		if json_content == nil {
			return 807, errors.New(fmt.Sprintf("Provided json data %s could not be parsed", content))
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
		// TODO: check content for schema
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
		id := r.URL.Query()["id"]
		if len(id) != 1 {
			logger.Error(fmt.Sprintf("query parameter \"id\" could not be determined correctly: http://%s%s?%s", r.Host, r.URL.Path, r.URL.RawQuery))
			content, _ := os.ReadFile("html/templates/api/v1/error.html")
			fmt.Fprintf(w, fmt.Sprintf(string(content), "InfoWatch could not process your request."))
			return
		}
		data, _ := ioutil.ReadAll(r.Body)

		if id_err := validateID(id[0]); id_err != nil {
			logger.Error(fmt.Sprintf("ID validation failed: %s", id_err.Error()))
			return
		}

		response_code, err := processData(id[0], string(data))
		if err != nil {
			logger.Error(fmt.Sprintf("Server Response Code: %d - %s", response_code, err.Error()))
		}
	} else {
		fmt.Fprintf(w, "denied\n")
	}

}
