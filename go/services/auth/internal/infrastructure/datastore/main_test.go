package datastore_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/phamquanandpad/training-project/go/services/auth/internal/infrastructure/datastore"
	"github.com/phamquanandpad/training-project/go/services/auth/internal/testutil"
)

var ctxWithReadDB context.Context

// For [table]_reader_test.go, because READ query might not affect to each other so every testcase might use Parallel and use ctxWithReadDB.
// For [table]_writer_test.go, because WRITE query might affect to each other so avoid using Parallel for each test case and generate a new database for each test function.
func TestMain(m *testing.M) {
	ctxWithReadDB = context.Background()

	// Initialize the template database.
	tmplDB, err := testutil.InitTemplateDB()
	if err != nil {
		log.Fatalf("cannot init template db: %s", err)
	}

	gormDB, _, closeDBFunc, err := testutil.InitReadDB()
	if err != nil {
		log.Fatal("cannot init read db", err)
	}

	ctxWithReadDB = datastore.WithAuthDB(ctxWithReadDB, gormDB)
	exitVal := m.Run()

	// Cleanup read db and template db.
	if err := closeDBFunc(); err != nil {
		log.Printf("close read db: %s", err)
	}
	tmplDB.Release()

	os.Exit(exitVal)
}
