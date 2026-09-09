package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uuid"

	"github.com/yabu1121/expense-api/internal/model"
)

type CategoryStore interface {
	ListCategories() ([]model.Category, error)
	GetCategoryByID(id uuid.UUID) (*model.Category, error)
	CreateCategory(category model.CategoryRequest) (*model.Category, error)
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
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
	var categoryRequest CreateCategoryRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&categoryRequest); err != nil {
		http.Error(
			w,
			"failed to decode request body",
			http.StatusBadRequest,
		)
		return
	}

	category := model.Category{
		Name: categoryRequest.Name,
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

	createdCategory, err := h.store.CreateCategory(model.CategoryRequest{Name: category.Name})
	if err != nil {
		if errors.Is(err, model.ErrCategoryAlreadyExists) {
			http.Error(
				w,
				err.Error(),
				http.StatusConflict,
			)
			return
		}
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

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		http.Error(
			w,
			"failed to parse uuid",
			http.StatusBadRequest,
		)
		return
	}
	category, err := h.store.GetCategoryByID(categoryID)
	if err != nil {
		if errors.Is(err, model.ErrCategoryNotFound) {
			http.Error(
				w,
				err.Error(),
				http.StatusNotFound,
			)
			return
		}
		http.Error(
			w,
			"failed to get category by id",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(category); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
