package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

// ------------------------------------------------------------------------
// ---------------------------------  DB  ---------------------------------
// ------------------------------------------------------------------------

func openDB() (*sql.DB, error) {
	db, db_err := sql.Open("sqlite3", DATABASE_PATH)
	if db_err != nil {
		return nil, fmt.Errorf("Unable to open database : %s", db_err.Error())
	}
	return db, nil
}

// ------------------------------------------------------------------------
// --------------------------------- Misc ---------------------------------
// ------------------------------------------------------------------------

func sortInterfaceSlice(s [][]interface{}, ind int) [][]interface{} {
	ret := [][]interface{}{}
	var max_min int = 0
	for i := 0; i < len(s); i++ {
		var x int
		var min int = 1000000
		var index int
		for a, b := range s {
			switch b[ind].(type) {
			case int64:
				x = int(b[ind].(int64))
			case int:
				x = int(b[ind].(int64))
			case float32:
				x = int(b[ind].(int64))
			case float64:
				x = int(b[ind].(int64))
			default:
				x = int(b[ind].(int64))
			}
			if x < min && x > max_min {
				min = x
				index = a
			}
		}
		max_min = min
		ret = append(ret, s[index])
	}
	return ret
}

// ---------------------------------------------------------------------------
// ---------------------------------   Web   ---------------------------------
// ---------------------------------------------------------------------------

func healthyRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "running")
}

// ---------------------------------------------------------------------------
// --------------------------------- Testing ---------------------------------
// ---------------------------------------------------------------------------

func prettify(expectation string, actual string) string {
	return fmt.Sprintf("\n\t[ WANTED ] %s\n\t[ ACTUAL ] %s\n\n", expectation, actual)
}
