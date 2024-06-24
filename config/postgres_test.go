package config

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
)

func TestPostgresBlock_Validate(t *testing.T) {
	// Test case 1: Valid PostgresBlock
	t.Run("Valid PostgresBlock", func(t *testing.T) {
		conn, d := hclsyntax.ParseExpression([]byte("foobar"), "", hcl.InitialPos)
		assert.False(t, d.HasErrors(), "Expected no error for valid HCL expression")
		p := PostgresBlock{
			Id:         "postgres1",
			Connection: conn,
			Tables: []TableBlock{
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
			Id:         "postgres2",
			Connection: conn,
			Tables: []TableBlock{
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
