package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/spoletum/dg/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresTableBlock_Snapshot(t *testing.T) {

	// Create a Postgres container and wait until it's ready
	ctx := context.Background()
	ctr, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage(imageName),
		postgres.WithInitScripts("postgres_test.sql"),
		readyCondition,
	)
	require.NoError(t, err)
	defer ctr.Terminate(ctx)

	t.Run("Run correct test case", func(t *testing.T) {

		// Create a new test database
		db, err := sql.Open("postgres", ctr.MustConnectionString(ctx, "sslmode=disable"))
		require.NoError(t, err)
		defer db.Close()

		// Create a new table block
		block := &PostgresTableBlock{
			Name:  "test_table",
			Query: testutils.MustGetHclExpression("SELECT ID, NAME, SURNAME, MIDDLE_NAME, BIRTH_DATE, DELETED FROM CLIENTS"),
		}

		// Execute the snapshot
		err = block.Snapshot(context.Background(), db, &hcl.EvalContext{})
		assert.NoError(t, err)
	})

}
