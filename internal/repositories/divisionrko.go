package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/donskova1ex/1cServices/internal/domain"
	"sync"
)

func (r *Repository) GetRkoByDivision(ctx context.Context, from string, to string) ([]*domain.DivisionRko, error) {
	var rkoSlise []*domain.DivisionRko
	var divisionRko *domain.DivisionRko

	var divisionId string
	var result float32
	var quantity int32

	mu := &sync.Mutex{}

	query := `SELECT
	  rfl.DepartamentId AS 'DivisionId'
	 ,SUM(rfl.Amount) AS 'Result'
	 ,COUNT(*) AS 'Quantity'
	FROM RkoForLoans rfl
	WHERE CAST(rfl.Date AS DATE) BETWEEN @from AND @to
	GROUP BY rfl.DepartamentId
	ORDER BY rfl.DepartamentId`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("from", fmt.Sprintf("'%s'", from)), sql.Named("to", fmt.Sprintf("'%s'", to)))
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	for rows.Next() {

		if err := rows.Scan(&divisionId, &quantity, &result); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return rkoSlise, nil
			}
			return nil, fmt.Errorf("scan: %w", err)
		}

		mu.Lock()
		divisionRko.DivisionId = divisionId
		divisionRko.Quantity = quantity
		divisionRko.Result = result
		mu.Unlock()

		rkoSlise = append(rkoSlise, divisionRko)
	}

	return rkoSlise, nil
}
