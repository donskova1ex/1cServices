package domain

type CalculationParameters struct {
	LoanId               string  `json:"LoanId" db:"Loanid"`
	Incomes              float32 `json:"Incomes" db:"Incomes"`
	Expenses             float32 `json:"Expenses" db:"Expenses"`
	IncomesTypeId        string  `json:"IncomesTypeId" db:"IncomesTypeId"`
	AverageRegionIncomes float32 `json:"AverageRegionIncomes" db:"AverageRegionIncomes"`
}
