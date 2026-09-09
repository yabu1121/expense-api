//go:build integration

package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yabu1121/expense-api/internal/handler"
	"github.com/yabu1121/expense-api/internal/model"
)

func TestCreateExpenseIntegration(t *testing.T) {
	store := newExpenseTestStore(t)

	expenseHandler := handler.NewExpenseHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /expenses", expenseHandler.CreateExpense)

	t.Run("success", func(t *testing.T) {
		createdCategory, err := store.CreateCategory(model.CategoryRequest{
			Name: "food",
		})
		if err != nil {
			t.Fatalf("failed to create cateogory to store: %v", err)
		}

		expense := model.ExpenseRequest{
			Title:      "coffee",
			Amount:     500,
			CategoryID: createdCategory.ID,
		}

		expenseBody, err := json.Marshal(expense)
		if err != nil {
			t.Fatalf("failed to marshal expense: %v", err)
		}

		expenseReader := bytes.NewReader(expenseBody)

		req := httptest.NewRequest(http.MethodPost, "/expenses", expenseReader)
		recorder := httptest.NewRecorder()

		mux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, recorder.Code)
		}

		var createdExpense model.Expense
		if err := json.NewDecoder(recorder.Body).Decode(&createdExpense); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if createdExpense.Title != expense.Title {
			t.Errorf("expected title %s, got %s", expense.Title, createdExpense.Title)
		}

		if createdExpense.Amount != expense.Amount {
			t.Errorf("expected amount %d, got %d", expense.Amount, createdExpense.Amount)
		}

		got, err := store.GetExpenseByID(createdExpense.ID)
		if err != nil {
			t.Fatalf("failed to get expense by id: %v", err)
		}
		if got.ID != createdExpense.ID {
			t.Fatalf("expected %d, got %d", createdExpense.ID, got.ID)
		}
		if got.Title != createdExpense.Title {
			t.Fatalf("expected %s, got %s", createdExpense.Title, got.Title)
		}
		if got.Amount != createdExpense.Amount {
			t.Fatalf("expected %d, got %d", createdExpense.Amount, got.Amount)
		}
	})
}

// func TestGetExpenseByIDIntegration(t *testing.T) {
// 	tests := []struct {
// 		name              string
// 		param             string
// 		existingExpenses  []model.Expense
// 		useCreatedExpense bool
// 		expectedStatus    int
// 	}{
// 		{
// 			name:           "not found",
// 			param:          "999",
// 			expectedStatus: http.StatusNotFound,
// 		},
// 		{
// 			name:           "invalid param",
// 			param:          "fafdsa",
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name: "success",
// 			existingExpenses: []model.Expense{
// 				{
// 					Title:    "coffee",
// 					Amount:   500,
// 					Category: "food",
// 				},
// 			},
// 			useCreatedExpense: true,
// 			expectedStatus:    http.StatusOK,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			expenseStore := newExpenseTestStore(t)

// 			expenseHandler := handler.NewExpenseHandler(expenseStore)

// 			mux := http.NewServeMux()
// 			mux.HandleFunc("GET /expenses/{id}", expenseHandler.GetExpenseByID)

// 			var createdExpenses []model.Expense
// 			for _, e := range tt.existingExpenses {
// 				createdExpense, err := expenseStore.CreateExpense(e)
// 				if err != nil {
// 					t.Fatalf("failed to create expense: %v", err)
// 				}
// 				createdExpenses = append(createdExpenses, *createdExpense)
// 			}

// 			requestParam := tt.param
// 			if tt.useCreatedExpense {
// 				requestParam = strconv.Itoa(createdExpenses[0].ID)
// 			}

// 			recorder := httptest.NewRecorder()
// 			req := httptest.NewRequest(http.MethodGet, "/expenses/"+requestParam, nil)

// 			mux.ServeHTTP(recorder, req)

// 			if recorder.Code != tt.expectedStatus {
// 				t.Fatalf(
// 					"expected %d, got %d",
// 					tt.expectedStatus,
// 					recorder.Code,
// 				)
// 			}

// 			if recorder.Code == http.StatusOK {
// 				var gotExpense model.Expense
// 				if err := json.NewDecoder(recorder.Body).Decode(&gotExpense); err != nil {
// 					t.Fatalf("failed to get expense by id: %v", err)
// 				}

// 				if gotExpense.ID != createdExpenses[0].ID {
// 					t.Fatalf(
// 						"expected %d, got %d",
// 						createdExpenses[0].ID,
// 						gotExpense.ID,
// 					)
// 				}
// 				if gotExpense.Title != createdExpenses[0].Title {
// 					t.Fatalf(
// 						"expected %s, got %s",
// 						createdExpenses[0].Title,
// 						gotExpense.Title,
// 					)
// 				}

// 				if gotExpense.Amount != createdExpenses[0].Amount {
// 					t.Fatalf(
// 						"expected %d, got %d",
// 						createdExpenses[0].Amount,
// 						gotExpense.Amount,
// 					)
// 				}

