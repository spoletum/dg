package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustGetHclExpression(t *testing.T) {

	t.Run("Valid HCL expression", func(t *testing.T) {
		// Test case 1: Valid HCL expression
		expr := MustGetHclExpression("hello ${name}")
		assert.NotNil(t, expr, "Expected HCL expression to be non-nil")
	})

	t.Run("Invalid HCL expression", func(t *testing.T) {
		// Test case 2: Invalid HCL expression
		assert.Panics(t, func() {
			MustGetHclExpression("hello ${name")
		})
	})
}
