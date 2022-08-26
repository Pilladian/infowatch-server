package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Pilladian/go-helper"
	"github.com/Pilladian/logger"
)

// Response Codes:
//
//	  0 : OK
//	604 : Unable to perform queryDatabaseProject
//	605 : Unable to obtain data from database
func queryDatabaseProject(db *sql.DB, stmt string) (string, map[int64]map[string]interface{}, int, error) {
	rows, query_err := db.Query(stmt)
	if query_err != nil {
		return "", nil, 604, fmt.Errorf("Unable to perform query : %s", query_err.Error())
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	data := make(map[int64]map[string]interface{})
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if read_data_err := rows.Scan(columnPointers...); read_data_err != nil {
			return "", nil, 605, fmt.Errorf("Unable to obtain data from database : %s", read_data_err.Error())
		}

		m := make(map[string]interface{})
		var ind int64
		for i, colName := range cols {
			val := *columnPointers[i].(*interface{})
			if colName == "ID" {
				ind = val.(int64)
				continue
			}
			m[colName] = val
		}
		data[ind] = m
	}

	data_b, _ := json.MarshalIndent(data, "", " ")
	data_s := string(data_b)

	return data_s, data, 0, nil
}

// Response Codes:
//
//	  0 : OK
//	604 : Unable to perform queryDatabase
//	605 : Unable to obtain data from database
func queryDatabaseNames(db *sql.DB) ([]string, int, error) {
	rows, query_err := db.Query("SELECT name FROM sqlite_schema WHERE type IN ('table','view')	AND name NOT LIKE 'sqlite_%';")
	if query_err != nil {
		return nil, 604, fmt.Errorf("Unable to perform query : %s", query_err.Error())
	}
	defer rows.Close()

	data := []string{}
	for rows.Next() {
		var tmp string
		if read_data_err := rows.Scan(&tmp); read_data_err != nil {
			return nil, 605, fmt.Errorf("Unable to obtain table names from database : %s", read_data_err.Error())
		}
		data = append(data, tmp)
	}
	return data, 0, nil
}

// Response Codes:
//
//	  0 : OK
//	601 : Unrecognized pid format
//	602 : Unable to connect to database
//	603 : Database does not exist
//	604 : Unable to perform query
//	605 : Unable to obtain data from database
func getAllFromDatabase(project_id string) (string, int, error) {
	if project_id == "" {
		return "", 601, fmt.Errorf("Unrecognized pid format \"%s\"", project_id)
	}

	existent, existent_err := helper.Exists(DATABASE_PATH)
	if existent_err != nil {
		return "", 602, fmt.Errorf("Unable to connect to database : %s", existent_err.Error())
	}
	if !existent {
		return "", 603, fmt.Errorf("Database does not exist")
	}

	db, db_err := openDB()
	if db_err != nil {
		return "", 602, fmt.Errorf("Unable to connect to database : %s", db_err.Error())
	}
	defer db.Close()

	stmt := fmt.Sprintf("SELECT * from \"%s\"", project_id)
	data_s, _, data_s_err_code, data_s_err := queryDatabaseProject(db, stmt)
	if data_s_err != nil {
		return "", data_s_err_code, data_s_err
	}

	return data_s, 0, nil
}

func queryRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "HEAD" {
		logger.Info(fmt.Sprintf("query API received a %s Request", r.Method))
		content, os_read_err := os.ReadFile("html/templates/api/v1/query.html")
		if os_read_err != nil {
			logger.Error(fmt.Sprintf("cannot access file query.html : %s", os_read_err.Error()))
		}
		fmt.Fprintf(w, string(content))
		return
	} else if r.Method == "GET" {
		pid := r.URL.Query()["pid"]
		if len(pid) != 1 {
			logger.Error(fmt.Sprintf("query parameter \"pid\" could not be determined correctly: http://%s%s?%s", r.Host, r.URL.Path, r.URL.RawQuery))
			content, _ := os.ReadFile("html/templates/api/v1/error.html")
			fmt.Fprintf(w, fmt.Sprintf(string(content), "InfoWatch could not process your request."))
			return
		}

		if id_err := validatePID(pid[0]); id_err != nil {
			logger.Error(fmt.Sprintf("pid validation failed: %s", id_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}

		data, query_err_code, query_err := getAllFromDatabase(pid[0])
		if query_err != nil {
			logger.Error(fmt.Sprintf("Server response code \"%d\" : %s", query_err_code, query_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}
		fmt.Fprintf(w, data+"\n")
		logger.Info(fmt.Sprintf("successfully queried data from server - project : %s", pid))
	} else {
		logger.Warning(fmt.Sprintf("query API received a %s Request", r.Method))
		fmt.Fprintf(w, "denied\n")
	}
}
