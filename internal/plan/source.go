package plan

import (
	"context"

	"github.com/hashicorp/hcl/v2"
)

type Source interface {
	Validate() error
	Snapshot(ctx context.Context, ec hcl.EvalContext) error
}
