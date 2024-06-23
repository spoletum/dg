package config

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

var ErrDuplicatedTable = errors.New("duplicated table")

// Postgres is a struct that holds a database source
type Postgres struct {
	Id         string         `hcl:"id,label"`
	Connection hcl.Expression `hcl:"connection"`
	Tables     []Table        `hcl:"table,block"`
}

type Table struct {
	Name  string         `hcl:"name,label"`
	Query hcl.Expression `hcl:"query"` // The query required to extract the data. It can use variables, e.g. `SELECT * FROM TABLE WHERE FOO = ${bar}`
}

// Validate validates the Postgres struct
func (p Postgres) Validate() error {
	// There must not be duplicated seen
	seen := make(map[string]any)
	for _, table := range p.Tables {
		if _, ok := seen[table.Name]; ok {
			return fmt.Errorf("duplicated table %s in PostgreSQL database %s", table.Name, p.Id)
		}
		seen[table.Name] = nil
	}
	return nil
}
