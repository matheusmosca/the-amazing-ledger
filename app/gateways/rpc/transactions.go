package rpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) CreateTransaction(ctx context.Context, req *proto.CreateTransactionRequest) (*emptypb.Empty, error) {
	defer newrelic.FromContext(ctx).StartSegment("CreateTransaction").End()

	logger := zerolog.Ctx(ctx)
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("handler", "CreateTransaction")
	})

	tid, err := uuid.Parse(req.Id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse transaction id")
		return nil, status.Error(codes.InvalidArgument, "invalid transaction id")
	}

	domainEntries := make([]entities.Entry, len(req.Entries))
	for i, entry := range req.Entries {
		entryID, entryErr := uuid.Parse(entry.Id)
		if entryErr != nil {
			logger.Error().Err(entryErr).Int("index", i).Msg("failed to parse entry id")
			return nil, status.Error(codes.InvalidArgument, "invalid entry id")
		}

		metadata, mErr := entry.Metadata.MarshalJSON()
		if mErr != nil {
			logger.Error().Err(mErr).Int("index", i).Msg("failed to marshal entry metadata")
			return nil, status.Error(codes.InvalidArgument, "invalid entry metadata")
		}

		domainEntry, domainErr := entities.NewEntry(
			entryID,
			vos.OperationType(proto.Operation_value[entry.Operation.String()]),
			entry.Account,
			vos.Version(entry.ExpectedVersion),
			int(entry.Amount),
			metadata,
		)
		if domainErr != nil {
			logger.Error().Err(mErr).Int("index", i).Msg("failed to create entry")
			return nil, status.Error(codes.InvalidArgument, domainErr.Error())
		}

		domainEntries[i] = domainEntry
	}

	competenceDate := time.Unix(req.CompetenceDate.Seconds, 0).UTC()
	if competenceDate.After(time.Now().UTC()) {
		return nil, status.Error(codes.InvalidArgument, "competence date set to the future")
	}

	tx, err := entities.NewTransaction(tid, req.Event, req.Company, competenceDate, domainEntries...)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if err := a.UseCase.CreateTransaction(ctx, tx); err != nil {
		logger.Error().Err(err).Msg("failed to save transaction")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &emptypb.Empty{}, nil
}
