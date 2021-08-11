package vos

import (
	"time"

	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/app/pagination"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

type AccountEntryRequest struct {
	Account   Account
	StartDate time.Time
	EndDate   time.Time
	Filter    AccountEntryFilter
	Page      pagination.Page
}

type AccountEntryFilter struct {
	Companies []string
	Events    []int32
	Operation OperationType
}

func NewEntryFilter(filter *proto.ListAccountEntriesRequest_Filter) AccountEntryFilter {
	if filter == nil {
		return AccountEntryFilter{}
	}

	return AccountEntryFilter{
		Companies: filter.Companies,
		Events:    filter.Events,
		Operation: OperationType(proto.Operation_value[filter.Operation.String()]),
	}
}

type AccountEntryResponse struct {
	Entries  []AccountEntry
	NextPage pagination.Cursor
}

type AccountEntry struct {
	ID             uuid.UUID
	Version        Version
	Operation      OperationType
	Amount         int
	Event          int
	CompetenceDate time.Time
	Metadata       map[string]interface{}
}
