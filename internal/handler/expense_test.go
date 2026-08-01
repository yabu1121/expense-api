package handler

import (
	"bytes"
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
	if f.err != nil {
		return nil, f.err
	}
	expense.ID = 1
	return &expense, nil
}

func (f *fakeExpenseStore) UpdateExpense(expense model.Expense) (*model.Expense, error) {
	return &f.expenses[0], f.err
}

func (f *fakeExpenseStore) DeleteExpense(id int) error {
	return f.err
}

func TestGetAllExpenses(t *testing.T) {
	tests := []struct{
			name string
			store *fakeExpenseStore
			expectedStatus int
		}{
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

// func TestCreateExpense(t *testing.T) {
// 	fakeStore := &fakeExpenseStore{}
// 	handler := NewExpenseHandler(fakeStore)

// 	expense := model.Expense{
// 		Title: "coffee",
// 		Amount: 500,
// 		Category: "food",
// 	}

// 	body, err := json.Marshal(expense)
// 	if err != nil {
// 		t.Fatalf("failed to marshal expense: %v", err)
// 	}

// 	reader := bytes.NewReader(body)

// 	req := httptest.NewRequest(
// 		http.MethodPost,
// 		"/expenses",
// 		reader,
// 	)
// 	recorder := httptest.NewRecorder()

// 	handler.ServeHTTP(recorder, req)

// 	if recorder.Code != http.StatusCreated {
// 		t.Fatalf(
// 			"expected code %d, got %d",
// 			http.StatusCreated,
// 			recorder.Code,
// 		)
// 	}

// 	var createdExpense model.Expense

// 	if err := json.NewDecoder(recorder.Body).Decode(&createdExpense); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if createdExpense.ID != 1 {
// 		t.Fatalf(
// 			"expected id %d, got %d",
// 			1,
// 			createdExpense.ID,
// 		)
// 	}
// 	if createdExpense.Title != "coffee" {
// 		t.Fatalf(
// 			"expected id %q, got %q",
// 			"coffee",
// 			createdExpense.Title,
// 		)
// 	}
// 	if createdExpense.Amount != 500 {
// 		t.Fatalf(
// 			"expected id %d, got %d",
// 			500,
// 			createdExpense.Amount,
// 		)
// 	}
// 	if createdExpense.Category != "food" {
// 		t.Fatalf(
// 			"expected id %q, got %q",
// 			"food",
// 			createdExpense.Category,
// 		)
// 	}
// }

// func TestCreateExpenseInvalidJSON(t *testing.T) {
// 	fakeStore := &fakeExpenseStore{}
// 	handler := NewExpenseHandler(fakeStore)

// 	invalidJSON := []byte(`{"title": "coffee"`)

// 	reader := bytes.NewReader(invalidJSON)

// 	req := httptest.NewRequest(http.MethodPost, "/expenses", reader)
// 	recorder := httptest.NewRecorder()

// 	handler.ServeHTTP(recorder, req)

// 	if recorder.Code != http.StatusBadRequest {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusBadRequest,
// 			recorder.Code,
// 		)
// 	}
// }

// func TestCreateExpenseInternalServerError(t *testing.T) {
// 	fakeStore := &fakeExpenseStore{
// 		err: errors.New("store error"),
// 	}
// 	handler := NewExpenseHandler(fakeStore)

// 	expense := model.Expense{
// 		Title: "coffee",
// 		Amount: 500,
// 		Category: "food",
// 	}

// 	body, err := json.Marshal(expense)
// 	if err != nil {
// 		t.Fatalf("failed to marshal expense: %v", err)
// 	}
// 	reader := bytes.NewReader(body)

// 	req := httptest.NewRequest(http.MethodPost, "/expenses", reader)
// 	recorder := httptest.NewRecorder()

// 	handler.ServeHTTP(recorder, req)

// 	if recorder.Code != http.StatusInternalServerError {
// 		t.Fatalf(
// 			"expected status %d, got %d",
// 			http.StatusInternalServerError,
// 			recorder.Code,
// 		)
// 	}
// }

func TestCreateExpense(t *testing.T) {
	expense := model.Expense{
		Title: "coffee",
		Amount: 500,
		Category: "food",
	}

	validBody, err := json.Marshal(expense)
	if err != nil {
		t.Fatalf("failed to marshal expense: %v", err)
	}

	invalidBody := []byte(`{"title": "coffee"`)

	tests := []struct{
		name string
		store *fakeExpenseStore
		body []byte
		expectedStatus int
	}{
		{
			name: "success",
			store: &fakeExpenseStore{},
			body: validBody,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid json",
			store: &fakeExpenseStore{},
			body: invalidBody,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "store error",
			store: &fakeExpenseStore{
				err: errors.New("store error"),
			},
			body: validBody,
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			handler := NewExpenseHandler(tt.store)

			reader := bytes.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/expenses", reader)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if tt.expectedStatus == http.StatusCreated {
				var createdExpense model.Expense
				if err := json.NewDecoder(recorder.Body).Decode(&createdExpense); err != nil {
					t.Fatalf("failed to decode created expense: %v", err)
				}
				if createdExpense.ID != 1 {
					t.Fatalf(
						"expected id %d, got %d",
						1,
						createdExpense.ID,
					)
				}
				if createdExpense.Title != "coffee" {
					t.Fatalf(
						"expected title %q, got %q",
						"coffee",
						createdExpense.Title,
					)
				}
				if createdExpense.Amount != 500 {
					t.Fatalf(
						"expected amount %d, got %d",
						500,
						createdExpense.Amount,
					)
				}
				if createdExpense.Category != "food" {
					t.Fatalf(
						"expected category %q, got %q",
						"food",
						createdExpense.Category,
					)
				}
			}
		})
	}
}
