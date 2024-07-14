package transform

import (
	"context"

	"github.com/hashicorp/hcl/v2"
)

// A transform directive instructs doppelganger to alter the value of this column as per the expression provided.
type TransformBlock struct {
	Name     string          `hcl:"name,label"`
	Function *hcl.Expression `hcl:"function,label"` // The function required to transform the data. Example: 'mask()'
}

func Transform(ctx context.Context, name string, kind string, value any, transformBlock *TransformBlock) error {
	return nil
}
