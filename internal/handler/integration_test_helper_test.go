package handler_test

import (
	"path/filepath"
	"testing"

	"github.com/yabu1121/expense-api/internal/store"
)

func newExpenseTestStore(t *testing.T) *store.SQLiteStore {
	t.Helper()

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "expenses.db")

	expenseStore, err := store.NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create expense store: %v", err)
	}

	t.Cleanup(func() {
		expenseStore.Close()
	})

	return expenseStore
}

func newCategoryTestStore(t *testing.T) *store.SQLiteStore {
	t.Helper()

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "categories.db")

	categoryStore, err := store.NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create category store: %v", err)
	}

	t.Cleanup(func() {
		categoryStore.Close()
	})

	return categoryStore
}
