package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yabu1121/expense-api/internal/model"
)

type ExpenseSummaryStore interface {
	GetExpenseSummary() (*model.ExpenseSummary, error)
}

type ExpenseSummaryHandler struct {
	store ExpenseSummaryStore
}

func NewExpenseSummaryHandler(store ExpenseSummaryStore) *ExpenseSummaryHandler {
	return &ExpenseSummaryHandler{
		store: store,
	}
}

func (h ExpenseSummaryHandler) getExpenseSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.store.GetExpenseSummary()
	if err != nil {
		log.Printf("failed to get expense summary: %v", err)
		http.Error(
			w,
			"failed to get expense summary",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(summary); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h ExpenseSummaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getExpenseSummary(w, r)
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
