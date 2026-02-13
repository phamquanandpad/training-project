package datastore_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/infrastructure/datastore"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/testutil"
)

var ctxWithReadDB context.Context

// For [table]_reader_test.go, because READ query might not affect to each other so every testcase might use Parallel and use ctxWithReadDB.
// example: /go/services/owner/internal/infrastructure/datastore/share_reader_test.go@TestListShareByOwnerID
// For [table]_writer_test.go, because WRITE query might affect to each other so avoid using Parallel for each test case and generate a new database for each test function, with each testcase, we expect to rollback data after it finished
// example: /go/services/owner/internal/infrastructure/datastore/share_writer_test.go@TestCreateShare
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

	ctxWithReadDB = datastore.WithTodoDB(ctxWithReadDB, gormDB)
	exitVal := m.Run()

	// Cleanupread db and template db.
	closeDBFunc()
	tmplDB.Release()

	os.Exit(exitVal)
}
