package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/stone-co/the-amazing-ledger/app/domain"
	"github.com/stone-co/the-amazing-ledger/app/domain/instrumentators"
	"github.com/stone-co/the-amazing-ledger/app/gateways/db/querybuilder"
)

const (
	collection = "entry"
)

var _ domain.Repository = &LedgerRepository{}

type LedgerRepository struct {
	db *pgxpool.Pool
	pb *instrumentators.LedgerInstrumentator
	qb querybuilder.QueryBuilder
}

func NewLedgerRepository(db *pgxpool.Pool, pb *instrumentators.LedgerInstrumentator) *LedgerRepository {
	qb := querybuilder.New(createTransactionQuery, numArgs)
	qb.Init(numDefaultQueries)

	return &LedgerRepository{
		db: db,
		pb: pb,
		qb: qb,
	}
}
