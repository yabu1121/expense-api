package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yabu1121/expense-api/internal/model"
)

type fakeExpenseStore struct {
	expenses []model.Expense
	err error
}

func (f *fakeExpenseStore) GetAllExpenses() ([]model.Expense, error) {
	return f.expenses, f.err
}

func (f *fakeExpenseStore) GetExpenseByID(id int) (*model.Expense, error) {
	for i := range f.expenses {
		if f.expenses[i].ID == id {
			return &f.expenses[i], nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *fakeExpenseStore) CreateExpense(expense model.Expense) (*model.Expense, error) {
	return &f.expenses[0], f.err
}

func (f *fakeExpenseStore) UpdateExpense(expense model.Expense) (*model.Expense, error) {
	return &f.expenses[0], f.err
}

func (f *fakeExpenseStore) DeleteExpense(id int) error {
	return f.err
}


func TestGetAllExpenses(t *testing.T) {
	type testTemplete struct{
		name string
		store *fakeExpenseStore
		expectedStatus int
	}

	tests := []testTemplete{
		{
			name: "success",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID: 1,
						Title: "coffee",
						Amount: 500,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "store error",
			store: &fakeExpenseStore{
				err: errors.New("store error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewExpenseHandler(tt.store)

			req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}
		})
	}
}

func TestGetExpenseByID(t *testing.T) {
	fakeStore := &fakeExpenseStore{
		expenses: []model.Expense{
			{
				ID: 1,
				Title: "coffee",
				Amount: 500,
				Category: "food",
			},
		},
	}

	handler := NewExpenseHandler(fakeStore)

	req := httptest.NewRequest(http.MethodGet, "/expenses/1", nil)
	req.SetPathValue("id", "1")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected to %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	var expense model.Expense

	if err := json.NewDecoder(recorder.Body).Decode(&expense); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if expense.ID != 1 {
		t.Fatalf(
			"expected id %d, got %d",
			1,
			expense.ID,
		)
	}
	if expense.Title != "coffee" {
		t.Fatalf(
			"expected title %v, got %v",
			"coffee",
			expense.Title,
		)
	}
	if expense.Amount != 500 {
		t.Fatalf(
			"expected amount %d, got %d",
			500,
			expense.Amount,
		)
	}
	if expense.Category != "food" {
		t.Fatalf(
			"expected category %v, got %v",
			"food",
			expense.Category,
		)
	}
}
