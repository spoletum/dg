package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
)

type PostgresTableBlock struct {
	Name  string         `hcl:"name,label"`
	Query hcl.Expression `hcl:"query"` // The query required to extract the data. It can use variables, e.g. `SELECT * FROM TABLE WHERE FOO = ${bar}`
}

func (b *PostgresTableBlock) Snapshot(ctx context.Context, db *sql.DB, ec *hcl.EvalContext) error {

	// Evaluate the query binding the variables
	query, diag := b.Query.Value(ec)
	if diag.HasErrors() {
		return diag
	}

	// Execute the query
	rows, err := db.QueryContext(ctx, query.AsString())
	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		slog.Info("Row", "table", b.Name)
	}

	// return fmt.Errorf("not implemented")
	return nil
}
