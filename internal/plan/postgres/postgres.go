package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/hashicorp/hcl/v2"
)

// PostgresBlock is a struct that holds a database source
type PostgresBlock struct {
	Id     string               `hcl:"id,label"`
	URL    hcl.Expression       `hcl:"url"` // TODO this should be a static string
	Tables []PostgresTableBlock `hcl:"table,block"`
}

// Validate validates the Postgres struct
func (p PostgresBlock) Validate() error {

	// Id must not be empty
	if p.Id == "" {
		return fmt.Errorf("id must not be empty")
	}

	// Connection must not be empty
	if p.URL == nil {
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

// Snapshot generates a snapshot of the database
func (p PostgresBlock) Snapshot(ctx context.Context, ec *hcl.EvalContext) error {

	// Notify the user we are starting a snapshot
	slog.Info("Snapshotting database", "name", p.Id)

	// Establish a connection to the database
	slog.Debug("Establishing connection to database")
	url, diag := p.URL.Value(ec)
	if diag.HasErrors() {
		return diag
	}
	db, err := sql.Open("postgres", url.AsString())
	if err != nil {
		return err
	}
	defer db.Close()
	slog.Debug("Successfully connected to database")

	// Snapshot each table
	for _, table := range p.Tables {
		if err = table.Snapshot(ctx, db, ec); err != nil {
			return err
		}
	}

	// Notify the user we are done
	return nil
	// return fmt.Errorf("not implemented")
}
