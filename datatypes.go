package main

type ProjectQuery struct {
	ProjectID string
	Header    []string
	Rows      [][]interface{}
}

type TableQuery struct {
	Header []string
	Rows   []TableRowQuery
}

type TableRowQuery struct {
	Name          string
	AmountEntries int
	AmountColumns int
}
