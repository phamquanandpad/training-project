package testutil

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/uuid"
)

// TemplateDBConfig represents a config for creating TemplateDB.
type TemplateDBConfig struct {
	DBHost       string
	DBPort       int
	DBUser       string
	DBPass       string
	DBNamePrefix string

	FixturesDir    string
	SchemaFilename string
}

// TemplateDB represents a template database for testing.
type TemplateDB struct {
	cfg *TemplateDBConfig

	dsn         string
	dbName      string       // name of db.
	db          *sql.DB      // connection to db.
	tableNames  []string     // list of table names of db.
	foreignKeys []ForeignKey // list of foreign keys of db.

	initialized bool
	mu          sync.RWMutex

	nRef int // number of references to this template db.
}

func NewTemplateDB(cfg *TemplateDBConfig) *TemplateDB {
	return &TemplateDB{cfg: cfg}
}

// Init initializes a new template database if it hasn't been initialized yet.
func (db *TemplateDB) Init() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.initialized {
		return nil
	}

	db.dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?parseTime=true&multiStatements=true",
		db.cfg.DBUser,
		db.cfg.DBPass,
		db.cfg.DBHost,
		db.cfg.DBPort,
	)

	var err error
	db.db, err = sql.Open("mysql", db.dsn)
	if err != nil {
		return fmt.Errorf("open db conn: %w", err)
	}

	id := uuid.New()
	db.dbName = fmt.Sprintf("%s_template_%x", db.cfg.DBNamePrefix, id[:])

	if err := db.createDB(); err != nil {
		return fmt.Errorf("create db: %w", err)
	}

	if err := db.applySchema(); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	if err := db.loadFixtures(); err != nil {
		return fmt.Errorf("load fixtures: %w", err)
	}

	if err := db.loadTableNames(); err != nil {
		return fmt.Errorf("load table names: %w", err)
	}

	if err := db.loadForeignKeys(); err != nil {
		return fmt.Errorf("load foreign keys: %w", err)
	}

	db.initialized = true

	return nil
}

// createDB creates a new MySQL database with UTF8MB4 encoding and switches to
// use it.
func (db *TemplateDB) createDB() error {
	createDBStmt := fmt.Sprintf(
		"CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;",
		db.dbName,
	)
	if _, err := db.db.ExecContext(
		context.TODO(),
		createDBStmt,
	); err != nil {
		return fmt.Errorf("create db: %w", err)
	}

	if _, err := db.db.ExecContext(
		context.TODO(),
		fmt.Sprintf("USE `%s`;", db.dbName),
	); err != nil {
		return fmt.Errorf("switch to db, dbName=%s: %w", db.dbName, err)
	}

	return nil
}

// applySchema applies the database schema from sql file.
func (db *TemplateDB) applySchema() error {
	schemaStmt, err := os.ReadFile(db.cfg.SchemaFilename)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	if _, err := db.db.ExecContext(
		context.TODO(),
		string(schemaStmt),
	); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	return nil
}

