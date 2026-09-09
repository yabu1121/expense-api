package store

import (
	"errors"
	"fmt"
	"testing"
	"uuid"

	"github.com/yabu1121/expense-api/internal/model"
)

// expense
func TestCreateExpense(t *testing.T) {
	store := newExpenseTestStore(t)

	createdCategory, err := store.CreateCategory(model.CategoryRequest{
		Name: "food",
	})

	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	successExpense := model.ExpenseRequest{
		Title:      "coffee",
		Amount:     500,
		CategoryID: createdCategory.ID,
	}

	tests := []struct {
		name string
		body model.ExpenseRequest
	}{
		{
			name: "success",
			body: successExpense,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdExpense, err := store.CreateExpense(tt.body)
			if err != nil {
				t.Fatalf("failed to create expense: %v", err)
			}

			if createdExpense.ID == uuid.Nil() {
				t.Fatalf("expected %v, got %v", uuid.Nil(), createdExpense.ID)
			}
			if createdExpense.Title != tt.body.Title {
				t.Fatalf("expected %v, got %v", tt.body.Title, createdExpense.Title)
			}
			if createdExpense.Amount != tt.body.Amount {
				t.Fatalf("expected %v, got %v", tt.body.Amount, createdExpense.Amount)
			}
			if createdExpense.CategoryID != tt.body.CategoryID {
				t.Fatalf("expected %v, got %v", tt.body.CategoryID, createdExpense.CategoryID)
			}
			if createdExpense.CreatedAt.IsZero() {
				t.Fatalf("expected non-zero CreatedAt, got zero")
			}
		})
	}
}

func TestGetExpenseByID(t *testing.T) {
	store := newExpenseTestStore(t)

	createdCategory, err := store.CreateCategory(model.CategoryRequest{
		Name: "food",
	})

	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	createdExpense, err := store.CreateExpense(model.ExpenseRequest{
		Title:      "coffee",
		Amount:     500,
		CategoryID: createdCategory.ID,
	})

	t.Run("success", func(t *testing.T) {
		gotExpense, err := store.GetExpenseByID(createdExpense.ID)
		if err != nil {
			t.Fatalf("failed to get expense: %v", err)
		}

		fmt.Print(gotExpense)

		if gotExpense.ID == uuid.Nil() {
			t.Fatal("uuid must not nil")
		}
		if gotExpense.Title != createdExpense.Title {
			t.Fatalf("expected %v, got %v", createdExpense.Title, gotExpense.Title)
		}
		if gotExpense.Amount != createdExpense.Amount {
			t.Fatalf("expected %d, got %d", createdExpense.Amount, gotExpense.Amount)
		}
		if gotExpense.CategoryID != createdExpense.CategoryID {
			t.Fatalf("expected %v, got %v", createdExpense.CategoryID, gotExpense.CategoryID)
		}
		if gotExpense.CreatedAt.IsZero() {
			t.Fatalf("expected non-zero CreatedAt, got zero")
		}
	})
}

func TestListExpenses(t *testing.T) {
	store := newExpenseTestStore(t)

	createdCategory, err := store.CreateCategory(model.CategoryRequest{
		Name: "food",
	})
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	expenses := []model.ExpenseRequest{
		{
			Title:      "coffee",
			Amount:     500,
			CategoryID: createdCategory.ID,
		},
		{
			Title:      "latte",
			Amount:     550,
			CategoryID: createdCategory.ID,
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, e := range expenses {
			_, err := store.CreateExpense(e)
			if err != nil {
				t.Fatalf("failed to create expense: %v", err)
			}
		}

		gotExpenses, err := store.ListExpenses()
		if err != nil {
			t.Fatalf("failed to get all expenses: %v", err)
		}

		if len(gotExpenses) != 2 {
			t.Fatalf("expected 2 expenses, got %d", len(gotExpenses))
		}
	})

	t.Run("empty", func(t *testing.T) {
		emptyStore := newExpenseTestStore(t)

		got, err := emptyStore.ListExpenses()

		if err != nil {
			t.Fatalf("failed to get all expenses: %v", err)
		}

		if len(got) != 0 {
			t.Fatalf("expected 0 expenses, got %d", len(got))
		}
	})
}

