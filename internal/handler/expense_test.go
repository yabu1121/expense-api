package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yabu1121/expense-api/internal/model"
)

type fakeExpenseStore struct {
	expenses []model.Expense
	err      error
}

func (f *fakeExpenseStore) GetAllExpenses() ([]model.Expense, error) {
	return f.expenses, f.err
}

func (f *fakeExpenseStore) GetExpenseByID(id int) (*model.Expense, error) {
	if f.err != nil {
		return nil, f.err
	}

	for i := range f.expenses {
		if f.expenses[i].ID == id {
			return &f.expenses[i], nil
		}
	}
	return nil, model.ErrExpenseNotFound
}

func (f *fakeExpenseStore) CreateExpense(expense model.Expense) (*model.Expense, error) {
	if f.err != nil {
		return nil, f.err
	}
	expense.ID = 1
	return &expense, nil
}

func (f *fakeExpenseStore) UpdateExpense(expense model.Expense) (*model.Expense, error) {
	if f.err != nil {
		return nil, f.err
	}

	id := expense.ID

	for i := range f.expenses {
		if f.expenses[i].ID == id {
			f.expenses[i] = expense
			return &f.expenses[i], nil
		}
	}
	return nil, model.ErrExpenseNotFound
}

func (f *fakeExpenseStore) DeleteExpense(id int) error {
	if f.err != nil {
		return f.err
	}

	var newExpenses []model.Expense
	var flag bool

	for i := range f.expenses {
		if f.expenses[i].ID == id {
			flag = true
			continue
		}
		newExpenses = append(newExpenses, f.expenses[i])
	}

	f.expenses = newExpenses

	if !flag {
		return model.ErrExpenseNotFound
	}

	return nil
}

func (f *fakeExpenseStore) GetExpenseSummary() (*model.ExpenseSummary, error) {
	if f.err != nil {
		return nil, f.err
	}

	var count, total_amount int

	for _, expense := range f.expenses {
		count += 1
		total_amount += expense.Amount
	}

	return &model.ExpenseSummary{
		Count:       count,
		TotalAmount: total_amount,
	}, nil
}

