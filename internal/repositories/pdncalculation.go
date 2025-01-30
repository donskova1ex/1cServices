package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/donskova1ex/1cServices/internal"
	"github.com/donskova1ex/1cServices/internal/domain"
)

func (r *Repository) GetPDNParameters(ctx context.Context, loanid string) (*domain.CalculationParameters, error) {
	pdnParameters := &domain.CalculationParameters{}
	query := `SELECT
 l.Id AS 'LoanId',
 lapcl.Income AS 'Incomes',
 lapcl.Expenses AS 'Expenses',
 CASE
   WHEN ci.ClientIncomeType = 0 THEN 'Зарплата'
   WHEN ci.ClientIncomeType = 1 THEN 'Пенсия'
   WHEN ci.ClientIncomeType = 2 THEN 'Выписка'
   ELSE 'Не опрпделено'
 END AS 'IncomesTypeId',
 rai.Value AS 'AverageRegionIncomes'
FROM
 LoanApplicationPdnCalcLogs lapcl
 LEFT JOIN Loans l ON lapcl.LoanApplicationId = l.LoanApplicationId
 LEFT JOIN ClientIncomes ci ON ci.ClientId = l.ClientId AND ci.CreateDate <= l.CreateDate
 LEFT JOIN Clients c ON c.Id = l.ClientId -- Получаем данные клиента
 LEFT JOIN Address regaddress ON c.RegAddressId = regaddress.Id -- Получаем регистрационный адрес клиента
 LEFT JOIN Regions region ON LEFT(regaddress.KladrId, 2) = region.KladrCode -- Сравниваем первые две цифры KladrId
 LEFT JOIN RegionAvgIncomes rai ON rai.RegionId = region.Id
   AND rai.Quarter = DATEPART(QUARTER, l.CreateDate) -- Квартал создания заявки
   AND rai.Year = DATEPART(YEAR, l.CreateDate) -- Год создания заявки
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
	err := row.Scan(
		&pdnParameters.LoanId,
		&pdnParameters.Incomes,
		&pdnParameters.Expenses,
		&pdnParameters.IncomesTypeId,
		&pdnParameters.AverageRegionIncomes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w with Id [%s]", internal.ErrNotFound, loanid)
	}
	return pdnParameters, nil
}
