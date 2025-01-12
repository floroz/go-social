package mocks

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func SetupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}
