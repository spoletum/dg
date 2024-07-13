package plan

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	_ "github.com/lib/pq"
	"github.com/spoletum/dg/internal/evalcontext"
	"github.com/spoletum/dg/internal/plan/postgres"
)

type Plan struct {
	Postgres []postgres.PostgresBlock `hcl:"postgres,block"`
}

// Validate validates the ConfigurationBlock struct
func (c Plan) Validate() error {

	// Validate each Postgres block
	for _, p := range c.Postgres {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Generates a snapshot using the configuration of the plan
func (c Plan) Snapshot(ctx context.Context, ec *hcl.EvalContext) error {
	slog.Info("Validating schema")
	if err := c.Validate(); err != nil {
		return err
	}
	slog.Debug("Snapshotting PostgreSQL databases")
	for _, p := range c.Postgres {
		if err := p.Snapshot(ctx, ec); err != nil {
			return err
		}
	}
	slog.Debug("Successfully snapshotted PostgreSQL databases")
	return errors.New("not implemented")
}

// Parse a plan from a reader and generate the data structure
func Parse(filename string, r io.Reader) (*Plan, error) {

	var config Plan

	// Create an evaluation context from the environment
	variables, err := evalcontext.FromEnvironment(os.Environ())
	if err != nil {
		return nil, err
	}

	// Create an evaluation context with the standard functions
	evalContext := evalcontext.Default()
	evalContext.Variables = variables

	// Read the source code
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Decode the source code and store the result in the config variable
	err = hclsimple.Decode(filename, src, evalContext, &config)
	if err != nil {
		return nil, err
	}

	// Return the configuration
	return &config, nil
}
