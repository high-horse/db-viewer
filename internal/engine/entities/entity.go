package entities

import "time"


type ConnectionConfig struct {
	ID   string
	Name string
	Type string

	Host string
	Port int

	User     string
	Password string
	Database string

	SSL bool
}

type QueryResult struct {
	Columns []string
	Rows [][]any
	Duration time.Duration
}


type Table struct {
	Name string
	Schema string
}

type Column struct {
	Name string
	Type string
}