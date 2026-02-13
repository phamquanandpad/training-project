package testutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestDB represents a database for testing, it's designed for running with
// a specific testcase.
type TestDB struct {
	dsn    string
	dbName string
	db     *sql.DB

	cleanedUp bool
	mu        sync.RWMutex
}

// newTestDB creates a new TestDB instance.
func newTestDB(
	dsn string,
	dbName string,
) (*TestDB, error) {
	db := &TestDB{
		dsn:    dsn,
		dbName: dbName,
	}

	var err error
	db.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	return db, nil
}

// Open establishes a GORM database connection to the test database.
func (db *TestDB) Open() (*gorm.DB, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.cleanedUp {
		return nil, errors.New("db already cleaned up")
	}

	return gorm.Open(
		mysql.Open(db.dsn),
		&gorm.Config{
			SkipDefaultTransaction: true,
		},
	)
}

// Cleanup drops the test database and closes the underlying database connection.
func (db *TestDB) Cleanup() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.cleanedUp {
		return nil
	}

	db.cleanedUp = true

	if _, err := db.db.ExecContext(
		context.TODO(),
		fmt.Sprintf("DROP DATABASE `%s`;", db.dbName),
	); err != nil {
		return fmt.Errorf("drop db: %w", err)
	}

	if err := db.db.Close(); err != nil {
		return fmt.Errorf("db conn close: %w", err)
	}

	return nil
}

// DBName returns the name of the test database.
func (db *TestDB) DBName() string {
	return db.dbName
}
