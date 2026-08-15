package store

import (
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "expenses.db")

	expenseStore, err := NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create expense store: %v", err)
	}

	t.Cleanup(func() {
		expenseStore.Close()
	})

	return expenseStore
}
