package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Pilladian/go-helper"
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
//	606 : Unable to obtain amount of entries from database
//	607 : Unable to obtain amount of columns from database
func queryDatabaseTables(db *sql.DB) (TableQuery, int, error) {
	tq := TableQuery{}
	tq.Header = append(tq.Header, "Table")
	tq.Header = append(tq.Header, "Entries")
	tq.Header = append(tq.Header, "Columns")

	rows, query_err := db.Query("SELECT name FROM sqlite_schema WHERE type IN ('table','view')	AND name NOT LIKE 'sqlite_%';")
	if query_err != nil {
		return TableQuery{}, 604, fmt.Errorf("Unable to perform query : %s", query_err.Error())
	}
	defer rows.Close()

	tables := []string{}
	for rows.Next() {
		var tmp string
		if read_data_err := rows.Scan(&tmp); read_data_err != nil {
			return TableQuery{}, 605, fmt.Errorf("Unable to obtain table names from database : %s", read_data_err.Error())
		}
		tables = append(tables, tmp)
	}

	entries := []int{}
	for _, table := range tables {
		rows, query_err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM \"%s\"", table))
		if query_err != nil {
			return TableQuery{}, 604, fmt.Errorf("Unable to perform query : %s", query_err.Error())
		}
		defer rows.Close()
		var tmp int
		for rows.Next() {
			if read_data_err := rows.Scan(&tmp); read_data_err != nil {
				return TableQuery{}, 606, fmt.Errorf("Unable to obtain amount of entries from database : %s", read_data_err.Error())
			}
		}
		entries = append(entries, tmp)
	}

	columns := []int{}
	for _, table := range tables {
		rows, query_err := db.Query(fmt.Sprintf("SELECT * FROM \"%s\"", table))
		if query_err != nil {
			return TableQuery{}, 604, fmt.Errorf("Unable to perform query : %s", query_err.Error())
		}
		defer rows.Close()
		tmp, tmp_err := rows.Columns()
		if tmp_err != nil {
			return TableQuery{}, 607, fmt.Errorf("Unable to obtain amount of columns from database : %s", tmp_err.Error())
		}
		columns = append(columns, len(tmp))
	}

	for i := 0; i < len(tables); i++ {
		var tmp TableRowQuery
		tmp.Name = tables[i]
		tmp.AmountEntries = entries[i]
		tmp.AmountColumns = columns[i]
		tq.Rows = append(tq.Rows, tmp)
	}

	return tq, 0, nil
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
