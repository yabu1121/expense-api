package handler

import (
	"net/http"
	"uuid"
)

type ExpenseTagStore interface {
	AddTagToExpense(expenseID, tagID uuid.UUID) error
}

type ExpenseTagHandler struct {
	store ExpenseTagStore
}

func NewExpenseTagHandler(store ExpenseTagStore) *ExpenseTagHandler {
	return &ExpenseTagHandler{
		store: store,
	}
}

func (h *ExpenseTagHandler) AddTagToExpense(w http.ResponseWriter, r *http.Request) {
	expenseIDParam := r.PathValue("expenseID")
	expenseID, err := uuid.Parse(expenseIDParam)
	if err != nil {
		http.Error(
			w,
			"failed to parse expense id",
			http.StatusBadRequest,
		)
		return
	}

	tagIDParam := r.PathValue("tagID")
	tagID, err := uuid.Parse(tagIDParam)
	if err != nil {
		http.Error(
			w,
			"failed to uuid parse tag id param",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.store.AddTagToExpense(expenseID, tagID); err != nil {
		http.Error(
			w,
			"failed to add tag to expense",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