// func TestListExpensesByCategory(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expenses []model.Expense
// 		category string
// 		want     []model.Expense
// 	}{
// 		{
// 			name: "two_food",
// 			expenses: []model.Expense{
// 				{
// 					Title:    "coffee",
// 					Amount:   500,
// 					CategoryID: "food",
// 				},
// 				{
// 					Title:    "latte",
// 					Amount:   550,
// 					CategoryID: "food",
// 				},
// 			},
// 			category: "food",
// 			want: []model.Expense{
// 				{
// 					ID:       1,
// 					Title:    "coffee",
// 					Amount:   500,
// 					CategoryID: "food",
// 				},
// 				{
// 					ID:       2,
// 					Title:    "latte",
// 					Amount:   550,
// 					CategoryID: "food",
// 				},
// 			},
// 		},
// 		{
// 			name: "food_and_drink",
// 			expenses: []model.Expense{
// 				{
// 					Title:    "coffee",
// 					Amount:   500,
// 					CategoryID: "food",
// 				},
// 				{
// 					Title:    "latte",
// 					Amount:   550,
// 					CategoryID: "drink",
// 				},
// 			},
// 			category: "food",
// 			want: []model.Expense{
// 				{
// 					ID:       1,
// 					Title:    "coffee",
// 					Amount:   500,
// 					CategoryID: "food",
// 				},
// 			},
// 		},
// 		{
// 			name: "category_not_found",
// 			expenses: []model.Expense{
// 				{
// 					Title:    "coffee",
// 					Amount:   500,
// 					CategoryID: "food",
// 				},
// 				{
// 					Title:    "latte",
// 					Amount:   550,
// 					CategoryID: "drink",
// 				},
// 			},
// 			category: "travel",
// 			want:     []model.Expense{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			expenseStore := newExpenseTestStore(t)
// 			for _, e := range tt.expenses {
// 				_, err := expenseStore.CreateExpense(e)
// 				if err != nil {
// 					t.Fatalf("failed to create expense: %v", err)
// 				}
// 			}

// 			gotExpenses, err := expenseStore.ListExpenses(model.ExpenseFilter{
// 				Category: tt.category,
// 			})
// 			if err != nil {
// 				t.Fatalf("failed to list expenses by category: %v", err)
// 			}

// 			if !slices.Equal(gotExpenses, tt.want) {
// 				t.Fatalf("expected %+v, got %+v", tt.want, gotExpenses)
// 			}
// 		})
// 	}
// }

func TestUpdateExpense(t *testing.T) {
	store := newExpenseTestStore(t)

	createdCategory, err := store.CreateCategory(model.CategoryRequest{
		Name: "food",
	})
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	createdExpense, err := store.CreateExpense(model.ExpenseRequest{
		Title:      "coffee",
		Amount:     500,
		CategoryID: createdCategory.ID,
	})
	if err != nil {
		t.Fatalf("failed to create expense: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		update := model.ExpenseRequest{
			Title:      "latte",
			Amount:     550,
			CategoryID: createdCategory.ID,
		}

		_, err := store.UpdateExpense(createdExpense.ID, update)
		if err != nil {
			t.Fatalf("failed to update expense: %v", err)
		}

		got, err := store.GetExpenseByID(createdExpense.ID)
		if err != nil {
			t.Fatalf("failed to get expense: %v", err)
		}

		if got.ID != createdExpense.ID {
			t.Fatalf("expected expense ID %d, got %d", createdExpense.ID, got.ID)
		}
		if got.Title != update.Title {
			t.Fatalf("expected expense title %s, got %s", update.Title, got.Title)
		}
		if got.Amount != update.Amount {
			t.Fatalf("expected expense amount %d, got %d", update.Amount, got.Amount)
		}
		if got.CategoryID != update.CategoryID {
			t.Fatalf("expected expense category %s, got %s", update.CategoryID, got.CategoryID)
		}
		if got.CreatedAt != createdExpense.CreatedAt {
			t.Fatalf("expected expense created at %v, got %v", createdExpense.CreatedAt, got.CreatedAt)
		}
	})

	t.Run("not found", func(t *testing.T) {
		notFoundID := uuid.NewV7()
		_, err := store.UpdateExpense(notFoundID, model.ExpenseRequest{
			Title:      "latte",
			Amount:     550,
			CategoryID: createdCategory.ID,
		})
		if !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected model.ErrExpenseNotFound: %v", err)
		}
	})
}

