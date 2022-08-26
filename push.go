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
	_ "github.com/mattn/go-sqlite3"
)

// Response Codes:
//
//	  0 : OK
//	801 : Unrecognized pid format
//	802 : Unable to connect to database
//	803 : Unable to create database
//	804 : Unknown type in json data
//	805 : Unable to create table in database
//	806 : Unable to store data in database
func processData(project_id string, content string) (int, error) {
	if project_id == "" {
		return 801, fmt.Errorf("Unrecognized pid format \"%s\"", project_id)
	}

	var json_content map[string]interface{}
	json_parse_err := json.Unmarshal([]byte(content), &json_content)
	if json_content == nil {
		return 807, fmt.Errorf("Provided json data %s could not be parsed : %s", content, json_parse_err.Error())
	}

	existent, existent_err := helper.Exists(DATABASE_PATH)
	if existent_err != nil {
		return 802, fmt.Errorf("Unable to connect to database : %s", existent_err.Error())
	}

	if !existent {
		f, db_create_err := os.Create(DATABASE_PATH)
		if db_create_err != nil {
			return 803, fmt.Errorf("Unable to create database : %s", db_create_err.Error())
		}
		f.Close()
	}

	db, db_err := openDB()
	if db_err != nil {
		return 802, fmt.Errorf("Unable to connect to database : %s", db_err.Error())
	}
	defer db.Close()

	create_table_stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\"(ID INTEGER PRIMARY KEY AUTOINCREMENT", project_id)
	for k, v := range json_content {
		switch v.(type) {
		case string:
			create_table_stmt += fmt.Sprintf(", %s TEXT", k)
		case int:
			create_table_stmt += fmt.Sprintf(", %s INTEGER", k)
		case float32:
			create_table_stmt += fmt.Sprintf(", %s INTEGER", k)
		case float64:
			create_table_stmt += fmt.Sprintf(", %s INTEGER", k)
		default:
			return 804, errors.New("Unknown type for json data")
		}
	}
	create_table_stmt += ");"
	_, create_table_err := db.Exec(create_table_stmt)
	if create_table_err != nil {
		return 805, fmt.Errorf("Unable to create table in database : %s", create_table_err.Error())
	}

	store_data_stmt := fmt.Sprintf("INSERT INTO \"%s\"(", project_id)
	tmp1 := ""
	tmp2 := " values("
	for k, v := range json_content {
		switch v.(type) {
		case string:
			tmp1 += fmt.Sprintf("%s, ", k)
			tmp2 += fmt.Sprintf("\"%s\", ", v)
		case int:
			tmp1 += fmt.Sprintf("%s, ", k)
			tmp2 += fmt.Sprintf("%d, ", v)
		case float32:
			tmp1 += fmt.Sprintf("%s, ", k)
			tmp2 += fmt.Sprintf("%f, ", v)
		case float64:
			tmp1 += fmt.Sprintf("%s, ", k)
			tmp2 += fmt.Sprintf("%f, ", v)
		default:
			return 804, errors.New("Unknown type for json data")
		}
	}
	store_data_stmt += fmt.Sprintf("%s) %s);", tmp1[:len(tmp1)-2], tmp2[:len(tmp2)-2])
	_, store_data_err := db.Exec(store_data_stmt)
	if store_data_err != nil {
		return 806, fmt.Errorf("Unable to store data in database : %s", store_data_err.Error())
	}

	return 0, nil
}

func pushRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "HEAD" {
		logger.Info(fmt.Sprintf("push API received a %s Request", r.Method))
		content, os_read_err := os.ReadFile("html/templates/api/v1/push.html")
		if os_read_err != nil {
			logger.Error(fmt.Sprintf("cannot access file push.html : %s", os_read_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}
		fmt.Fprintf(w, string(content))
		return
	} else if r.Method == "POST" {
		pid := r.URL.Query()["pid"]
		if len(pid) != 1 {
			logger.Error(fmt.Sprintf("query parameter \"pid\" could not be determined correctly: http://%s%s?%s", r.Host, r.URL.Path, r.URL.RawQuery))
			content, _ := os.ReadFile("html/templates/api/v1/error.html")
			fmt.Fprintf(w, fmt.Sprintf(string(content), "InfoWatch could not process your request."))
			return
		}
		data, _ := ioutil.ReadAll(r.Body)

		if pid_err := validatePID(pid[0]); pid_err != nil {
			logger.Error(fmt.Sprintf("pid validation failed: %s", pid_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}

		if data_err := validateData(string(data)); data_err != nil {
			logger.Error(fmt.Sprintf("Data validation failed: %s", data_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}

		response_code, err := processData(pid[0], string(data))
		if err != nil {
			logger.Error(fmt.Sprintf("Server Response Code: %d - %s", response_code, err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}
		logger.Info(fmt.Sprintf("successfully pushed data to server - project : %s", pid))
		fmt.Fprintf(w, "success\n")
	} else {
		logger.Warning(fmt.Sprintf("push API received a %s Request", r.Method))
		fmt.Fprintf(w, "denied\n")
	}
}
