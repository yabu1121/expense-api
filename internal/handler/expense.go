package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uuid"

	"github.com/yabu1121/expense-api/internal/model"
)

type ExpenseStore interface {
	ListExpenses() ([]model.Expense, error)
	GetExpenseByID(id uuid.UUID) (*model.Expense, error)
	CreateExpense(expense model.ExpenseRequest) (*model.Expense, error)
	UpdateExpense(id uuid.UUID, expense model.ExpenseRequest) (*model.Expense, error)
	DeleteExpense(id uuid.UUID) error
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
	expenses, err := h.store.ListExpenses()
	if err != nil {
		log.Printf("failed to get expenses: %v", err)
		http.Error(w, "failed to get expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expenses); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.ExpenseRequest

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	expense.Normalize()
	if err := expense.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdExpense, err := h.store.CreateExpense(expense)
	if err != nil {
		log.Printf("failed to create expense: %v", err)
		http.Error(w, "failed to create expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdExpense); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ExpenseHandler) GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
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
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		log.Println(err)
		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expense); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *ExpenseHandler) DeleteExpenseByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteExpense(id); err != nil {
		if errors.Is(err, model.ErrExpenseNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		log.Printf("failed to delete expense: %v", err)
		http.Error(w, "failed to delete expense", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) UpdateExpenseByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var expense model.ExpenseRequest

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	expense.Normalize()
	if err := expense.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedExpense, err := h.store.UpdateExpense(id, expense)
	if err != nil {
		if errors.Is(err, model.ErrExpenseNotFound) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}
		log.Printf("failed to update expenses: %v", err)
		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updatedExpense); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
