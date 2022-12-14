package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Pilladian/logger"
)

// ----------------------------------------------------------------------------------
// ---------------------------------  /api/v1/push  ---------------------------------
// ----------------------------------------------------------------------------------

func pushRequestHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		logger.Warning(fmt.Sprintf("Authentication Attempt: %s with password %s", username, password))
		expectedUsernameHash := sha256.Sum256([]byte(BASIC_AUTH_USER))
		expectedPasswordHash := sha256.Sum256([]byte(BASIC_AUTH_PASS))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
			logger.Info("Successfully logged in")
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
				logger.Info(fmt.Sprintf("Successfully pushed data to server - project : %s", pid[0]))
				fmt.Fprintf(w, "success\n")
			} else {
				logger.Warning(fmt.Sprintf("Push API received a %s Request", r.Method))
				fmt.Fprintf(w, "denied\n")
			}
			return
		} else {
			logger.Error("Wrong credentails")
		}
	} else {
		logger.Error("No basic authentication was provided")
	}
	fmt.Fprintf(w, "denied\n")
}

// ----------------------------------------------------------------------------------
// --------------------------------- /api/v1/query  ---------------------------------
// ----------------------------------------------------------------------------------

func queryRequestHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		logger.Warning(fmt.Sprintf("Authentication Attempt: %s with password %s", username, password))
		expectedUsernameHash := sha256.Sum256([]byte(BASIC_AUTH_USER))
		expectedPasswordHash := sha256.Sum256([]byte(BASIC_AUTH_PASS))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
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
				logger.Info(fmt.Sprintf("successfully queried data from server - project : %s", pid[0]))
			} else {
				logger.Warning(fmt.Sprintf("query API received a %s Request", r.Method))
				fmt.Fprintf(w, "denied\n")
			}
		} else {
			logger.Error("Wrong credentails")
		}
	} else {
		logger.Error("No basic authentication was provided")
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	fmt.Fprintf(w, "denied\n")
}

// ----------------------------------------------------------------------------------
// ---------------------------------     /view      ---------------------------------
// ----------------------------------------------------------------------------------

func viewRequestHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		logger.Warning(fmt.Sprintf("Authentication Attempt: %s with password %s", username, password))
		expectedUsernameHash := sha256.Sum256([]byte(BASIC_AUTH_USER))
		expectedPasswordHash := sha256.Sum256([]byte(BASIC_AUTH_PASS))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
			if r.Method == "POST" || r.Method == "HEAD" {
				logger.Info(fmt.Sprintf("view received a %s Request", r.Method))
				fmt.Fprintf(w, "denied\n")
			} else if r.Method == "GET" {
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

				//t, t_err := template.ParseFiles("html/main/results.html")
				lp := filepath.Join("html/main", "results.html")
				fp := filepath.Join("html/main", "css/styles.css")
				t, t_err := template.ParseFiles(lp, fp)
				if t_err != nil {
					logger.Error(t_err.Error())
					fmt.Fprintf(w, "error\n")
					return
				}

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
				q := getQueryForTemplate(pid[0], data)
				t.Execute(w, q)
				logger.Info(fmt.Sprintf("successfully viewed server data - project : %s", pid[0]))
			} else {
				logger.Warning(fmt.Sprintf("view received a %s Request", r.Method))
				fmt.Fprintf(w, "denied\n")
			}
		} else {
			logger.Error("Wrong credentails")
		}
	} else {
		logger.Error("No basic authentication was provided")
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	fmt.Fprintf(w, "denied\n")
}

// ----------------------------------------------------------------------------------
// ---------------------------------       /        ---------------------------------
// ----------------------------------------------------------------------------------

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		logger.Warning(fmt.Sprintf("Authentication Attempt: %s with password %s", username, password))
		expectedUsernameHash := sha256.Sum256([]byte(BASIC_AUTH_USER))
		expectedPasswordHash := sha256.Sum256([]byte(BASIC_AUTH_PASS))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {

			t, _ := template.ParseFiles("html/main/index.html")

			db, db_err := openDB()
			if db_err != nil {
				logger.Error(db_err.Error())
				fmt.Fprintf(w, "error\n")
				return
			}
			defer db.Close()

			tq, query_err_code, query_err := queryDatabaseTables(db)
			if query_err != nil {
				logger.Error(fmt.Sprintf("Server response code \"%d\" : %s", query_err_code, query_err.Error()))
				fmt.Fprintf(w, "error\n")
				return
			}
			t.Execute(w, tq)
		} else {
			logger.Error("Wrong credentails")
		}
	} else {
		logger.Error("No basic authentication was provided")
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "", http.StatusUnauthorized)
}
