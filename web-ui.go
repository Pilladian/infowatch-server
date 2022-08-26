package main

import (
	"sort"
)

func getQueryForTemplate(pid string, data_json map[int64]map[string]interface{}) ProjectQuery {
	q := ProjectQuery{}
	q.ProjectID = pid
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
			if i != "ID" {
				tmp = append(tmp, b[i])
			}
		}
		m2 = append(m2, tmp)
	}
	q.Rows = sortInterfaceSlice(m2, 0)
	return q
}
