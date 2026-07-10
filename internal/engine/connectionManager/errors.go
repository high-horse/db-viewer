package manager

import "errors"

var (
	ErrConnectionExists = errors.New("connection already exists")
	ErrConnectionNotFound =errors.New("connection not found")
	ErrInvalidConnection = errors.New("invalid connection")
)