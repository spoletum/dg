package postgres

import (
	"context"
	"log/slog"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	_ "github.com/lib/pq"
	"github.com/spoletum/dg/internal/evalcontext"
	"github.com/spoletum/dg/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	imageName    = "postgres:16.3"
	readyMessage = "database system is ready to accept connections"
)

var readyCondition testcontainers.CustomizeRequestOption = testcontainers.WithWaitStrategy(wait.ForLog(readyMessage).WithOccurrence(2))

func TestPostgresBlock_Validate(t *testing.T) {
	// Test case 1: Valid PostgresBlock
	t.Run("Valid PostgresBlock", func(t *testing.T) {
		conn, d := hclsyntax.ParseExpression([]byte("foobar"), "", hcl.InitialPos)
		assert.False(t, d.HasErrors(), "Expected no error for valid HCL expression")
		p := PostgresBlock{
			Id:  "postgres1",
			URL: conn,
			Tables: []PostgresTableBlock{
				{
					Name:  "table1",
					Query: nil,
				},
				{
					Name:  "table2",
					Query: nil,
				},
			},
		}
		err := p.Validate()
		assert.NoError(t, err, "Expected no error for valid PostgresBlock")
	})

	// Test case 2: Invalid PostgresBlock with duplicated table names
	t.Run("Invalid PostgresBlock with duplicated table names", func(t *testing.T) {
		conn, d := hclsyntax.ParseExpression([]byte("foobar"), "", hcl.InitialPos)
		assert.False(t, d.HasErrors(), "Expected no error for valid HCL expression")
		p := PostgresBlock{
			Id:  "postgres2",
			URL: conn,
			Tables: []PostgresTableBlock{
				{
					Name:  "table1",
					Query: nil,
				},
				{
					Name:  "table1",
					Query: nil,
				},
			},
		}
		err := p.Validate()
		assert.EqualError(t, err, "duplicated table table1 in PostgreSQL database postgres2", "Expected error for duplicated table names")
	})
}

func TestPostgresBlock_Snapshot(t *testing.T) {

	// Start a PostgreSQL container
	ctx := context.Background()
	ctr, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage(imageName),
		postgres.WithInitScripts("postgres_test.sql"),
		readyCondition,
	)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer ctr.Terminate(ctx)
	slog.Info("PostgreSQL container started", "url", ctr.MustConnectionString(ctx, "sslmode=disable"))

	t.Run("Valid PostgresBlock", func(t *testing.T) {

		// Create the evaluation context
		ec := evalcontext.Default()

		// Create a PostgresBlock
		p := PostgresBlock{
			Id:  "postgres1",
			URL: testutils.MustGetHclExpression(ctr.MustConnectionString(ctx, "sslmode=disable")),
			Tables: []PostgresTableBlock{
				{
					Name:  "table1",
					Query: testutils.MustGetHclExpression("SELECT * FROM TRANSACTIONS T"),
				},
			},
		}

		// Run the snapshot
		err = p.Snapshot(ctx, ec)
		assert.NoError(t, err, "Expected no error for valid PostgresBlock")
	})
}
