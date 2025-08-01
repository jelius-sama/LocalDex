package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var Conn *sql.DB

// Initialize all tables and triggers
func InitializeSchema(db *sql.DB) error {
	// Set PRAGMA settings to improve concurrency and reliability
	pragmas := []string{
		`PRAGMA journal_mode = WAL;`,  // concurrent read/write
		`PRAGMA foreign_keys = ON;`,   // enforce FKs
		`PRAGMA busy_timeout = 3000;`, // wait 3s if db is locked
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to set pragma %q: %w", pragma, err)
		}
	}
	return nil
}

func WithRetryWrite(execFunc func() error) error {
	const maxAttempts = 5
	const retryDelay = 200 * time.Millisecond

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := execFunc()
		if err == nil {
			return nil
		}

		// Only retry if it's a SQLITE_BUSY error
		if !strings.Contains(err.Error(), "database is locked") {
			return err
		}

		time.Sleep(retryDelay)
	}

	return errors.New("write failed after retries")
}
