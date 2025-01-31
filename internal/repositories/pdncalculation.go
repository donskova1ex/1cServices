package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/donskova1ex/1cServices/internal"
	"github.com/donskova1ex/1cServices/internal/domain"
	"golang.org/x/sync/errgroup"
)

func (r *Repository) GetPDNParameters(ctx context.Context, loanid string) (*domain.CalculationParameters, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	g := new(errgroup.Group)

	pdnParameters := &domain.CalculationParameters{}
	mu := &sync.Mutex{}

	g.Go(func() error {
		return r.getClientsIncomes(ctx, loanid, pdnParameters, mu)
	})

	g.Go(func() error {
		return r.getLoanDetails(ctx, loanid, pdnParameters, mu)
	})

	g.Go(func() error {
		return r.getRegionIncomes(ctx, loanid, pdnParameters, mu)
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return pdnParameters, nil

}

func (r *Repository) getClientsIncomes(
	ctx context.Context, loanid string,
	pdnParameters *domain.CalculationParameters,
	mu *sync.Mutex,
) error {

	query := `SELECT
				lapcl.Income AS 'Incomes',
				lapcl.Expenses AS 'Expenses',
				CASE
					WHEN ci.ClientIncomeType = 0 THEN 'Зарплата'
					WHEN ci.ClientIncomeType = 1 THEN 'Пенсия'
					WHEN ci.ClientIncomeType = 2 THEN 'Выписка'
					ELSE 'Не опрпделено'
				END AS 'IncomesTypeId'
				FROM 
				LoanApplicationPdnCalcLogs lapcl
				LEFT JOIN Loans l ON lapcl.LoanApplicationId = l.LoanApplicationId
				LEFT JOIN ClientIncomes ci ON ci.ClientId = l.ClientId AND ci.CreateDate <= l.CreateDate
				WHERE 
				l.Id = @id
				AND ci.Id = (
					SELECT TOP 1 ci2.Id 
					FROM ClientIncomes ci2 
					WHERE ci2.ClientId = l.ClientId 
					ORDER BY ci2.CreateDate DESC
				) -- Выбираем последний доход клиента
				AND lapcl.id = (
				SELECT TOP 1 lapcl2.Id 
					FROM LoanApplicationPdnCalcLogs lapcl2 
					WHERE lapcl2.LoanApplicationId = l.LoanApplicationId
					ORDER BY lapcl2.Date DESC)`
	row := r.db.QueryRowContext(ctx, query, sql.Named("id", loanid))
	var incomes float32
	var expenses float32
	var incomesTypeId string
	err := row.Scan(&incomes, &expenses, &incomesTypeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("client incomes not found for loan Id [%s]", loanid)

		}
		return fmt.Errorf("failed to get client incomes: %w", internal.ErrClientIncomes)

	}
	mu.Lock()
	pdnParameters.Incomes = incomes
	pdnParameters.Expenses = expenses
	pdnParameters.IncomesTypeId = incomesTypeId
	mu.Unlock()
	return nil
}

func (r *Repository) getRegionIncomes(
	ctx context.Context, loanid string,
	pdnParameters *domain.CalculationParameters,
	mu *sync.Mutex,
) error {
	query := `SELECT rai.Value FROM RegionAvgIncomes rai
					JOIN Regions region ON rai.RegionId = region.Id
					JOIN Address regaddress ON LEFT(regaddress.KladrId, 2) = region.KladrCode
					JOIN Clients c ON c.RegAddressId = regaddress.Id
					JOIN Loans l ON l.ClientId = c.Id
					WHERE l.Id = @id 
					AND rai.Year = YEAR(l.StartDate) 
					AND rai.Quarter = DATEPART(quarter, l.StartDate)`
	row := r.db.QueryRowContext(ctx, query, sql.Named("id", loanid))
	var averageRegionIncomes float32
	err := row.Scan(&averageRegionIncomes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("err is [%w]; region incomes not found for loan Id [%s]", internal.ErrNotFound, loanid)

		}
		return fmt.Errorf("failed to get region incomes: %w", internal.ErrRegionIncomes)

	}
	mu.Lock()
	pdnParameters.AverageRegionIncomes = averageRegionIncomes
	mu.Unlock()
	return nil
}

func (r *Repository) getLoanDetails(
	ctx context.Context, loanid string,
	pdnParameters *domain.CalculationParameters,
	mu *sync.Mutex,
) error {
	query := `SELECT
				  l.Id AS 'LoanId',
				  CASE
					WHEN ci.ClientIncomeType = 0 THEN 'Зарплата'
					WHEN ci.ClientIncomeType = 1 THEN 'Пенсия'
					WHEN ci.ClientIncomeType = 2 THEN 'Выписка'
					ELSE 'Не определено'
				  END AS 'IncomesTypeId'
				FROM 
				  LoanApplicationPdnCalcLogs lapcl
				  LEFT JOIN Loans l ON lapcl.LoanApplicationId = l.LoanApplicationId
				  LEFT JOIN ClientIncomes ci ON ci.ClientId = l.ClientId AND ci.CreateDate <= l.CreateDate
				WHERE 
				  l.Id = @id
				  AND ci.Id = (
					SELECT TOP 1 ci2.Id 
					FROM ClientIncomes ci2 
					WHERE ci2.ClientId = l.ClientId 
					ORDER BY ci2.CreateDate DESC
				  )
				  AND lapcl.id = (
				  SELECT TOP 1 lapcl2.Id 
					FROM LoanApplicationPdnCalcLogs lapcl2 
					WHERE lapcl2.LoanApplicationId = l.LoanApplicationId
					ORDER BY lapcl2.Date DESC
				  )`
	row := r.db.QueryRowContext(ctx, query, sql.Named("id", loanid))
	var loanId string
	var incomesTypeId string
	err := row.Scan(&loanId, &incomesTypeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("loan details not found for loan Id [%s]", loanid)

		}
		return fmt.Errorf("failed to get loan details: %w", internal.ErrLoanDetails)
	}
	mu.Lock()
	pdnParameters.IncomesTypeId = incomesTypeId
	pdnParameters.LoanId = loanId
	mu.Unlock()
	return nil
}
