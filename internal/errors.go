package internal

import "errors"

var ErrNotFound = errors.New("entity not found")
var ErrDBConnection = errors.New("database connection error")
var ErrDBPing = errors.New("database ping error")
var ErrClientIncomes = errors.New("client incomes error")
var ErrRegionIncomes = errors.New("region incomes error")
var ErrLoanDetails = errors.New("loan detail error")
