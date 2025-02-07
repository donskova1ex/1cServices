package domain

type DivisionRko struct {
	DivisionId string  `json:"DivisionId" db:"DivisionId"`
	Result     float32 `json:"Result" db:"Result"`
	Quantity   int32   `json:"Quantity" db:"Quantity"`
}
