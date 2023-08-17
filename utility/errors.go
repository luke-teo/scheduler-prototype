package utility

import "errors"

var (
	ErrNotFound = errors.New("requested item was not found")
	ErrConflict = errors.New("item already exists")
)
