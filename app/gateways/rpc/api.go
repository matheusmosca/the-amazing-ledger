package rpc

import (
	"github.com/stone-co/the-amazing-ledger/app/domain"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

var _ proto.LedgerServiceServer = &API{}

type API struct {
	UseCase domain.UseCase
}

func NewAPI(useCase domain.UseCase) *API {
	return &API{
		UseCase: useCase,
	}
}
