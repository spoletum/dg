package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// PostgresBlock is a struct that holds a database source
type PostgresBlock struct {
	Id         string         `hcl:"id,label"`
	Connection hcl.Expression `hcl:"connection"`
	Tables     []TableBlock   `hcl:"table,block"`
}

type TableBlock struct {
	Name  string         `hcl:"name,label"`
	Query hcl.Expression `hcl:"query"` // The query required to extract the data. It can use variables, e.g. `SELECT * FROM TABLE WHERE FOO = ${bar}`
}

// Validate validates the Postgres struct
func (p PostgresBlock) Validate() error {

	// Id must not be empty
	if p.Id == "" {
		return fmt.Errorf("id must not be empty")
	}

	// Connection must not be empty
	if p.Connection == nil {
		return fmt.Errorf("connection must not be empty")
	}

	// There must not be duplicated table names
	seen := make(map[string]any)
	for _, table := range p.Tables {
		if _, ok := seen[table.Name]; ok {
			return fmt.Errorf("duplicated table %s in PostgreSQL database %s", table.Name, p.Id)
		}
		seen[table.Name] = nil
	}
	return nil
}
