package model

import "errors"

var (
	ErrNotFound = errors.New("entity not found")
	ErrInvalid  = errors.New("invalid input")
	ErrConflict = errors.New("state conflict")
	ErrCanceled = errors.New("operation canceled")
)
