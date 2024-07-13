package evalcontext

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestFromEnvironment(t *testing.T) {
	// Existing test case
	t.Run("Empty source", func(t *testing.T) {
		source := []string{}
		expected := map[string]cty.Value{}
		result, err := FromEnvironment(source)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	// Additional test cases
	t.Run("Single environment variable", func(t *testing.T) {
		source := []string{"FOO=bar"}
		expected := map[string]cty.Value{
			"FOO": cty.StringVal("bar"),
		}
		result, err := FromEnvironment(source)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Multiple environment variables", func(t *testing.T) {
		source := []string{"FOO=bar", "BAZ=qux", "ABC=123"}
		expected := map[string]cty.Value{
			"FOO": cty.StringVal("bar"),
			"BAZ": cty.StringVal("qux"),
			"ABC": cty.StringVal("123"),
		}
		result, err := FromEnvironment(source)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Invalid environment variable", func(t *testing.T) {
		source := []string{"INVALID"}
		expected := errors.New("could not parse environment variable: INVALID")
		result, err := FromEnvironment(source)
		assert.Error(t, err)
		assert.Equal(t, expected, err)
		assert.Nil(t, result)
	})
}

func TestDefault(t *testing.T) {
	// Existing test case
	t.Run("Default evaluation context", func(t *testing.T) {
		result := Default()
		assert.NotNil(t, result)
		assert.Len(t, result.Variables, 0)
		assert.Len(t, result.Functions, 6)
		assert.Contains(t, result.Functions, "upper")
		assert.Contains(t, result.Functions, "lower")
		assert.Contains(t, result.Functions, "min")
		assert.Contains(t, result.Functions, "max")
		assert.Contains(t, result.Functions, "strlen")
		assert.Contains(t, result.Functions, "substr")
	})
}
