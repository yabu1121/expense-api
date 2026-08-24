package store

import (
	"path/filepath"
	"testing"
)

func newExpenseTestStore(t *testing.T) *SQLiteStore {
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

func newCategoryTestStore(t *testing.T) *SQLiteStore {
	t.Helper()

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "categories.db")

	categoryStore, err := NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create category store: %v", err)
	}

	t.Cleanup(func() {
		categoryStore.Close()
	})

	return categoryStore
}
