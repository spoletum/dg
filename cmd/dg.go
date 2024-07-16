package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/spf13/cobra"
	"github.com/spoletum/dg/internal/plan"
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

	var config plan.Plan

	filePath := args[0]

	err := hclsimple.DecodeFile(filePath, nil, &config)
	if err != nil {
		return err
	}

	return nil
}
