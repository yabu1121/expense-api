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

func (f *fakeCategoryStore) CreateCategory(request model.CategoryRequest) (*model.Category, error) {
	if f.err != nil {
		return nil, f.err
	}
	category := model.Category{Name: request.Name}
	category.ID = uuid.NewV7()
	return &category, nil
}

func (f *fakeCategoryStore) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	if f.err != nil {
		return nil, f.err
	}

	for i := range f.categories {
		if f.categories[i].ID == id {
			return &f.categories[i], nil
		}
	}
	return nil, model.ErrCategoryNotFound
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
	validCategory := CreateCategoryRequest{
		Name: "food",
	}
	validBody, err := json.Marshal(validCategory)
	if err != nil {
		t.Fatalf("failed to marshal valid category: %v", err)
	}

	invalidBody := []byte(`{"name": "food"`)

	spaceOnlyCategory := CreateCategoryRequest{
		Name: " ",
	}
	spaceOnlyBody, err := json.Marshal(spaceOnlyCategory)
	if err != nil {
		t.Fatalf("failed to marshal space only category: %v", err)
	}

	unknownFieldBody := []byte(
		`{"id":"01900000-0000-7000-8000-000000000000","name":"food"}`,
	)

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
		{
			name:           "unknown field",
			body:           unknownFieldBody,
			store:          &fakeCategoryStore{},
			expectedStatus: http.StatusBadRequest,
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

func TestGetCategoryByID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		store          *fakeCategoryStore
		wantCategory   model.Category
		expectedStatus int
	}{
		{
			name: "success",
			id:   uuid.Max().String(),
			store: &fakeCategoryStore{
				categories: []model.Category{
					{
						ID:   uuid.Max(),
						Name: "food",
					},
				},
			},
			wantCategory: model.Category{
				ID:   uuid.Max(),
				Name: "food",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid uuid",
			id:   "invalid-uuid-dayo",
			store: &fakeCategoryStore{
				categories: []model.Category{
					{
						Name: "food",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			id:             uuid.Max().String(),
			store:          &fakeCategoryStore{},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "store error",
			id:   uuid.Max().String(),
			store: &fakeCategoryStore{
				err: errors.New("store error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewCategoryHandler(tt.store)

			req := httptest.NewRequest(http.MethodGet, "/categories/"+tt.id, nil)
			req.SetPathValue("id", (tt.id))
			recorder := httptest.NewRecorder()

			handler.GetCategoryByID(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status to %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if tt.expectedStatus == http.StatusOK {
				var category model.Category
				if err := json.NewDecoder(recorder.Body).Decode(&category); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if category != tt.wantCategory {
					t.Fatalf("category is not unmathced")
				}
			}
		})
	}
}
