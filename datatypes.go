package main

type Query struct {
	ProjectID string
	Header    []string
	Rows      [][]interface{}
}
