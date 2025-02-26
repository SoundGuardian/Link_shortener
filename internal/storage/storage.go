package storage

import "errors"

var (
	ErrURLNotFound = errors.New("urls not found")
	ErrURLExist    = errors.New("url exists")
)
