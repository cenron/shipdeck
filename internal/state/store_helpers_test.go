package state

import (
	"context"
	"strings"
	"testing"
)

func TestExecuteTxInsertsAndReturnsID(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())
	ctx := context.Background()

	if _, err := db.ExecContext(ctx, `CREATE TABLE tx_test (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`); err != nil {
		t.Fatalf("create tx_test table: %v", err)
	}

	id, err := s.executeTx(ctx, `INSERT INTO tx_test (name) VALUES (?)`, "example")
	if err != nil {
		t.Fatalf("executeTx() returned error: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected inserted id > 0, got %d", id)
	}
}

func TestExecuteTxReturnsErrorOnQueryFailure(t *testing.T) {
	t.Parallel()

	db := mustOpenTestDB(t)
	t.Cleanup(func() { _ = db.Close() })

	s := NewStore(db, testLogger())

	_, err := s.executeTx(context.Background(), `INSERT INTO does_not_exist (name) VALUES (?)`, "bad")
	if err == nil {
		t.Fatal("expected executeTx() to return error")
	}
	if !strings.Contains(err.Error(), "inserting snapshot") {
		t.Fatalf("expected insert failure context, got: %v", err)
	}
}

func TestParseDBTimeRejectsInvalidFormat(t *testing.T) {
	t.Parallel()

	if _, err := parseDBTime("not-a-time"); err == nil {
		t.Fatal("expected parseDBTime() to fail for invalid input")
	}
}
