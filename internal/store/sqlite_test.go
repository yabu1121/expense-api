package store

import (
	"cmp"
	"errors"
	"slices"
	"testing"
	"uuid"

	"github.com/yabu1121/expense-api/internal/model"
)

// expense
func TestCreateExpense(t *testing.T) {
	expenseStore := newExpenseTestStore(t)

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
			createdExpense, err := expenseStore.CreateExpense(tt.body)
			if err != nil {
				t.Fatalf("failed to insert expense: %v", err)
			}
			expected := tt.body
			expected.ID = 1

			if *createdExpense != expected {
				t.Fatalf("expected %+v, got %+v", expected, *createdExpense)
			}
		})
	}
}

func TestGetExpenseByID(t *testing.T) {
	expenseStore := newExpenseTestStore(t)

	tests := []struct {
		name             string
		existingExpenses []model.Expense
		id               int
		wantExpense      model.Expense
		wantErr          error
	}{
		{
			name: "success",
			existingExpenses: []model.Expense{
				{
					ID:       1,
					Title:    "coffee",
					Amount:   500,
					Category: "food",
				},
			},
			id: 1,
			wantExpense: model.Expense{
				ID:       1,
				Title:    "coffee",
				Amount:   500,
				Category: "food",
			},
		},
		{
			name:    "not found",
			wantErr: model.ErrExpenseNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, e := range tt.existingExpenses {
				_, err := expenseStore.CreateExpense(e)
				if err != nil {
					t.Fatalf("failed to create expense: %v", err)
				}
			}

			if tt.wantErr == nil {
				got, err := expenseStore.GetExpenseByID(tt.id)
				if err != nil {
					t.Fatalf("failed to get expense by id: %v", err)
				}

				if *got != tt.wantExpense {
					t.Fatalf("expected got == wantExpense. but didn't.")
				}
			} else {
				got, err := expenseStore.GetExpenseByID(tt.id)
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v: %v", tt.wantErr, err)
				}

				if got != nil {
					t.Fatalf("expected nil expense, got: %+v", got)
				}
			}
		})
	}
}

func TestListExpenses(t *testing.T) {
	expenseStore := newExpenseTestStore(t)

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

		gotExpenses, err := expenseStore.ListExpenses(model.ExpenseFilter{})
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
		emptyStore := newExpenseTestStore(t)

		got, err := emptyStore.ListExpenses(model.ExpenseFilter{})

		if err != nil {
			t.Fatalf("failed to get all expenses: %v", err)
		}

		if len(got) != 0 {
			t.Fatalf("expected 0 expenses, got %d", len(got))
		}
	})
}

