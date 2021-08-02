package instrumentators

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain"
)

var _ domain.Instrumentator = &LedgerInstrumentator{}

type LedgerInstrumentator struct {
	logger   *logrus.Logger
	newrelic *newrelic.Application
}

func NewLedgerInstrumentator(l *logrus.Logger, nr *newrelic.Application) *LedgerInstrumentator {
	return &LedgerInstrumentator{
		logger:   l,
		newrelic: nr,
	}
}

func (lp LedgerInstrumentator) Log(ctx context.Context, value string) {
	lp.logger.Infof(value)
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

func (lp LedgerInstrumentator) Logger() *logrus.Logger {
	return lp.logger
}
