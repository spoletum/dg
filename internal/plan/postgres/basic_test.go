package postgres

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// This is a very basic and somehow redundant test. It just checks that we can connect to a Postgres container
// without dependencies like the init scripts. To be used in case of issues with the more complex tests.
func TestConnection(t *testing.T) {

	// Create a Postgres container and wait until it's ready
	ctx := context.Background()
	ctr, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage(imageName),
		readyCondition,
	)
	require.NoError(t, err)
	defer ctr.Terminate(ctx)

	// Open the connection to the database
	url := ctr.MustConnectionString(ctx, "sslmode=disable")
	db, err := sql.Open("postgres", url)
	require.NoError(t, err)
	defer db.Close()

	// Run a "SELECT 1" query and assert the result
	var result int
	err = db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, 1, result)

}
