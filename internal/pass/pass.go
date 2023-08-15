package pass

import (
	"context"

	"github.com/llumus/x-plugin-wallet-passes/internal/components"
)

type Passbook interface {
	CreatePass(ctx context.Context, request components.CreatePassbookRequestObject) (string, string, error)
}