func TestDeleteExpense(t *testing.T) {
	store := newExpenseTestStore(t)

	createdCategory, err := store.CreateCategory(model.CategoryRequest{
		Name: "food",
	})
	if err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	createdExpense, err := store.CreateExpense(model.ExpenseRequest{
		Title:      "coffee",
		Amount:     500,
		CategoryID: createdCategory.ID,
	})
	if err != nil {
		t.Fatalf("failed to create expense: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := store.DeleteExpense(createdExpense.ID); err != nil {
			t.Fatalf("failed to delete expense: %v", err)
		}

		got, err := store.GetExpenseByID(createdExpense.ID)
		if !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected model.ErrExpenseNotFound, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil expense, got: %+v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		notFoundID := uuid.NewV7()
		if err := store.DeleteExpense(notFoundID); !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected model.ErrExpenseNotFound, got: %v", err)
		}
	})
}

// func TestGetExpenseSummary(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		expenses       []model.Expense
// 		expectedResult model.ExpenseSummary
// 	}{
// 		{
// 			name: "success1",
// 			expectedResult: model.ExpenseSummary{
// 				Count:       0,
// 				TotalAmount: 0,
// 			},
// 		},
// 		{
// 			name: "success2",
// 			expenses: []model.Expense{
// 				{
// 					ID:       1,
// 					Title:    "coffee",
// 					Amount:   500,
// 					CategoryID: "food",
// 				},
// 				{
// 					ID:       2,
// 					Title:    "latte",
// 					Amount:   550,
// 					CategoryID: "food",
// 				},
// 			},
// 			expectedResult: model.ExpenseSummary{
// 				Count:       2,
// 				TotalAmount: 1050,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			expenseStore := newExpenseTestStore(t)
// 			for _, expense := range tt.expenses {
// 				_, err := expenseStore.CreateExpense(expense)
// 				if err != nil {
// 					t.Fatalf("failed to create expense: %v", err)
// 				}
// 			}

// 			getExpenseSummary, err := expenseStore.GetExpenseSummary()
// 			if err != nil {
// 				t.Fatalf("failed to get expense summary: %v", err)
// 			}

// 			if getExpenseSummary.Count != tt.expectedResult.Count {
// 				t.Fatalf(
// 					"expected expense summary count %d, got %d",
// 					tt.expectedResult.Count,
// 					getExpenseSummary.Count,
// 				)
// 			}

// 			if getExpenseSummary.TotalAmount != tt.expectedResult.TotalAmount {
// 				t.Fatalf(
// 					"expected expense summary total amount %d, got %d",
// 					tt.expectedResult.TotalAmount,
// 					getExpenseSummary.TotalAmount,
// 				)
// 			}
// 		})
// 	}
// }

// // category
// func TestCreateCategory(t *testing.T) {
// 	category := model.Category{
// 		Name: "food",
// 	}

