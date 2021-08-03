package domain

import (
	"context"
)

type Instrumentator interface {
	MonitorSegment(ctx context.Context) Segment
	MonitorDataSegment(ctx context.Context, collection, operation, query string) Segment
}

type Segment interface {
	End()
}
