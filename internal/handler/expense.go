package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/yabu1121/expense-api/internal/model"
)

type ExpenseStore interface {
	ListExpenses(filter model.ExpenseFilter) ([]model.Expense, error)
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

func (h *ExpenseHandler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	categoryFilter := r.URL.Query().Get("category")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	filter := model.ExpenseFilter{
		Category: categoryFilter,
	}

	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			http.Error(
				w,
				"invalid limit",
				http.StatusBadRequest,
			)
			return
		}
		filter.Limit = limit
	}

	if offsetParam != "" {
		if limitParam == "" {
			http.Error(
				w,
				"offset requires limit",
				http.StatusBadRequest,
			)
			return
		}
		offset, err := strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			http.Error(
				w,
				"invalid offset",
				http.StatusBadRequest,
			)
			return
		}
		filter.Offset = offset
	}

	expenses, err := h.store.ListExpenses(filter)
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

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	expense.Normalize()
	if err := expense.Validate(); err != nil {
		http.Error(
			w,
			err.Error(),
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

func (h *ExpenseHandler) GetExpenseByID(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, model.ErrExpenseNotFound) {
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

func (h *ExpenseHandler) DeleteExpenseByID(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, model.ErrExpenseNotFound) {
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

func (h *ExpenseHandler) UpdateExpenseByID(w http.ResponseWriter, r *http.Request) {
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

	expense.Normalize()
	if err := expense.Validate(); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	updatedExpense, err := h.store.UpdateExpense(expense)
	if err != nil {
		if errors.Is(err, model.ErrExpenseNotFound) {
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
