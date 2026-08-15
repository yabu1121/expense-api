package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yabu1121/expense-api/internal/model"
)

func TestWriteFile(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte(`hello world`)

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("got %v, expected %v", string(got), string(content))
	}
}

func TestCreateExpense(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "expenses.db")

	expenseStore, err := NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	t.Cleanup(func() {
		expenseStore.Close()
	})

	expense := model.Expense{
		Title:    "coffee",
		Amount:   500,
		Category: "food",
	}

	tests := []struct {
		name  string
		body  model.Expense
		store SQLiteStore
	}{
		{
			name:  "success",
			store: SQLiteStore{db: expenseStore.db},
			body:  expense,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdExpense, err := expenseStore.CreateExpense(expense)
			if err != nil {
				t.Fatalf("failed to insert expense: %v", err)
			}

			if createdExpense.ID != 1 {
				t.Fatalf(
					"expected expense ID %d, got %d",
					1,
					createdExpense.ID,
				)
			}
			if createdExpense.Title != expense.Title {
				t.Fatalf(
					"expected expense title %s, got %s",
					expense.Title,
					createdExpense.Title,
				)
			}
			if createdExpense.Amount != expense.Amount {
				t.Fatalf(
					"expected expense amount %d, got %d",
					expense.Amount,
					createdExpense.Amount,
				)
			}
			if createdExpense.Category != expense.Category {
				t.Fatalf(
					"expected expense category %s, got %s",
					expense.Category,
					createdExpense.Category,
				)
			}
		})
	}
}

func TestGetExpenseByID(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "expenses.db")

	expenseStore, err := NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	t.Cleanup(func() {
		expenseStore.Close()
	})

	t.Run("success", func(t *testing.T) {
		expense := model.Expense{
			Title:    "coffee",
			Amount:   500,
			Category: "food",
		}

		_, err := expenseStore.CreateExpense(expense)
		if err != nil {
			t.Fatalf("failed to create expense: %v", err)
		}

		got, err := expenseStore.GetExpenseByID(1)
		if err != nil {
			t.Fatalf("failed to get expense by id: %v", err)
		}

		if got.ID != 1 {
			t.Fatalf(
				"expected expense ID %d, got %d",
				1,
				got.ID,
			)
		}

		if got.Title != "coffee" {
			t.Fatalf(
				"expected expense title %s, got %s",
				"coffee",
				got.Title,
			)
		}

		if got.Amount != 500 {
			t.Fatalf(
				"expected expense amount %d, got %d",
				500,
				got.Amount,
			)
		}

		if got.Category != "food" {
			t.Fatalf(
				"expected expense category %s, got %s",
				"food",
				got.Category,
			)
		}
	})

	t.Run("not found", func(t *testing.T) {
		got, err := expenseStore.GetExpenseByID(999)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows: %v", err)
		}

		if got != nil {
			t.Fatalf("expected nil expense, got: %+v", got)
		}
	})
}

func TestGetAllExpenses(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "expenses.db")
	expenseStore, err := NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	t.Cleanup(func() {
		expenseStore.Close()
	})

	expenses := []model.Expense{
		{
			Title:    "coffee",
			Amount:   500,
			Category: "food",
		},
		{
			Title:    "latte",
			Amount:   550,
			Category: "food",
		},
	}

	t.Run("Success", func(t *testing.T) {
		for _, e := range expenses {
			_, err := expenseStore.CreateExpense(e)
			if err != nil {
				t.Fatalf("failed to create expense: %v", err)
			}
		}

		gotExpenses, err := expenseStore.GetAllExpenses()
		if err != nil {
			t.Fatalf("failed to get all expenses: %v", err)
		}

		if len(gotExpenses) != 2 {
			t.Fatalf("expected 2 expenses length, got %d", len(gotExpenses))
		}

		for i, expected := range expenses {
			if gotExpenses[i].ID != expected.ID {
				t.Fatalf(
					"expected expense ID %d, got %d",
					expected.ID,
					gotExpenses[i].ID,
				)
			}
			if gotExpenses[i].Title != expected.Title {
				t.Fatalf(
					"expected expense title %s, got %s",
					expected.Title,
					gotExpenses[i].Title,
				)
			}
			if gotExpenses[i].Amount != expected.Amount {
				t.Fatalf(
					"expected expense amount %d, got %d",
					expected.Amount,
					gotExpenses[i].Amount,
				)
			}
			if gotExpenses[i].Category != expected.Category {
				t.Fatalf(
					"expected expense category %s, got %s",
					expected.Category,
					gotExpenses[i].Category,
				)
			}

		}
	})
}
