package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/yabu1121/expense-api/internal/handler"
	"github.com/yabu1121/expense-api/internal/model"
	"github.com/yabu1121/expense-api/internal/store"
)

func TestCreateExpenseIntegration(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "expenses.db")

	expenseStore, err := store.NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	t.Cleanup(func() {
		expenseStore.Close()
	})

	expenseHandler := handler.NewExpenseHandler(expenseStore)

	mux := http.NewServeMux()
	mux.Handle("/expenses", expenseHandler)
	mux.Handle("/expenses/{id}", expenseHandler)

	t.Run("success", func(t *testing.T) {
		expense := model.Expense{
			Title:    "coffee",
			Amount:   500,
			Category: "food",
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
			t.Fatalf(
				"expected %d, got %d",
				http.StatusCreated,
				recorder.Code,
			)
		}

		var createdExpense model.Expense
		if err := json.NewDecoder(recorder.Body).Decode(&createdExpense); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if createdExpense.ID <= 0 {
			t.Errorf("expected id, must be positive")
		}

		if createdExpense.Title != expense.Title {
			t.Errorf(
				"expected title %s, got %s",
				expense.Title,
				createdExpense.Title,
			)
		}

		if createdExpense.Amount != expense.Amount {
			t.Errorf(
				"expected amount %d, got %d",
				expense.Amount,
				createdExpense.Amount,
			)
		}

		if createdExpense.Category != expense.Category {
			t.Errorf(
				"expected category %s, got %s",
				expense.Category,
				createdExpense.Category,
			)
		}

		got, err := expenseStore.GetExpenseByID(createdExpense.ID)
		if err != nil {
			t.Fatalf("failed to get expense by id: %v", err)
		}

		if got.ID != createdExpense.ID {
			t.Fatalf(
				"expected %d, got %d",
				createdExpense.ID,
				got.ID,
			)
		}
		if got.Title != createdExpense.Title {
			t.Fatalf(
				"expected %s, got %s",
				createdExpense.Title,
				got.Title,
			)
		}
		if got.Amount != createdExpense.Amount {
			t.Fatalf(
				"expected %d, got %d",
				createdExpense.Amount,
				got.Amount,
			)
		}
		if got.Category != createdExpense.Category {
			t.Fatalf(
				"expected %s, got %s",
				createdExpense.Category,
				got.Category,
			)
		}
	})
}

func TestGetExpenseByIDIntegration(t *testing.T) {
	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "expenses.db")

	expenseStore, err := store.NewSQLiteStore(filePath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	t.Cleanup(func() {
		expenseStore.Close()
	})

	expenseHandler := handler.NewExpenseHandler(expenseStore)

	mux := http.NewServeMux()
	mux.Handle("/expenses/{id}", expenseHandler)

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses/999", nil)
		recorder := httptest.NewRecorder()

		mux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusNotFound,
				recorder.Code,
			)
		}
	})

	t.Run("invalid param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/expenses/fadfa", nil)
		recorder := httptest.NewRecorder()

		mux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})

	t.Run("success", func(t *testing.T) {
		expense := model.Expense{
			ID:       1,
			Title:    "coffee",
			Amount:   500,
			Category: "food",
		}

		createdExpense, err := expenseStore.CreateExpense(expense)
		if err != nil {
			t.Fatalf("failed to create expense: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/expenses/"+strconv.Itoa(createdExpense.ID), nil)
		recorder := httptest.NewRecorder()

		mux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		var gotExpense model.Expense
		if err := json.NewDecoder(recorder.Body).Decode(&gotExpense); err != nil {
			t.Fatalf("failed to get expense by id: %v", err)
		}

		if gotExpense.ID != createdExpense.ID {
			t.Fatalf(
				"expected %d, got %d",
				createdExpense.ID,
				gotExpense.ID,
			)
		}

		if gotExpense.Title != createdExpense.Title {
			t.Fatalf(
				"expected %s, got %s",
				createdExpense.Title,
				gotExpense.Title,
			)
		}

		if gotExpense.Amount != createdExpense.Amount {
			t.Fatalf(
				"expected %d, got %d",
				createdExpense.Amount,
				gotExpense.Amount,
			)
		}

		if gotExpense.Category != createdExpense.Category {
			t.Fatalf(
				"expected %s, got %s",
				createdExpense.Category,
				gotExpense.Category,
			)
		}
	})
}
