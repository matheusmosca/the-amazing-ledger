package usecases

import (
	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
)

var _ domain.UseCase = &LedgerUseCase{}

type LedgerUseCase struct {
	instrumentator *instrumentators.LedgerInstrumentator
	repository     domain.Repository
}

func NewLedgerUseCase(repository domain.Repository, instrumentator *instrumentators.LedgerInstrumentator) *LedgerUseCase {
	return &LedgerUseCase{
		repository:     repository,
		instrumentator: instrumentator,
	}
}
