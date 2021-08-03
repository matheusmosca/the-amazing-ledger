package instrumentators

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/stone-co/the-amazing-ledger/app/domain"
)

var _ domain.Instrumentator = &LedgerInstrumentator{}

type LedgerInstrumentator struct {
	newrelic *newrelic.Application
}

func NewLedgerInstrumentator(nr *newrelic.Application) *LedgerInstrumentator {
	return &LedgerInstrumentator{
		newrelic: nr,
	}
}

func (lp LedgerInstrumentator) MonitorSegment(ctx context.Context) domain.Segment {
	txn := newrelic.FromContext(ctx)
	seg := &newrelic.Segment{}
	seg.StartTime = txn.StartSegmentNow()
	return seg
}

func (lp LedgerInstrumentator) MonitorDataSegment(ctx context.Context, collection, operation, query string) domain.Segment {
	txn := newrelic.FromContext(ctx)
	seg := &newrelic.DatastoreSegment{
		Product:            newrelic.DatastorePostgres,
		Collection:         collection,
		Operation:          operation,
		ParameterizedQuery: query,
	}
	seg.StartTime = txn.StartSegmentNow()
	return seg
}
