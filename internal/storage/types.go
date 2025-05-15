package storage

import "errors"

var (
	// marshal json data for pebble
	ErrJsonMarshal   = errors.New("json marshal error")
	ErrJsonUnMarshal = errors.New("json unmarshal error")
	ErrCollNotFound  = errors.New("collection not found")
)

// Sql query condition define
type Condition struct {
	Field    string
	Operator string // "=", ">", "<"
	Value    interface{}
}

type Query struct {
	Conditions []Condition
	OrderBy    string
	Limit      int
	Offset     int
}
