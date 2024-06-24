package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	imageName    = "postgres:16.3"
	readyMessage = "database system is ready to accept connections"
)

func TestEndToEnd(t *testing.T) {

	ctx := context.Background()

	// Wait for the container to be ready
	ctr, err := postgres.RunContainer(ctx, testcontainers.WithImage(imageName), testcontainers.WithWaitStrategy(wait.ForLog(readyMessage)))
	assert.NoError(t, err)

	t.Run("TestEndToEnd", func(t *testing.T) {
		assert.NotEmpty(t, ctr.MustConnectionString(ctx))
		fmt.Println(ctr.MustConnectionString(ctx))
	})
}
