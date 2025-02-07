package processors

import (
	"context"
	"fmt"
	"github.com/donskova1ex/1cServices/internal/domain"
	"log/slog"
)

type DivisionRkoRepository interface {
	GetRkoByDivision(ctx context.Context, from string, to string) ([]*domain.DivisionRko, error)
}
type DivisionRkoLogger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}
type divisionRko struct {
	dvisionRkoRepository DivisionRkoRepository
	log                  DivisionRkoLogger
}

func NewDivisionRko(divisionRkoRepository DivisionRkoRepository, log DivisionRkoLogger) *divisionRko {
	return &divisionRko{
		dvisionRkoRepository: divisionRkoRepository,
		log:                  log,
	}
}

func (d *divisionRko) DivisionRko(ctx context.Context, from string, to string) ([]*domain.DivisionRko, error) {
	divisionRko, err := d.dvisionRkoRepository.GetRkoByDivision(ctx, from, to)
	if err != nil {
		d.log.Error(
			"error getting division rko",
			slog.String("err", err.Error()),
			slog.String("from", from),
			slog.String("to", to),
		)
		return nil, fmt.Errorf("error getting division rko: %w", err)
	}
	return divisionRko, nil
}
