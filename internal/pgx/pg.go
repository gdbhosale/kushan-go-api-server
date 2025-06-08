package pgx

import (
	"go-api/internal"

	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Seed SQL File
func Seed(db *sqlx.DB, sqlFile string) error {
	internal.Debug("Pgx::Seed", sqlFile)

	// Read SQL File
	file, err := os.ReadFile(sqlFile)
	if err != nil {
		internal.Error("Pgx::Seed", "Read File Error for "+sqlFile, err)
		return err
	}

	// Start SQL Transaction
	tx, err := db.Begin()
	if err != nil {
		internal.Error("Pgx::Seed", "DB Begin Error for "+sqlFile, err)
		return err
	}

	// Rollback on failure
	defer func() {
		// internal.Error("Pgx::Seed", sqlFile, "Rollback")
		tx.Rollback()
	}()

	// Run SQL Queries
	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(q); err != nil {
			internal.Error("Pgx::Seed", "Error in Execution for "+sqlFile, err)
			return err
		}
	}

	// Execute
	return tx.Commit()
}

// FormatLimitOffset returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}
