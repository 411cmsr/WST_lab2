package database

import "errors"

var (
	ErrPersonNotFound = errors.New("person not found")
	ErrPersonExists   = errors.New("person exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrEmptyQuery     = errors.New("empty query")
	ErrQueryTooLong   = errors.New("query too long")
	ErrEmailExists    = errors.New("email exists")
	
)
