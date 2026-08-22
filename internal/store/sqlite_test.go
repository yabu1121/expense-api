package store

import (
	"errors"
	"testing"

	"github.com/yabu1121/expense-api/internal/model"
)

func TestCreateExpense(t *testing.T) {
	expenseStore := newTestStore(t)

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
	expenseStore := newTestStore(t)

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
		if !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected model.ErrExpenseNotFuond: %v", err)
		}

		if got != nil {
			t.Fatalf("expected nil expense, got: %+v", got)
		}
	})
}

func TestGetAllExpenses(t *testing.T) {
	expenseStore := newTestStore(t)

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
			t.Fatalf("expected 2 expenses, got %d", len(gotExpenses))
		}

		for i, expected := range expenses {
			if gotExpenses[i].ID != i+1 {
				t.Fatalf(
					"expected expense ID %d, got %d",
					i+1,
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

	t.Run("empty", func(t *testing.T) {
		emptyStore := newTestStore(t)

		got, err := emptyStore.GetAllExpenses()

		if err != nil {
			t.Fatalf("failed to get all expenses: %v", err)
		}

		if len(got) != 0 {
			t.Fatalf("expected 0 expenses, got %d", len(got))
		}
	})
}
func TestUpdateExpense(t *testing.T) {
	expenseStore := newTestStore(t)

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

		update := model.Expense{
			ID:       1,
			Title:    "latte",
			Amount:   550,
			Category: "food",
		}

		updatedExpense, err := expenseStore.UpdateExpense(update)
		if err != nil {
			t.Fatalf("failed to update expense: %v", err)
		}

		got, err := expenseStore.GetExpenseByID(1)
		if err != nil {
			t.Fatalf("failed to get expense: %v", err)
		}

		if got.ID != updatedExpense.ID {
			t.Fatalf(
				"expected expense ID %d, got %d",
				1,
				got.ID,
			)
		}
		if got.Title != updatedExpense.Title {
			t.Fatalf(
				"expected expense title %s, got %s",
				updatedExpense.Title,
				got.Title,
			)
		}
		if got.Amount != updatedExpense.Amount {
			t.Fatalf(
				"expected expense amount %d, got %d",
				updatedExpense.Amount,
				got.Amount,
			)
		}
		if got.Category != updatedExpense.Category {
			t.Fatalf(
				"expected expense category %s, got %s",
				updatedExpense.Category,
				got.Category,
			)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := expenseStore.UpdateExpense(model.Expense{ID: 999})
		if !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected model.ErrExpenseNotFound: %v", err)
		}
	})
}

func TestDeleteExpense(t *testing.T) {
	expenseStore := newTestStore(t)
	_, err := expenseStore.CreateExpense(model.Expense{
		Title:    "coffee",
		Amount:   500,
		Category: "food",
	})
	if err != nil {
		t.Fatalf("failed to create expense: %v", err)
	}
	t.Run("success", func(t *testing.T) {
		if err := expenseStore.DeleteExpense(1); err != nil {
			t.Fatalf("failed to delete expense: %v", err)
		}

		got, err := expenseStore.GetExpenseByID(1)
		if !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected model.ErrExpenseNotFound, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil expense, got: %+v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		if err := expenseStore.DeleteExpense(999); !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected model.ErrExpenseNotFound, got: %v", err)
		}
	})
}

func TestGetExpenseSummary(t *testing.T) {
	tests := []struct {
		name           string
		expenses       []model.Expense
		expectedResult model.ExpenseSummary
	}{
		{
			name:  "success1",
			expectedResult: model.ExpenseSummary{
				Count:       0,
				TotalAmount: 0,
			},
		},
		{
			name:  "success2",
			expenses: []model.Expense{
				{
					ID:       1,
					Title:    "coffee",
					Amount:   500,
					Category: "food",
				},
				{
					ID:       2,
					Title:    "latte",
					Amount:   550,
					Category: "food",
				},
			},
			expectedResult: model.ExpenseSummary{
				Count:       2,
				TotalAmount: 1050,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenseStore := newTestStore(t)
			for _, expense := range tt.expenses {
				_, err := expenseStore.CreateExpense(expense)
				if err != nil {
					t.Fatalf("failed to create expense: %v", err)
				}
			}

			getExpenseSummary, err := expenseStore.GetExpenseSummary()
			if err != nil {
				t.Fatalf("failed to get expense summary: %v", err)
			}

			if getExpenseSummary.Count != tt.expectedResult.Count {
				t.Fatalf(
					"expected expense summary count %d, got %d",
					tt.expectedResult.Count,
					getExpenseSummary.Count,
				)
			}

			if getExpenseSummary.TotalAmount != tt.expectedResult.TotalAmount {
				t.Fatalf(
					"expected expense summary total amount %d, got %d",
					tt.expectedResult.TotalAmount,
					getExpenseSummary.TotalAmount,
				)
			}
		})
	}
}
