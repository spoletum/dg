package config

import (
	"io"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spoletum/dg/utils/evalcontext"
)

func Parse(filename string, r io.Reader) (*Configuration, error) {

	var config Configuration

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
