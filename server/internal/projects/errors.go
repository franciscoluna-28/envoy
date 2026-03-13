package projects

import "errors"

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrUnauthorized    = errors.New("unauthorized")
)
