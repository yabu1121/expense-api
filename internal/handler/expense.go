package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/yabu1121/expense-api/internal/model"
)

type ExpenseStore interface {
	GetAllExpenses() ([]model.Expense, error)
	GetExpenseByID(id int) (*model.Expense, error)
	CreateExpense(expense model.Expense) (*model.Expense, error)
	UpdateExpense(expense model.Expense) (*model.Expense, error)
	DeleteExpense(id int) error
}

type ExpenseHandler struct {
	store ExpenseStore
}

func NewExpenseHandler(store ExpenseStore) *ExpenseHandler {
	return &ExpenseHandler{
		store: store,
	}
}

func (h *ExpenseHandler) getAllExpense(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.store.GetAllExpenses()
	if err != nil {
		log.Printf("failed to get expenses: %v", err)
		http.Error(
			w,
			"failed to get expenses",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expenses); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ExpenseHandler) createExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	trimmedTitle := strings.TrimSpace(expense.Title)
	expense.Title = trimmedTitle

	if expense.Title == "" {
		http.Error(
			w,
			"title must not be empty",
			http.StatusBadRequest,
		)
		return
	}

	if expense.Amount <= 0 {
		http.Error(
			w,
			"amount must be greater than zero",
			http.StatusBadRequest,
		)
		return
	}

	trimmedCategory := strings.TrimSpace(expense.Category)
	expense.Category = trimmedCategory

	if expense.Category == "" {
		http.Error(
			w,
			"category must not be empty",
			http.StatusBadRequest,
		)
		return
	}

	createdExpense, err := h.store.CreateExpense(expense)
	if err != nil {
		log.Printf("failed to create expense: %v", err)
		http.Error(
			w,
			"failed to create expense",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdExpense); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ExpenseHandler) getExpenseByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			"invalid expense id",
			http.StatusBadRequest,
		)
		return
	}

	expense, err := h.store.GetExpenseByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(
				w,
				"expense not found",
				http.StatusNotFound,
			)
			return
		}
		log.Println(err)
		http.Error(
			w,
			"failed to get expense",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expense); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ExpenseHandler) deleteExpenseByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			"invalid expense id",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.store.DeleteExpense(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(
				w,
				"expense not found",
				http.StatusNotFound,
			)
			return
		}
		log.Printf("failed to delete expense: %v", err)
		http.Error(
			w,
			"failed to delete expense",
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) updateExpenseByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			"invalid expense id",
			http.StatusBadRequest,
		)
		return
	}

	var expense model.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}
	expense.ID = id

	updatedExpense, err := h.store.UpdateExpense(expense)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(
				w,
				"expense not found",
				http.StatusNotFound,
			)
			return
		}
		log.Printf("failed to update expenses: %v", err)
		http.Error(
			w,
			"failed to update expense",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updatedExpense); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ExpenseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		if id != "" {
			h.getExpenseByID(w, r)
			return
		}
		h.getAllExpense(w, r)
		return
	case http.MethodPost:
		if id != "" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.createExpense(w, r)
		return
	case http.MethodDelete:
		if id != "" {
			h.deleteExpenseByID(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodPut:
		if id != "" {
			h.updateExpenseByID(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
