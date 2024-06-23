package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spf13/cobra"
	"github.com/spoletum/dg/config"
	"github.com/spoletum/dg/utils/evalcontext"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func main() {

	app := cobra.Command{}

	app.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validates a configuration file",
		RunE:  validate,
		Args:  cobra.ExactArgs(1),
	})

	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Validate validates the configuration file provided as an argument
func validate(cmd *cobra.Command, args []string) error {

	var config config.Configuration

	// Create an evaluation context from the environment
	variables, err := evalcontext.FromEnvironment(os.Environ())
	if err != nil {
		return err
	}

	evalContext := &hcl.EvalContext{
		Variables: variables,
		Functions: map[string]function.Function{
			"upper":  stdlib.UpperFunc,
			"lower":  stdlib.LowerFunc,
			"min":    stdlib.MinFunc,
			"max":    stdlib.MaxFunc,
			"strlen": stdlib.StrlenFunc,
			"substr": stdlib.SubstrFunc,
		},
	}

	filePath := args[0]

	err = hclsimple.DecodeFile(filePath, nil, &config)

	if err != nil {
		return err
	}

	conn, d := config.Postgres[0].Connection.Value(evalContext)
	if d.HasErrors() {
		return fmt.Errorf(d.Error())
	}

	query, d := config.Postgres[0].Tables[0].Query.Value(evalContext)
	if d.HasErrors() {
		return fmt.Errorf(d.Error())
	}

	fmt.Printf("Evaluated connection: %s\n", conn.AsString())
	fmt.Printf("Evaluated query: %s\n", query.AsString())
	return nil
}
