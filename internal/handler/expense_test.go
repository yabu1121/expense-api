package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"uuid"

	"github.com/yabu1121/expense-api/internal/model"
)

type fakeExpenseStore struct {
	expenses []model.Expense
	err      error
}

func (f *fakeExpenseStore) ListExpenses() ([]model.Expense, error) {
	if f.err != nil {
		return nil, f.err
	}
	expenses := f.expenses

	return expenses, nil
}

func (f *fakeExpenseStore) GetExpenseByID(id uuid.UUID) (*model.Expense, error) {
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

func (f *fakeExpenseStore) CreateExpense(expense model.ExpenseRequest) (*model.Expense, error) {
	if f.err != nil {
		return nil, f.err
	}
	result := model.Expense{
		ID:         uuid.NewV7(),
		Title:      expense.Title,
		Amount:     expense.Amount,
		CategoryID: expense.CategoryID,
		CreatedAt:  time.Now(),
	}

	return &result, nil
}

func (f *fakeExpenseStore) UpdateExpense(id uuid.UUID, expense model.ExpenseRequest) (*model.Expense, error) {
	if f.err != nil {
		return nil, f.err
	}

	for i := range f.expenses {
		if f.expenses[i].ID == id {
			f.expenses[i].Title = expense.Title
			f.expenses[i].Amount = expense.Amount
			f.expenses[i].CategoryID = expense.CategoryID
			return &f.expenses[i], nil
		}
	}
	return nil, model.ErrExpenseNotFound
}

func (f *fakeExpenseStore) DeleteExpense(id uuid.UUID) error {
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

// func (f *fakeExpenseStore) GetExpenseSummary() (*model.ExpenseSummary, error) {
// 	if f.err != nil {
// 		return nil, f.err
// 	}

// 	var count, total_amount int

// 	for _, expense := range f.expenses {
// 		count += 1
// 		total_amount += expense.Amount
// 	}

// 	return &model.ExpenseSummary{
// 		Count:       count,
// 		TotalAmount: total_amount,
// 	}, nil
// }

func TestListExpenses(t *testing.T) {
	expenseID1 := uuid.NewV7()
	expenseID2 := uuid.NewV7()
	expenseID3 := uuid.NewV7()
	categoryFoodID := uuid.NewV7()

	store := &fakeExpenseStore{
		expenses: []model.Expense{
			{
				ID:         expenseID1,
				Title:      "coffee",
				Amount:     500,
				CategoryID: categoryFoodID,
				CreatedAt:  time.Now(),
			},
			{
				ID:         expenseID2,
				Title:      "latte",
				Amount:     550,
				CategoryID: categoryFoodID,
				CreatedAt:  time.Now(),
			},
			{
				ID:         expenseID3,
				Title:      "moca",
				Amount:     600,
				CategoryID: categoryFoodID,
				CreatedAt:  time.Now(),
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		handler := NewExpenseHandler(store)

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		recorder := httptest.NewRecorder()

		handler.ListExpenses(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("empty", func(t *testing.T) {
		expenses := make([]model.Expense, 0)
		emptyStore := &fakeExpenseStore{
			expenses: expenses,
		}
		handler := NewExpenseHandler(emptyStore)

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		recorder := httptest.NewRecorder()

		handler.ListExpenses(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("store error", func(t *testing.T) {
		errorStore := &fakeExpenseStore{
			err: errors.New("store error"),
		}

		handler := NewExpenseHandler(errorStore)

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		recorder := httptest.NewRecorder()

		handler.ListExpenses(recorder, req)

		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf("expected %d, got %d", http.StatusInternalServerError, recorder.Code)
		}
	})
}

func TestGetExpenseByID(t *testing.T) {
	id1 := uuid.NewV7()
	id2 := uuid.NewV7()
	id3 := uuid.NewV7()
	foodCategoryID := uuid.NewV7()
	tests := []struct {
		id             string
		name           string
		store          *fakeExpenseStore
		wantExpense    model.Expense
		expectedStatus int
	}{
		{
			id:   id2.String(),
			name: "success",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:         id1,
						Title:      "coffee",
						Amount:     500,
						CategoryID: foodCategoryID,
					},
					{
						ID:         id2,
						Title:      "latte",
						Amount:     550,
						CategoryID: foodCategoryID,
					},
				},
			},
			wantExpense: model.Expense{
				ID:         id2,
				Title:      "latte",
				Amount:     550,
				CategoryID: foodCategoryID,
			},
			expectedStatus: http.StatusOK,
		},
		{
			id:             id3.String(),
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
			id:   id1.String(),
			name: "store error",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:         id1,
						Title:      "coffee",
						Amount:     500,
						CategoryID: foodCategoryID,
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

			handler.GetExpenseByID(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("expected status to %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var expense model.Expense
				if err := json.NewDecoder(recorder.Body).Decode(&expense); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if expense != tt.wantExpense {
					t.Fatalf("want expense is not matched")
				}
			}
		})
	}
}

func TestCreateExpense(t *testing.T) {
	categoryID := uuid.NewV7()
	expense := model.ExpenseRequest{
		Title:      "coffee",
		Amount:     500,
		CategoryID: categoryID,
	}

	validBody, err := json.Marshal(expense)
	if err != nil {
		t.Fatalf("failed to marshal expense: %v", err)
	}

	invalidBody := []byte(`{"title": "coffee"`)

	emptyTitleExpense := model.Expense{
		Title:      "",
		Amount:     500,
		CategoryID: categoryID,
	}

	emptyTitleBody, err := json.Marshal(emptyTitleExpense)
	if err != nil {
		t.Fatalf("failed to marshal empty title expense: %v", err)
	}

	spaceOnlyTitleExpense := model.Expense{
		Title:      "　",
		Amount:     500,
		CategoryID: categoryID,
	}

	spaceOnlyTitleBody, err := json.Marshal(spaceOnlyTitleExpense)
	if err != nil {
		t.Fatalf("failed to marshal space only title expense: %v", err)
	}

	zeroAmountExpense := model.Expense{
		Title:      "coffee",
		Amount:     0,
		CategoryID: categoryID,
	}

	zeroAmountBody, err := json.Marshal(zeroAmountExpense)
	if err != nil {
		t.Fatalf("failed to marshal zero amount expense: %v", err)
	}

	negativeAmountExpense := model.Expense{
		Title:      "coffee",
		Amount:     -1,
		CategoryID: categoryID,
	}

	negativeAmountBody, err := json.Marshal(negativeAmountExpense)
	if err != nil {
		t.Fatalf("failed to marshal negative amount expense: %v", err)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewExpenseHandler(tt.store)

			reader := bytes.NewReader(tt.body)

			req := httptest.NewRequest(http.MethodPost, "/expenses", reader)
			recorder := httptest.NewRecorder()

			handler.CreateExpense(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var createdExpense model.Expense
				if err := json.NewDecoder(recorder.Body).Decode(&createdExpense); err != nil {
					t.Fatalf("failed to decode created expense: %v", err)
				}
				if createdExpense.ID == uuid.Nil() {
					t.Fatalf("expected id %q, got %q", uuid.Nil(), createdExpense.ID)
				}
				if createdExpense.Title != "coffee" {
					t.Fatalf("expected title %q, got %q", "coffee", createdExpense.Title)
				}
				if createdExpense.Amount != 500 {
					t.Fatalf("expected amount %d, got %d", 500, createdExpense.Amount)
				}
				if createdExpense.CategoryID != categoryID {
					t.Fatalf("expected category %q, got %q", categoryID, createdExpense.CategoryID)
				}
			}
		})
	}
}

func TestDeleteExpenseByID(t *testing.T) {
	id1 := uuid.NewV7()
	foodCategoryID := uuid.NewV7()
	tests := []struct {
		id             string
		name           string
		store          *fakeExpenseStore
		expectedStatus int
	}{
		{
			id:   id1.String(),
			name: "success",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:         id1,
						Title:      "coffee",
						Amount:     500,
						CategoryID: foodCategoryID,
					},
				},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			id:             id1.String(),
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
			id:   id1.String(),
			name: "store error",
			store: &fakeExpenseStore{
				expenses: []model.Expense{
					{
						ID:         id1,
						Title:      "coffee",
						Amount:     500,
						CategoryID: foodCategoryID,
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

			handler.DeleteExpenseByID(recorder, req)

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
	id1 := uuid.NewV7()
	foodCategoryID := uuid.NewV7()

	store := &fakeExpenseStore{
		expenses: []model.Expense{
			{
				ID:         id1,
				Title:      "coffee",
				Amount:     500,
				CategoryID: foodCategoryID,
			},
		},
	}

	expense := model.ExpenseRequest{
		Title:      "latte",
		Amount:     550,
		CategoryID: foodCategoryID,
	}

	validBody, err := json.Marshal(expense)
	if err != nil {
		t.Fatalf("failed to marshal expense %v", err)
	}

	t.Run("success", func(t *testing.T) {
		handler := NewExpenseHandler(store)

		reader := bytes.NewReader(validBody)

		req := httptest.NewRequest(http.MethodPut, "/expenses/"+id1.String(), reader)
		req.SetPathValue("id", id1.String())
		recorder := httptest.NewRecorder()

		handler.UpdateExpenseByID(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}
	})
}

// func TestGetExpenseSummary(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		store          *fakeExpenseStore
// 		expectedResult model.ExpenseSummary
// 		expectedStatus int
// 	}{
// 		{
// 			name: "success",
// 			store: &fakeExpenseStore{
// 				expenses: []model.Expense{
// 					{
// 						ID:       1,
// 						Title:    "coffee",
// 						Amount:   500,
// 						Category: "food",
// 					},
// 					{
// 						ID:       2,
// 						Title:    "latte",
// 						Amount:   550,
// 						Category: "food",
// 					},
// 				},
// 			},
// 			expectedResult: model.ExpenseSummary{
// 				Count:       2,
// 				TotalAmount: 1050,
// 			},
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name: "success2",
// 			store: &fakeExpenseStore{
// 				expenses: []model.Expense{
// 					{
// 						ID:       1,
// 						Title:    "coffee",
// 						Amount:   500,
// 						Category: "food",
// 					},
// 					{
// 						ID:       2,
// 						Title:    "latte",
// 						Amount:   550,
// 						Category: "food",
// 					},
// 					{
// 						ID:       3,
// 						Title:    "moca",
// 						Amount:   700,
// 						Category: "food",
// 					},
// 				},
// 			},
// 			expectedResult: model.ExpenseSummary{
// 				Count:       3,
// 				TotalAmount: 1750,
// 			},
// 			expectedStatus: http.StatusOK,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			handler := NewExpenseSummaryHandler(tt.store)

// 			recorder := httptest.NewRecorder()
// 			req := httptest.NewRequest(http.MethodGet, "/expenses/summary", nil)

// 			handler.GetExpenseSummary(recorder, req)

// 			if recorder.Code != tt.expectedStatus {
// 				t.Fatalf(
// 					"expected status %d, got %d",
// 					tt.expectedStatus,
// 					recorder.Code,
// 				)
// 			}

// 			var res model.ExpenseSummary
// 			if err := json.NewDecoder(recorder.Body).Decode(&res); err != nil {
// 				t.Fatalf("failed to decode expense summary: %v", err)
// 			}

// 			if res.Count != tt.expectedResult.Count {
// 				t.Fatalf("result count is not matched")
// 			}

// 			if res.TotalAmount != tt.expectedResult.TotalAmount {
// 				t.Fatalf("result total amount is not matched")
// 			}
// 		})
// 	}
// }
