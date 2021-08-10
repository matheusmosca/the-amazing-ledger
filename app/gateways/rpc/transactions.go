package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

func (a *API) CreateTransaction(ctx context.Context, req *proto.CreateTransactionRequest) (*emptypb.Empty, error) {
	tid, err := uuid.Parse(req.Id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to parse transaction id")
		return nil, status.Error(codes.InvalidArgument, "invalid transaction id")
	}

	if req.CompetenceDate == nil {
		return nil, status.Error(codes.InvalidArgument, "competence_date must have a value")
	} else if !req.CompetenceDate.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "competence_date must be valid")
	}

	domainEntries := make([]entities.Entry, len(req.Entries))
	for i, entry := range req.Entries {
		entryID, entryErr := uuid.Parse(entry.Id)
		if entryErr != nil {
			zerolog.Ctx(ctx).Error().Err(entryErr).Int("index", i).Msg("failed to parse entry id")
			return nil, status.Error(codes.InvalidArgument, "invalid entry id")
		}

		metadata, mErr := entry.Metadata.MarshalJSON()
		if mErr != nil {
			zerolog.Ctx(ctx).Error().Err(mErr).Int("index", i).Msg("failed to marshal entry metadata")
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
			zerolog.Ctx(ctx).Error().Err(domainErr).Int("index", i).Msg("failed to create entry")
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
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create transaction")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := a.UseCase.CreateTransaction(ctx, tx); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to save transaction")
		switch {
		case errors.Is(err, app.ErrInvalidVersion):
			return nil, status.Error(codes.InvalidArgument, "invalid account version")
		case errors.Is(err, app.ErrIdempotencyKeyViolation):
			return nil, status.Error(codes.InvalidArgument, "invalid idempotency key")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &emptypb.Empty{}, nil
}
