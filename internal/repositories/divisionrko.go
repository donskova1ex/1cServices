package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/donskova1ex/1cServices/internal/domain"
)

func (r *Repository) GetRkoByDivision(ctx context.Context, from string, to string) ([]*domain.DivisionRko, error) {
	var rkoSlise []*domain.DivisionRko

	var divisionId string
	var result float32
	var quantity int32

	query := `SELECT
	  rfl.DepartamentId AS 'DivisionId'
	 ,SUM(rfl.Amount) AS 'Result'
	 ,COUNT(*) AS 'Quantity'
	FROM RkoForLoans rfl
	WHERE CAST(rfl.Date AS DATE) BETWEEN @from AND @to
	GROUP BY rfl.DepartamentId
	ORDER BY rfl.DepartamentId`
	//TODO: ctx time 5s & body len, indexes for math operations, ctx with timeout
	rows, err := r.db.QueryContext(ctx, query, sql.Named("from", from), sql.Named("to", to))
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&divisionId, &result, &quantity); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return rkoSlise, nil
			}
			return nil, fmt.Errorf("scan: %w", err)
		}
		divisionRko := &domain.DivisionRko{}
		divisionRko.Quantity = quantity
		divisionRko.Result = result
		divisionRko.DivisionId = divisionId
		rkoSlise = append(rkoSlise, divisionRko)
	}

	return rkoSlise, nil
}
