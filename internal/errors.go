package internal

import "errors"

var ErrNotFound = errors.New("entity not found")
var ErrGetById = errors.New("could not get entity by id")
var ErrReadRows = errors.New("can not read rows")
