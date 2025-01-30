package repositories

import (
	"context"
	"fmt"

	"github.com/donskova1ex/1cServices/internal/domain"
)

func (r *Repository) GetPDNParameters(ctx context.Context, loanid string) (*domain.CalculationParameters, error) {
	pdnParameters := &domain.CalculationParameters{}
	query := "SELECT * FROM PDNParameters WHERE LoanID = $1"
	row, err := r.db.QueryContext(ctx, query, loanid)
	if err != nil {
		return nil, fmt.Errorf("sql query error: %w", err)
	}
	err = row.Scan(&pdnParameters)
	if err != nil {
		return nil, fmt.Errorf("row reading error: %w", err)

	}
	return pdnParameters, nil
}
