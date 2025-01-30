package internal

import "errors"

var ErrNotFound = errors.New("entity not found")
var ErrDBConnection = errors.New("database connection error")
var ErrDBPing = errors.New("database ping error")