// 	tests := []struct {
// 		name     string
// 		existing []model.Category
// 		body     model.Category
// 		wantErr  error
// 	}{
// 		{
// 			name:    "success",
// 			body:    category,
// 			wantErr: nil,
// 		},
// 		{
// 			name: "conflict",
// 			existing: []model.Category{
// 				{
// 					Name: "food",
// 				},
// 			},
// 			body:    category,
// 			wantErr: model.ErrCategoryAlreadyExists,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			categoryStore := newCategoryTestStore(t)
// 			if tt.existing != nil {
// 				for _, c := range tt.existing {
// 					id := uuid.NewV7()
// 					_, err := categoryStore.db.Exec(`
// 						insert into categories (id, name) values (?, ?)
// 					`, id, c.Name)
// 					if err != nil {
// 						t.Fatalf("failed to prepare category: %v", err)
// 					}
// 				}
// 			}

// 			createdCategory, err := categoryStore.CreateCategory(tt.body)
// 			if !errors.Is(err, tt.wantErr) {
// 				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
// 			}

// 			if err != nil {
// 				return
// 			}

// 			if createdCategory.ID == uuid.Nil() {
// 				t.Fatal("expected category ID to be generated")
// 			}

// 			if createdCategory.Name != tt.body.Name {
// 				t.Fatalf(
// 					"expected category name %s, got %s",
// 					tt.body.Name,
// 					createdCategory.Name,
// 				)
// 			}
// 		})
// 	}
// }

// func TestListCategories(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		categories []model.Category
// 	}{
// 		{
// 			name: "none",
// 		},
// 		{
// 			name: "one category",
// 			categories: []model.Category{
// 				{
// 					Name: "food",
// 				},
// 			},
// 		},
// 		{
// 			name: "two category",
// 			categories: []model.Category{
// 				{
// 					Name: "food",
// 				},
// 				{
// 					Name: "drink",
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			categoryStore := newCategoryTestStore(t)
// 			want := make([]model.Category, 0)
// 			for _, category := range tt.categories {
// 				createdCategory, err := categoryStore.CreateCategory(category)
// 				if err != nil {
// 					t.Fatalf("failed to create category: %v", err)
// 				}
// 				want = append(want, *createdCategory)
// 			}

// 			categories, err := categoryStore.ListCategories()
// 			if err != nil {
// 				t.Fatalf("failed to list category: %v", err)
// 			}

// 			slices.SortFunc(want, func(a, b model.Category) int {
// 				return cmp.Compare(a.Name, b.Name)
// 			})

// 			if !slices.Equal(categories, want) {
// 				t.Fatalf(
// 					"expected %+v, got %+v",
// 					want,
// 					categories,
// 				)
// 			}
// 		})
// 	}
// }

// func TestGetCategoryByID(t *testing.T) {

// 	tests := []struct {
// 		name         string
// 		categoryName string
// 		wantErr      error
// 	}{
// 		{
// 			name:         "success",
// 			categoryName: "food",
// 		},
// 		{
// 			name:    "not found",
// 			wantErr: model.ErrCategoryNotFound,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			categoryStore := newCategoryTestStore(t)

// 			createdCategory, err := categoryStore.CreateCategory(model.Category{
// 				Name: tt.categoryName,
// 			})
// 			if err != nil {
// 				t.Fatalf("failed to create category: %v", err)
// 			}

// 			if tt.wantErr != nil {
// 				got, err := categoryStore.GetCategoryByID(uuid.New())
// 				if !errors.Is(err, tt.wantErr) {
// 					t.Fatalf("expected %v: %v", tt.wantErr, err)
// 				}

// 				if got != nil {
// 					t.Fatalf("expected nil expense, got: %+v", got)
// 				}
// 			} else {
// 				got, err := categoryStore.GetCategoryByID(createdCategory.ID)
// 				if err != nil {
// 					t.Fatalf("failed to get category by id: %v", err)
// 				}

// 				if *got != *createdCategory {
// 					t.Fatalf("got and wantCategory is unmatched")
// 				}
// 			}
// 		})
// 	}

// }
