package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"uuid"

	"github.com/yabu1121/expense-api/internal/model"
)

type fakeCategoryStore struct {
	categories []model.Category
	err        error
}

func (f *fakeCategoryStore) ListCategories() ([]model.Category, error) {
	return f.categories, f.err
}

func (f *fakeCategoryStore) CreateCategory(category model.Category) (*model.Category, error) {
	if f.err != nil {
		return nil, f.err
	}
	category.ID = uuid.NewV7()
	return &category, nil
}

func TestListCategories(t *testing.T) {
	tests := []struct {
		name           string
		store          *fakeCategoryStore
		expectedStatus int
	}{
		{
			name: "success",
			store: &fakeCategoryStore{
				categories: []model.Category{
					{
						ID:   uuid.NewV7(),
						Name: "food",
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "store error",
			store: &fakeCategoryStore{
				err: errors.New("store error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewCategoryHandler(tt.store)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/categories", nil)

			handler.ListCategories(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}
			if tt.expectedStatus != http.StatusOK {
				return
			}

			var res []model.Category
			if err := json.NewDecoder(recorder.Body).Decode(&res); err != nil {
				t.Fatalf("failed to decode recorder body: %v", err)
			}
			if !slices.Equal(res, tt.store.categories) {
				t.Fatalf("expected %+v, got %+v", tt.store.categories, res)
			}
		})
	}
}

func TestCreateCategory(t *testing.T) {
	validCategory := model.Category{
		Name: "food",
	}
	validBody, err := json.Marshal(validCategory)
	if err != nil {
		t.Fatalf("failed to marshal valid category: %v", err)
	}

	invalidBody := []byte(`{"name": "food"`)

	spaceOnlyCategory := model.Category{
		Name: " ",
	}
	spaceOnlyBody, err := json.Marshal(spaceOnlyCategory)
	if err != nil {
		t.Fatalf("failed to marshal space only category: %v", err)
	}

	tests := []struct {
		name           string
		body           []byte
		store          *fakeCategoryStore
		expectedStatus int
	}{
		{
			name:           "success",
			body:           validBody,
			store:          &fakeCategoryStore{},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid",
			body:           invalidBody,
			store:          &fakeCategoryStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "space only",
			body:           spaceOnlyBody,
			store:          &fakeCategoryStore{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "store error",
			body: validBody,
			store: &fakeCategoryStore{
				err: errors.New("store error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "conflict",
			body: validBody,
			store: &fakeCategoryStore{
				err: model.ErrCategoryAlreadyExists,
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reader := bytes.NewReader(tt.body)

			handler := NewCategoryHandler(tt.store)

			recoder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/categories", reader)

			handler.CreateCategory(recoder, req)

			if recoder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.expectedStatus,
					recoder.Code,
				)
			}

			if tt.expectedStatus != http.StatusCreated {
				return
			}

			var res model.Category
			if err := json.NewDecoder(recoder.Body).Decode(&res); err != nil {
				t.Fatalf("failed to decode recorder body: %v", err)
			}

			if res.Name != validCategory.Name {
				t.Fatalf("expected name %s, got %s", validCategory.Name, res.Name)
			}

			if res.ID == uuid.Nil() {
				t.Fatal("expected category ID to be generated")
			}
		})
	}
}
