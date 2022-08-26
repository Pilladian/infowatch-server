package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"

	"github.com/Pilladian/logger"
)

func getQueryForTemplate(data_json map[int64]map[string]interface{}) Query {
	q := Query{}
	m := []string{"ID"}
	for i := range data_json[int64(1)] {
		m = append(m, i)
	}
	sort.Strings(m)
	q.Header = m

	m2 := [][]interface{}{}
	for a, b := range data_json {
		tmp := []interface{}{a}
		for _, i := range m {
			tmp = append(tmp, b[i])
		}
		m2 = append(m2, tmp)
	}
	q.Rows = sortInterfaceSlice(m2, 0)
	return q
}

func viewRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" || r.Method == "HEAD" {
		logger.Warning(fmt.Sprintf("push API received a %s Request", r.Method))
		fmt.Fprintf(w, "denied\n")
	} else if r.Method == "GET" {
		logger.Info("test")
		pid := r.URL.Query()["pid"]
		if len(pid) != 1 {
			logger.Error(fmt.Sprintf("query parameter \"pid\" could not be determined correctly: http://%s%s?%s", r.Host, r.URL.Path, r.URL.RawQuery))
			content, _ := os.ReadFile("html/templates/api/v1/error.html")
			fmt.Fprintf(w, fmt.Sprintf(string(content), "InfoWatch could not process your request."))
			return
		}

		if pid_err := validatePID(pid[0]); pid_err != nil {
			logger.Error(fmt.Sprintf("pid validation failed: %s", pid_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}

		t, _ := template.ParseFiles("html/templates/root/results.html")

		db, db_err := openDB()
		if db_err != nil {
			logger.Error(db_err.Error())
			fmt.Fprintf(w, "error\n")
			return
		}
		defer db.Close()

		_, data, query_err_code, query_err := queryDatabaseProject(db, fmt.Sprintf("SELECT * FROM \"%s\";", pid[0]))
		if query_err != nil {
			logger.Error(fmt.Sprintf("Server response code \"%d\" : %s", query_err_code, query_err.Error()))
			fmt.Fprintf(w, "error\n")
			return
		}
		q := getQueryForTemplate(data)
		t.Execute(w, q)
	} else {
		logger.Warning(fmt.Sprintf("push API received a %s Request", r.Method))
		fmt.Fprintf(w, "denied\n")
	}
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/templates/root/index.html")

	db, db_err := openDB()
	if db_err != nil {
		logger.Error(db_err.Error())
		fmt.Fprintf(w, "error\n")
		return
	}
	defer db.Close()

	table_names, query_err_code, query_err := queryDatabaseNames(db)
	if query_err != nil {
		logger.Error(fmt.Sprintf("Server response code \"%d\" : %s", query_err_code, query_err.Error()))
		fmt.Fprintf(w, "error\n")
		return
	}
	t.Execute(w, table_names)
}
