package testutil

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"gorm.io/gorm"
)

var getTemplateDB = sync.OnceValues(func() (*TemplateDB, error) {
	_, currentFilename, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller error")
	}

	schemaFilename := filepath.Join(
		filepath.Dir(currentFilename),
		"../../database/test/sqls/import/create_tables.sql",
	)

	fixturesDir := path.Join(
		path.Dir(currentFilename),
		"../../testdata/todo_fixtures",
	)

	env := LoadEnv()
	db := NewTemplateDB(
		&TemplateDBConfig{
			DBHost:       env.DBHost,
			DBPort:       env.DBPort,
			DBUser:       env.DBUser,
			DBPass:       env.DBPass,
			DBNamePrefix: env.DBName,

			FixturesDir:    fixturesDir,
			SchemaFilename: schemaFilename,
		},
	)

	return db, nil
})

func InitDB(t *testing.T) (*gorm.DB, string) {
	tmplDB, err := getTemplateDB()
	if err != nil {
		t.Fatalf("get template db: %s", err)
	}

	if err := tmplDB.Init(); err != nil {
		t.Fatalf("init template db: %s", err)
	}

	db, err := tmplDB.NewTestDB()
	if err != nil {
		t.Fatalf("new test db: %s", err)
	}

	gormDB, err := db.Open()
	if err != nil {
		t.Fatalf("new gorm db: %s", err)
	}

	t.Cleanup(func() {
		sqlDB, err := gormDB.DB()
		if err != nil {
			t.Fatalf("get gorm underlying db: %s", err)
		}

		if err := sqlDB.Close(); err != nil {
			t.Fatalf("close gorm underlying db: %s", err)
		}

		if err := db.Cleanup(); err != nil {
			t.Fatalf("clean up db: %s", err)
		}
	})

	return gormDB, db.DBName()
}

func InitReadDB() (gormDB *gorm.DB, dbName string, dbCloseFunc func() error, err error) {
	tmplDB, err := getTemplateDB()
	if err != nil {
		return nil, "", nil, fmt.Errorf("get template db: %w", err)
	}

	if err := tmplDB.Init(); err != nil {
		return nil, "", nil, fmt.Errorf("init template db: %w", err)
	}

	db, err := tmplDB.NewTestDB()
	if err != nil {
		return nil, "", nil, fmt.Errorf("new test db: %w", err)
	}

	gormDB, err = db.Open()
	if err != nil {
		return nil, "", nil, fmt.Errorf("new gorm db: %w", err)
	}

	dbCloseFunc = func() error {
		sqlDB, err := gormDB.DB()
		if err != nil {
			return fmt.Errorf("get gorm underlying db: %w", err)
		}

		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("close gorm underlying db: %w", err)
		}

		if err := db.Cleanup(); err != nil {
			return fmt.Errorf("clean up db: %w", err)
		}

		return nil
	}

	return gormDB, db.DBName(), dbCloseFunc, nil
}

// InitTemplateDB initializes the template database and returns it.
func InitTemplateDB() (*TemplateDB, error) {
	tmplDB, err := getTemplateDB()
	if err != nil {
		return nil, fmt.Errorf("get template db: %w", err)
	}

	if err := tmplDB.Init(); err != nil {
		return nil, fmt.Errorf("init template db: %w", err)
	}

	return tmplDB.Acquire()
}
