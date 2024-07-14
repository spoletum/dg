package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
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

	// Create a cty object structure using the data types from the resultset
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Create an array of any to scan the values into
	row := make([]any, len(columns))
	rowPtrs := make([]interface{}, len(columns))
	for i := range row {
		rowPtrs[i] = &row[i]
	}
	// Get the list of column names
	columns, err = rows.Columns()
	if err != nil {
		return err
	}

	// Iterate over the rows and build the data model
	for rows.Next() {

		// Scan the row values into the row array
		if err := rows.Scan(rowPtrs...); err != nil {
			return err
		}

		// row := make([]any, len(columns))     // where row values will be scanned into
		values := make(map[string]cty.Value) // where the object model will be built

		// Iterate through columns and wrap the values around a cty Value
		for i, column := range columns {
			slog.Warn("Processing column", "name", column)
			value, err := b.toCtyValue(row[i])
			if err != nil {
				return err
			}
			values[column] = value
		}

		// Assemble the values into a cty object
		obj := cty.ObjectVal(values)

		slog.Info("Row built", "data", obj)
	}

	return nil

}

func (b *PostgresTableBlock) toCtyValue(value any) (cty.Value, error) {
	// Return an inferred type based on the SQL type
	switch valueType := value.(type) {
	case string:
		return gocty.ToCtyValue(value, cty.String)
	case int:
		return gocty.ToCtyValue(value, cty.Number)
	case float64, float32:
		return gocty.ToCtyValue(value, cty.Number)
	case bool:
		return gocty.ToCtyValue(value, cty.Bool)
	case time.Time:
		iso8601 := value.(time.Time).Format(time.RFC3339)
		return gocty.ToCtyValue(iso8601, cty.String)
	default:
		slog.Warn("Unsupported type", "type", valueType)
		return cty.NilVal, fmt.Errorf("unsupported type: %v", valueType)
	}
}
