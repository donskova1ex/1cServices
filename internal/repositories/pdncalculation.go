package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/donskova1ex/1cServices/internal"
	"github.com/donskova1ex/1cServices/internal/domain"
)

func (r *Repository) GetPDNParameters(ctx context.Context, loanid string) (*domain.CalculationParameters, error) {

	pdnParameters := &domain.CalculationParameters{}
	resultChan := make(chan *domain.CalculationParameters, 3)
	errChan := make(chan error, 3)
	defer close(resultChan)
	defer close(errChan)

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go r.getClientsIncomes(ctx, loanid, pdnParameters, wg, resultChan, errChan)
	go r.getLoanDetails(ctx, loanid, pdnParameters, wg, resultChan, errChan)
	go r.getRegionIncomes(ctx, loanid, pdnParameters, wg, resultChan, errChan)
	wg.Wait()

	select {
	case err := <-errChan:
		return nil, err
	case pdnParameters := <-resultChan:
		return pdnParameters, nil
	}
}

func (r *Repository) getClientsIncomes(
	ctx context.Context, loanid string,
	pdnParameters *domain.CalculationParameters,
	wg *sync.WaitGroup,
	resultChan chan<- *domain.CalculationParameters,
	errChan chan<- error) {

	defer wg.Done()
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
    ORDER BY lapcl2.Date DESC
  )`
	row := r.db.QueryRowContext(ctx, query, sql.Named("id", loanid))
	err := row.Scan(&pdnParameters.Incomes, &pdnParameters.Expenses, &pdnParameters.IncomesTypeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errChan <- fmt.Errorf("client incomes not found for loan Id [%s]", loanid)
			return
		}
		errChan <- fmt.Errorf("failed to get client incomes: %w", internal.ErrClientIncomes)
	}
	resultChan <- pdnParameters
}

func (r *Repository) getRegionIncomes(
	ctx context.Context, loanid string,
	pdnParameters *domain.CalculationParameters,
	wg *sync.WaitGroup,
	resultChan chan<- *domain.CalculationParameters,
	errChan chan<- error) {
	defer wg.Done()
	query := `SELECT rai.Value FROM RegionAvgIncomes rai
  				JOIN Regions region ON rai.RegionId = region.Id
  				JOIN Address regaddress ON LEFT(regaddress.KladrId, 2) = region.KladrCode
  				JOIN Clients c ON c.RegAddressId = regaddress.Id
  				JOIN Loans l ON l.ClientId = c.Id
  				WHERE l.Id = @id 
  				AND rai.Year = YEAR(l.StartDate) 
  				AND rai.Quarter = DATEPART(quarter, l.StartDate)`
	row := r.db.QueryRowContext(ctx, query, sql.Named("id", loanid))
	err := row.Scan(&pdnParameters.AverageRegionIncomes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errChan <- fmt.Errorf("err is [%w]; region incomes not found for loan Id [%s]", internal.ErrNotFound, loanid)
			return
		}
		errChan <- fmt.Errorf("failed to get region incomes: %w", internal.ErrRegionIncomes)
		return
	}

	resultChan <- pdnParameters
}

func (r *Repository) getLoanDetails(
	ctx context.Context, loanid string,
	pdnParameters *domain.CalculationParameters,
	wg *sync.WaitGroup,
	resultChan chan<- *domain.CalculationParameters,
	errChan chan<- error) {
	defer wg.Done()
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
	err := row.Scan(&pdnParameters.LoanId, &pdnParameters.IncomesTypeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errChan <- fmt.Errorf("loan details not found for loan Id [%s]", loanid)
			return
		}
		errChan <- fmt.Errorf("failed to get loan details: %w", internal.ErrLoanDetails)
		return
	}
	resultChan <- pdnParameters
}
