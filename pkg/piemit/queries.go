package piemit

import (
	"context"
)

type Querier interface {
	Create(ctx context.Context, e event) error
}
