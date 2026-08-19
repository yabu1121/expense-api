package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestDeleteExpenseByIDIntegration(t *testing.T) {
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
		// arrange
		expense := model.Expense{
			Title:    "coffee",
			Amount:   500,
			Category: "food",
		}

		createdExpense, err := expenseStore.CreateExpense(expense)
		if err != nil {
			t.Fatalf("failed to create expense to the store: %v", err)
		}

		// act delete
		deleteReq := httptest.NewRequest(http.MethodDelete, "/expenses/"+strconv.Itoa(createdExpense.ID), nil)
		deleteRecorder := httptest.NewRecorder()
		mux.ServeHTTP(deleteRecorder, deleteReq)

		if deleteRecorder.Code != http.StatusNoContent {
			t.Fatalf(
				"expected to %d, got %d",
				http.StatusNoContent,
				deleteRecorder.Code,
			)
		}

		// assert
		got, err := expenseStore.GetExpenseByID(createdExpense.ID)
		if !errors.Is(err, model.ErrExpenseNotFound) {
			t.Fatalf("expected ErrExpenseNotFound, got %v", err)
		}

		if got != nil {
			t.Fatal("expected expense to be nil")
		}
	})
}

func TestUpdateExpenseIntegration(t *testing.T) {
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
		// arrange
		expense := model.Expense{
			Title:    "coffee",
			Amount:   500,
			Category: "food",
		}

		// create to the store
		createdExpense, err := expenseStore.CreateExpense(expense)
		if err != nil {
			t.Fatalf("failed to create expense in the store: %v", err)
		}

		// update in http
		pendingUpdateExpense := model.Expense{
			Title: "latte",
			Amount: 550,
			Category: "food",
		}

		pendingUpdateExpenseBody, err := json.Marshal(pendingUpdateExpense)
		if err != nil {
			t.Fatalf("failed to marshal update expense: %v", err)
		}

		pendingUpdateExpenseReader := bytes.NewReader(pendingUpdateExpenseBody)
		updateReq := httptest.NewRequest(http.MethodPut, "/expenses/" + strconv.Itoa(createdExpense.ID), pendingUpdateExpenseReader)
		updateRecorder := httptest.NewRecorder()
		mux.ServeHTTP(updateRecorder, updateReq)

		if updateRecorder.Code != http.StatusOK {
			t.Fatalf(
				"expected to %d, got %d",
				http.StatusOK,
				updateRecorder.Code,
			)
		}

		var updatedExpense model.Expense
		if err := json.NewDecoder(updateRecorder.Body).Decode(&updatedExpense); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		// assert
		got, err := expenseStore.GetExpenseByID(createdExpense.ID)
		if err != nil {
			t.Fatalf("failed to get expense by id from store: %v", err)
		}

		if got.Title != pendingUpdateExpense.Title {
			t.Fatalf(
				"expected expense title %s, got %s",
				pendingUpdateExpense.Title,
				got.Title,
			)
		}
		if got.Amount != pendingUpdateExpense.Amount {
			t.Fatalf(
				"expected expense amount %d, got %d",
				pendingUpdateExpense.Amount,
				got.Amount,
			)
		}
		if got.Category != pendingUpdateExpense.Category {
			t.Fatalf(
				"expected expense category %s, got %s",
				pendingUpdateExpense.Category,
				got.Category,
			)
		}
	})
}