func TestListExpensesByCategory(t *testing.T) {
	tests := []struct {
		name     string
		expenses []model.Expense
		category string
		want     []model.Expense
	}{
		{
			name: "two_food",
			expenses: []model.Expense{
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
			},
			category: "food",
			want: []model.Expense{
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
		},
		{
			name: "food_and_drink",
			expenses: []model.Expense{
				{
					Title:    "coffee",
					Amount:   500,
					Category: "food",
				},
				{
					Title:    "latte",
					Amount:   550,
					Category: "drink",
				},
			},
			category: "food",
			want: []model.Expense{
				{
					ID:       1,
					Title:    "coffee",
					Amount:   500,
					Category: "food",
				},
			},
		},
		{
			name: "category_not_found",
			expenses: []model.Expense{
				{
					Title:    "coffee",
					Amount:   500,
					Category: "food",
				},
				{
					Title:    "latte",
					Amount:   550,
					Category: "drink",
				},
			},
			category: "travel",
			want:     []model.Expense{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenseStore := newExpenseTestStore(t)
			for _, e := range tt.expenses {
				_, err := expenseStore.CreateExpense(e)
				if err != nil {
					t.Fatalf("failed to create expense: %v", err)
				}
			}

			gotExpenses, err := expenseStore.ListExpenses(model.ExpenseFilter{
				Category: tt.category,
			})
			if err != nil {
				t.Fatalf("failed to list expenses by category: %v", err)
			}

			if !slices.Equal(gotExpenses, tt.want) {
				t.Fatalf("expected %+v, got %+v", tt.want, gotExpenses)
			}
		})
	}
}

func TestUpdateExpense(t *testing.T) {
	expenseStore := newExpenseTestStore(t)

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
	expenseStore := newExpenseTestStore(t)
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
			name: "success1",
			expectedResult: model.ExpenseSummary{
				Count:       0,
				TotalAmount: 0,
			},
		},
		{
			name: "success2",
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
			expenseStore := newExpenseTestStore(t)
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

// category
func TestCreateCategory(t *testing.T) {
	category := model.Category{
		Name: "food",
	}

	tests := []struct {
		name     string
		existing []model.Category
		body     model.Category
		wantErr  error
	}{
		{
			name:    "success",
			body:    category,
			wantErr: nil,
		},
		{
			name: "conflict",
			existing: []model.Category{
				{
					Name: "food",
				},
			},
			body:    category,
			wantErr: model.ErrCategoryAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryStore := newCategoryTestStore(t)
			if tt.existing != nil {
				for _, c := range tt.existing {
					id := uuid.NewV7()
					_, err := categoryStore.db.Exec(`
						insert into categories (id, name) values (?, ?)
					`, id, c.Name)
					if err != nil {
						t.Fatalf("failed to prepare category: %v", err)
					}
				}
			}

			createdCategory, err := categoryStore.CreateCategory(tt.body)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if err != nil {
				return
			}

			if createdCategory.ID == uuid.Nil() {
				t.Fatal("expected category ID to be generated")
			}

			if createdCategory.Name != tt.body.Name {
				t.Fatalf(
					"expected category name %s, got %s",
					tt.body.Name,
					createdCategory.Name,
				)
			}
		})
	}
}

func TestListCategories(t *testing.T) {
	tests := []struct {
		name       string
		categories []model.Category
	}{
		{
			name: "none",
		},
		{
			name: "one category",
			categories: []model.Category{
				{
					Name: "food",
				},
			},
		},
		{
			name: "two category",
			categories: []model.Category{
				{
					Name: "food",
				},
				{
					Name: "drink",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryStore := newCategoryTestStore(t)
			want := make([]model.Category, 0)
			for _, category := range tt.categories {
				createdCategory, err := categoryStore.CreateCategory(category)
				if err != nil {
					t.Fatalf("failed to create category: %v", err)
				}
				want = append(want, *createdCategory)
			}

			categories, err := categoryStore.ListCategories()
			if err != nil {
				t.Fatalf("failed to list category: %v", err)
			}

			slices.SortFunc(want, func(a, b model.Category) int {
				return cmp.Compare(a.Name, b.Name)
			})

			if !slices.Equal(categories, want) {
				t.Fatalf(
					"expected %+v, got %+v",
					want,
					categories,
				)
			}
		})
	}
}

func TestGetCategoryByID(t *testing.T) {

	tests := []struct {
		name         string
		categoryName string
		wantErr      error
	}{
		{
			name:         "success",
			categoryName: "food",
		},
		{
			name:    "not found",
			wantErr: model.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryStore := newCategoryTestStore(t)

			createdCategory, err := categoryStore.CreateCategory(model.Category{
				Name: tt.categoryName,
			})
			if err != nil {
				t.Fatalf("failed to create category: %v", err)
			}

			if tt.wantErr != nil {
				got, err := categoryStore.GetCategoryByID(uuid.New())
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v: %v", tt.wantErr, err)
				}

				if got != nil {
					t.Fatalf("expected nil expense, got: %+v", got)
				}
			} else {
				got, err := categoryStore.GetCategoryByID(createdCategory.ID)
				if err != nil {
					t.Fatalf("failed to get category by id: %v", err)
				}

				if *got != *createdCategory {
					t.Fatalf("got and wantCategory is unmatched")
				}
			}
		})
	}

}