func TestGetAllExpenses(t *testing.T) {
	tests := []struct {
		name           string
		store          *fakeExpenseStore
		expectedStatus int
	}{
		{
			name: "success",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "coffee",
						Amount:   500,
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
	tests := []struct {
		id             string
		name           string
		store          *fakeExpenseStore
		expectedStatus int
	}{
		{
			id:   "1",
			name: "success",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "coffee",
						Amount:   500,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			id:             "999",
			name:           "not found",
			store:          &fakeExpenseStore{},
			expectedStatus: http.StatusNotFound,
		},
		{
			id:             "abd",
			name:           "invalid id",
			store:          &fakeExpenseStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   "1",
			name: "store error",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "coffee",
						Amount:   500,
						Category: "food",
					},
				},
				err: errors.New("store error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewExpenseHandler(tt.store)

			req := httptest.NewRequest(http.MethodGet, "/expenses/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status to %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if tt.expectedStatus == http.StatusOK {
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
		})
	}
}

func TestCreateExpense(t *testing.T) {
	expense := model.Expense{
		Title:    "coffee",
		Amount:   500,
		Category: "food",
	}

	validBody, err := json.Marshal(expense)
	if err != nil {
		t.Fatalf("failed to marshal expense: %v", err)
	}

	invalidBody := []byte(`{"title": "coffee"`)

	emptyTitleExpense := model.Expense{
		Title:    "",
		Amount:   500,
		Category: "food",
	}

	emptyTitleBody, err := json.Marshal(emptyTitleExpense)
	if err != nil {
		t.Fatalf("failed to marshal empty title expense: %v", err)
	}

	spaceOnlyTitleExpense := model.Expense{
		Title:    "　",
		Amount:   500,
		Category: "food",
	}

	spaceOnlyTitleBody, err := json.Marshal(spaceOnlyTitleExpense)
	if err != nil {
		t.Fatalf("failed to marshal space only title expense: %v", err)
	}

	zeroAmountExpense := model.Expense{
		Title:    "coffee",
		Amount:   0,
		Category: "food",
	}

	zeroAmountBody, err := json.Marshal(zeroAmountExpense)
	if err != nil {
		t.Fatalf("failed to marshal zero amount expense: %v", err)
	}

	negativeAmountExpense := model.Expense{
		Title:    "coffee",
		Amount:   -1,
		Category: "food",
	}

	negativeAmountBody, err := json.Marshal(negativeAmountExpense)
	if err != nil {
		t.Fatalf("failed to marshal negative amount expense: %v", err)
	}

	emptyCategoryExpense := model.Expense{
		Title:    "coffee",
		Amount:   500,
		Category: "",
	}

	emptyCategoryBody, err := json.Marshal(emptyCategoryExpense)
	if err != nil {
		t.Fatalf("failed to marshal empty category expense: %v", err)
	}

	spaceOnlyCategoryExpense := model.Expense{
		Title:    "coffee",
		Amount:   500,
		Category: " ",
	}

	spaceOnlyCategoryBody, err := json.Marshal(spaceOnlyCategoryExpense)
	if err != nil {
		t.Fatalf("failed to marshal space only category expense: %v", err)
	}

	tests := []struct {
		name           string
		store          *fakeExpenseStore
		body           []byte
		expectedStatus int
	}{
		{
			name:           "success",
			store:          &fakeExpenseStore{},
			body:           validBody,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			store:          &fakeExpenseStore{},
			body:           invalidBody,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "store error",
			store: &fakeExpenseStore{
				err: errors.New("store error"),
			},
			body:           validBody,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "empty title",
			store:          &fakeExpenseStore{},
			body:           emptyTitleBody,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "space only title",
			store:          &fakeExpenseStore{},
			body:           spaceOnlyTitleBody,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "zero amount",
			store:          &fakeExpenseStore{},
			body:           zeroAmountBody,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative amount",
			store:          &fakeExpenseStore{},
			body:           negativeAmountBody,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty category",
			store:          &fakeExpenseStore{},
			body:           emptyCategoryBody,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "space only category",
			store:          &fakeExpenseStore{},
			body:           spaceOnlyCategoryBody,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestDeleteExpenseByID(t *testing.T) {
	tests := []struct {
		id             string
		name           string
		store          *fakeExpenseStore
		expectedStatus int
	}{
		{
			id:   "1",
			name: "success",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "coffee",
						Amount:   500,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			id:             "999",
			name:           "not found",
			store:          &fakeExpenseStore{},
			expectedStatus: http.StatusNotFound,
		},
		{
			id:             "fja",
			name:           "invalid id",
			store:          &fakeExpenseStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   "1",
			name: "store error",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "coffee",
						Amount:   500,
						Category: "food",
					},
				},
				err: errors.New("store error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewExpenseHandler(tt.store)

			req := httptest.NewRequest(http.MethodDelete, "/expenses/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status to %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}
		})
	}
}

func TestUpdateExpenseByID(t *testing.T) {
	expense := model.Expense{
		Title:    "coffee",
		Amount:   500,
		Category: "food",
	}

	validBody, err := json.Marshal(expense)
	if err != nil {
		t.Fatalf("failed to marshal expense %v", err)
	}

	invalidBody := []byte(`{"title": "coffee"`)

	emptyTitleExpense := model.Expense{
		ID:       1,
		Title:    "",
		Amount:   500,
		Category: "food",
	}

	emptyTitleExpenseBody, err := json.Marshal(emptyTitleExpense)
	if err != nil {
		t.Fatalf("failed to marshal empty title expense: %v", err)
	}

	tests := []struct {
		id             string
		name           string
		body           []byte
		store          fakeExpenseStore
		expectedStatus int
	}{
		{
			id:   "1",
			name: "success",
			body: validBody,
			store: fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "latte",
						Amount:   550,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			id:   "999",
			name: "not found",
			body: validBody,
			store: fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "latte",
						Amount:   550,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			id:   "adj",
			name: "invalid id",
			body: validBody,
			store: fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "latte",
						Amount:   550,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   "1",
			name: "store error",
			body: validBody,
			store: fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "latte",
						Amount:   550,
						Category: "food",
					},
				},
				err: errors.New("store error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   "1",
			name: "invalid json request",
			body: invalidBody,
			store: fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "latte",
						Amount:   550,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   "1",
			name: "empty title",
			body: emptyTitleExpenseBody,
			store: fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:       1,
						Title:    "latte",
						Amount:   550,
						Category: "food",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewExpenseHandler(&tt.store)

			reader := bytes.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPut, "/expenses/"+tt.id, reader)
			req.SetPathValue("id", tt.id)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if tt.expectedStatus == http.StatusOK {
				var updatedExpense model.Expense
				if err := json.NewDecoder(recorder.Body).Decode(&updatedExpense); err != nil {
					t.Fatalf("failed to decode updated expense: %v", err)
				}

				if updatedExpense.ID != 1 {
					t.Fatalf(
						"expected expense ID %d, got %d",
						1,
						updatedExpense.ID,
					)
				}
				if updatedExpense.Title != "coffee" {
					t.Fatalf(
						"expected expense title %s, got %s",
						"coffee",
						updatedExpense.Title,
					)
				}
				if updatedExpense.Amount != 500 {
					t.Fatalf(
						"expected expense amount %d, got %d",
						500,
						updatedExpense.Amount,
					)
				}
				if updatedExpense.Category != "food" {
					t.Fatalf(
						"expected expense category %s, got %s",
						"food",
						updatedExpense.Category,
					)
				}
			}
		})
	}
}

func TestGetExpenseSummary(t *testing.T) {
	tests := []struct {
		name           string
		store          *fakeExpenseStore
		expectedResult model.ExpenseSummary
		expectedStatus int
	}{
		{
			name: "success",
			store: &fakeExpenseStore{
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
			},
			expectedResult: model.ExpenseSummary{
				Count:       2,
				TotalAmount: 1050,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "success2",
			store: &fakeExpenseStore{
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
					{
						ID:       3,
						Title:    "moca",
						Amount:   700,
						Category: "food",
					},
				},
			},
			expectedResult: model.ExpenseSummary{
				Count:       3,
				TotalAmount: 1750,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewExpenseSummaryHandler(tt.store)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/expenses/summary", nil)

			handler.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			var res model.ExpenseSummary
			if err := json.NewDecoder(recorder.Body).Decode(&res); err != nil {
				t.Fatalf("failed to decode expense summary: %v", err)
			}

			if res.Count != tt.expectedResult.Count {
				t.Fatalf("result count is not matched")
			}

			if res.TotalAmount != tt.expectedResult.TotalAmount {
				t.Fatalf("result total amount is not matched")
			}
		})
	}
}