// loadFixtures loads test data into the database from fixture files.
func (db *TemplateDB) loadFixtures() error {
	// Load fixtures.
	fixtures, err := testfixtures.New(
		testfixtures.Database(db.db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(db.cfg.FixturesDir),
	)
	if err != nil {
		return fmt.Errorf("create fixtures loader: %w", err)
	}

	// Load fixtures.
	if err := fixtures.Load(); err != nil {
		return fmt.Errorf("load fixtures: %w", err)
	}

	return nil
}

// loadTableNames loads all table names in database.
func (db *TemplateDB) loadTableNames() error {
	rows, err := db.db.QueryContext(
		context.TODO(),
		fmt.Sprintf(`SHOW TABLES FROM %s`, db.dbName),
	)
	if err != nil {
		return fmt.Errorf("show tables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}
		db.tableNames = append(db.tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("sql rows: %w", err)
	}

	return nil
}

type ForeignKey struct {
	TableName            string
	ConstraintName       string
	ColumnName           string
	ReferencedTableName  string
	ReferencedColumnName string
	UpdateRule           string
	DeleteRule           string
}

// loadForeignKeys gets all foreign keys from database.
func (db *TemplateDB) loadForeignKeys() error {
	for _, tableName := range db.tableNames {
		foreignKeys, err := db.getForeignKeysFromTable(tableName)
		if err != nil {
			return fmt.Errorf("get foreign key, table = %s: %w", tableName, err)
		}
		db.foreignKeys = append(db.foreignKeys, foreignKeys...)
	}

	return nil
}

// getForeignKeysFromTable gets all foreign keys of a table from database.
func (db *TemplateDB) getForeignKeysFromTable(tableName string) ([]ForeignKey, error) {
	query := `
        SELECT
            kcu.constraint_name,
            kcu.column_name,
            kcu.referenced_table_name,
            kcu.referenced_column_name,
            rc.update_rule,
            rc.delete_rule
        FROM
            information_schema.TABLE_CONSTRAINTS tc
        JOIN
            information_schema.KEY_COLUMN_USAGE kcu ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema AND tc.table_name = kcu.table_name
        JOIN
            information_schema.REFERENTIAL_CONSTRAINTS rc ON tc.constraint_name = rc.constraint_name AND tc.table_schema = rc.constraint_schema
        WHERE
            tc.constraint_type = 'FOREIGN KEY'
            AND tc.table_schema = ?
            AND tc.table_name = ?
    `
	rows, err := db.db.QueryContext(
		context.TODO(),
		query,
		db.dbName,
		tableName,
	)
	if err != nil {
		return nil, fmt.Errorf("get foreign keys: %w", err)
	}

	defer rows.Close()

	var foreignKeys []ForeignKey
	for rows.Next() {
		fk := ForeignKey{
			TableName: tableName,
		}
		if err := rows.Scan(
			&fk.ConstraintName,
			&fk.ColumnName,
			&fk.ReferencedTableName,
			&fk.ReferencedColumnName,
			&fk.UpdateRule,
			&fk.DeleteRule,
		); err != nil {
			return nil, err
		}
		foreignKeys = append(foreignKeys, fk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sql rows: %w", err)
	}

	return foreignKeys, nil
}

// NewTestDB creates a new test database by cloning the schema and data from the
// template database. It returns a TestDB instance connected to the newly
// created database.
func (db *TemplateDB) NewTestDB() (*TestDB, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if !db.initialized {
		return nil, errors.New("db is not initialized")
	}

	id := uuid.New()
	dbName := fmt.Sprintf("%s_%x", db.cfg.DBNamePrefix, id[:])

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		"CREATE DATABASE %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;",
		dbName,
	))

	for _, tableName := range db.tableNames {
		sb.WriteString(fmt.Sprintf(
			"\nCREATE TABLE %s.%s LIKE %s.%s;",
			dbName,
			tableName,
			db.dbName,
			tableName,
		))
		sb.WriteString(fmt.Sprintf(
			"\nINSERT INTO %s.%s SELECT * FROM %s.%s;",
			dbName,
			tableName,
			db.dbName,
			tableName,
		))
	}

	for _, fk := range db.foreignKeys {
		sb.WriteString(fmt.Sprintf(
			"\nALTER TABLE %s.%s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s.%s (%s)",
			dbName,
			fk.TableName,
			fk.ConstraintName,
			fk.ColumnName,
			dbName,
			fk.ReferencedTableName,
			fk.ReferencedColumnName,
		))

		if fk.DeleteRule != "" {
			sb.WriteString(fmt.Sprintf(" ON DELETE %s", fk.DeleteRule))
		}

		if fk.UpdateRule != "" {
			sb.WriteString(fmt.Sprintf(" ON UPDATE %s", fk.UpdateRule))
		}

		sb.WriteString(";")
	}

	if _, err := db.db.ExecContext(
		context.TODO(),
		sb.String(),
	); err != nil {
		return nil, fmt.Errorf("clone db: %w", err)
	}

	return newTestDB(
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s",
			db.cfg.DBUser,
			db.cfg.DBPass,
			db.cfg.DBHost,
			db.cfg.DBPort,
			dbName,
			url.QueryEscape("Asia/Tokyo"),
		),
		dbName,
	)
}

// Cleanup drops the template database and closes the underlying database
// connection.
func (db *TemplateDB) cleanup() error {
	if !db.initialized {
		return nil
	}

	db.initialized = false
	db.tableNames = nil
	db.foreignKeys = nil

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

// Acquire acquires a reference to the template database.
func (db *TemplateDB) Acquire() (*TemplateDB, error) {
	log.Printf("Acquire template db, nRef = %d", db.nRef)
	db.mu.Lock()
	defer db.mu.Unlock()

	if !db.initialized {
		return nil, errors.New("db is not initialized")
	}

	db.nRef++

	return db, nil
}

// Release releases a reference to the template database.
func (db *TemplateDB) Release() error {
	log.Printf("Release template db, nRef = %d", db.nRef)
	db.mu.Lock()
	defer db.mu.Unlock()

	if !db.initialized {
		return errors.New("db is not initialized")
	}

	if db.nRef == 0 {
		return nil
	}

	db.nRef--
	if db.nRef > 0 {
		return nil
	}

	if err := db.cleanup(); err != nil {
		return fmt.Errorf("cleanup db: %w", err)
	}

	return nil
}