// 				if gotExpense.Category != createdExpenses[0].Category {
// 					t.Fatalf(
// 						"expected %s, got %s",
// 						createdExpenses[0].Category,
// 						gotExpense.Category,
// 					)
// 				}
// 			}
// 		})
// 	}
// }

// func TestDeleteExpenseByIDIntegration(t *testing.T) {
// 	tempDir := t.TempDir()

// 	filePath := filepath.Join(tempDir, "expenses.db")

// 	expenseStore, err := store.NewSQLiteStore(filePath)
// 	if err != nil {
// 		t.Fatalf("failed to create store: %v", err)
// 	}

// 	t.Cleanup(func() {
// 		expenseStore.Close()
// 	})

// 	expenseHandler := handler.NewExpenseHandler(expenseStore)

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("DELETE /expenses/{id}", expenseHandler.DeleteExpenseByID)

// 	t.Run("success", func(t *testing.T) {
// 		// arrange
// 		expense := model.Expense{
// 			Title:    "coffee",
// 			Amount:   500,
// 			Category: "food",
// 		}

// 		createdExpense, err := expenseStore.CreateExpense(expense)
// 		if err != nil {
// 			t.Fatalf("failed to create expense to the store: %v", err)
// 		}

// 		// act delete
// 		deleteReq := httptest.NewRequest(http.MethodDelete, "/expenses/"+strconv.Itoa(createdExpense.ID), nil)
// 		deleteRecorder := httptest.NewRecorder()
// 		mux.ServeHTTP(deleteRecorder, deleteReq)

// 		if deleteRecorder.Code != http.StatusNoContent {
// 			t.Fatalf(
// 				"expected to %d, got %d",
// 				http.StatusNoContent,
// 				deleteRecorder.Code,
// 			)
// 		}

// 		// assert
// 		got, err := expenseStore.GetExpenseByID(createdExpense.ID)
// 		if !errors.Is(err, model.ErrExpenseNotFound) {
// 			t.Fatalf("expected ErrExpenseNotFound, got %v", err)
// 		}

// 		if got != nil {
// 			t.Fatal("expected expense to be nil")
// 		}
// 	})
// }

// func TestUpdateExpenseIntegration(t *testing.T) {
// 	tempDir := t.TempDir()

// 	filePath := filepath.Join(tempDir, "expenses.db")

// 	expenseStore, err := store.NewSQLiteStore(filePath)
// 	if err != nil {
// 		t.Fatalf("failed to create store: %v", err)
// 	}

// 	t.Cleanup(func() {
// 		expenseStore.Close()
// 	})

// 	expenseHandler := handler.NewExpenseHandler(expenseStore)

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("PUT /expenses/{id}", expenseHandler.UpdateExpenseByID)

// 	t.Run("success", func(t *testing.T) {
// 		// arrange
// 		expense := model.Expense{
// 			Title:    "coffee",
// 			Amount:   500,
// 			Category: "food",
// 		}

// 		// create to the store
// 		createdExpense, err := expenseStore.CreateExpense(expense)
// 		if err != nil {
// 			t.Fatalf("failed to create expense in the store: %v", err)
// 		}

// 		// update in http
// 		pendingUpdateExpense := model.Expense{
// 			Title:    "latte",
// 			Amount:   550,
// 			Category: "food",
// 		}

// 		pendingUpdateExpenseBody, err := json.Marshal(pendingUpdateExpense)
// 		if err != nil {
// 			t.Fatalf("failed to marshal update expense: %v", err)
// 		}

// 		pendingUpdateExpenseReader := bytes.NewReader(pendingUpdateExpenseBody)
// 		updateReq := httptest.NewRequest(http.MethodPut, "/expenses/"+strconv.Itoa(createdExpense.ID), pendingUpdateExpenseReader)
// 		updateRecorder := httptest.NewRecorder()
// 		mux.ServeHTTP(updateRecorder, updateReq)

// 		if updateRecorder.Code != http.StatusOK {
// 			t.Fatalf(
// 				"expected to %d, got %d",
// 				http.StatusOK,
// 				updateRecorder.Code,
// 			)
// 		}

// 		var updatedExpense model.Expense
// 		if err := json.NewDecoder(updateRecorder.Body).Decode(&updatedExpense); err != nil {
// 			t.Fatalf("failed to decode body: %v", err)
// 		}

// 		// assert
// 		got, err := expenseStore.GetExpenseByID(createdExpense.ID)
// 		if err != nil {
// 			t.Fatalf("failed to get expense by id from store: %v", err)
// 		}

// 		if got.Title != pendingUpdateExpense.Title {
// 			t.Fatalf(
// 				"expected expense title %s, got %s",
// 				pendingUpdateExpense.Title,
// 				got.Title,
// 			)
// 		}
// 		if got.Amount != pendingUpdateExpense.Amount {
// 			t.Fatalf(
// 				"expected expense amount %d, got %d",
// 				pendingUpdateExpense.Amount,
// 				got.Amount,
// 			)
// 		}
// 		if got.Category != pendingUpdateExpense.Category {
// 			t.Fatalf(
// 				"expected expense category %s, got %s",
// 				pendingUpdateExpense.Category,
// 				got.Category,
// 			)
// 		}
// 	})
// }
