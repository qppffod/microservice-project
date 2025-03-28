package client

import (
	"context"

	"github.com/qppffod/microservice-project/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
