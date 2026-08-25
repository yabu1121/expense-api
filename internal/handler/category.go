package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yabu1121/expense-api/internal/model"
)

type CategoryStore interface {
	ListCategories() ([]model.Category, error)
	CreateCategory(category model.Category) (*model.Category, error)
}

type CategoryHandler struct {
	store CategoryStore
}

func NewCategoryHandler(store CategoryStore) *CategoryHandler {
	return &CategoryHandler{
		store: store,
	}
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.ListCategories()
	if err != nil {
		http.Error(
			w,
			"failed to list categories",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(
			w,
			"failed to decode request body",
			http.StatusBadRequest,
		)
		return
	}

	category.Normalize()
	if err := category.Validate(); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	createdCategory, err := h.store.CreateCategory(category)
	if err != nil {
		http.Error(
			w,
			"failed to create category",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdCategory); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}