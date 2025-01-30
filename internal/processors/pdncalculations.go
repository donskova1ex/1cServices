package processors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/donskova1ex/1cServices/internal/domain"
)

type PdnCalculationRepository interface {
	GetPDNParameters(ctx context.Context, loanid string) (*domain.CalculationParameters, error)
}

type PdnCalculationLogger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

type pdnCalculation struct {
	pdnCalculationRepository PdnCalculationRepository
	log                      PdnCalculationLogger
}

func NewPDNParametres(pdnCalculationRepository PdnCalculationRepository, log PdnCalculationLogger) *pdnCalculation {
	return &pdnCalculation{pdnCalculationRepository: pdnCalculationRepository, log: log}
}

func (p *pdnCalculation) PDNCalculationByLoanId(ctx context.Context, loanid string) (*domain.CalculationParameters, error) {

	pdnCalculation, err := p.pdnCalculationRepository.GetPDNParameters(ctx, loanid)
	if err != nil {
		p.log.Error(
			"Error getting PDN Calculation Parameters",
			slog.String("err", err.Error()),
			slog.String("loanid", loanid),
		)
		return nil, fmt.Errorf("error getting PDN Calculation Parameters: %w", err)
	}
	return pdnCalculation, nil
}
