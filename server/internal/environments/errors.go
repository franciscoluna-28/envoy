package environments

import "errors"

var ErrEnvironmentNotFound = errors.New("environment not found")
var ErrInvalidConnection = errors.New("invalid database connection")
var ErrEnvironmentAlreadyExists = errors.New("environment already exists")
