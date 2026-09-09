//go:build integration

package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/yabu1121/expense-api/internal/handler"
	"github.com/yabu1121/expense-api/internal/model"
)

func TestCreateCategoryIntegration(t *testing.T) {
	tests := []struct {
		name               string
		existingCategories []model.Category
		category           handler.CreateCategoryRequest
		expectedStatus     int
	}{
		{
			name: "success",
			category: handler.CreateCategoryRequest{
				Name: "food",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "conflict",
			existingCategories: []model.Category{
				{
					Name: "food",
				},
			},
			category: handler.CreateCategoryRequest{
				Name: "food",
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryStore := newCategoryTestStore(t)

			categoryHandler := handler.NewCategoryHandler(categoryStore)

			mux := http.NewServeMux()
			mux.HandleFunc("POST /categories", categoryHandler.CreateCategory)

			for _, c := range tt.existingCategories {
				_, err := categoryStore.CreateCategory(model.CategoryRequest{Name: c.Name})
				if err != nil {
					t.Fatalf("failed to create category to the store: %v", err)
				}
			}

			categoryBody, err := json.Marshal(tt.category)
			if err != nil {
				t.Fatalf("failed to marshal category: %v", err)
			}

			categoryReader := bytes.NewReader(categoryBody)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/categories", categoryReader)

			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if recorder.Code != http.StatusCreated {
				return
			}

			var createdCategory model.Category
			if err := json.NewDecoder(recorder.Body).Decode(&createdCategory); err != nil {
				t.Fatalf("failed to decode body: %v", err)
			}

			if createdCategory.Name != tt.category.Name {
				t.Fatalf(
					"expected category name %s, got %s",
					tt.category.Name,
					createdCategory.Name,
				)
			}

			if createdCategory.ID == uuid.Nil() {
				t.Fatal("expected category ID to be generated")
			}
		})
	}
}

func TestGetCategoryByIDIntegration(t *testing.T) {
	tests := []struct {
		name               string
		existingCategories []model.Category
		param              string
		expectedStatus     int
	}{
		{
			name: "success",
			existingCategories: []model.Category{
				{
					Name: "food",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "invalid uuid",
			param: "u-u-id",
			existingCategories: []model.Category{
				{
					Name: "food",
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "not found",
			param: uuid.NewV4().String(),
			existingCategories: []model.Category{
				{
					Name: "food",
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			categoryStore := newCategoryTestStore(t)

			categoryHandler := handler.NewCategoryHandler(categoryStore)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /categories/{id}", categoryHandler.GetCategoryByID)

			var createdCategories []model.Category
			for _, c := range tt.existingCategories {
				createdCategory, err := categoryStore.CreateCategory(model.CategoryRequest{Name: c.Name})
				if err != nil {
					t.Fatalf("failed to create category to the store: %v", err)
				}
				createdCategories = append(createdCategories, *createdCategory)
			}

			requestParam := tt.param
			if requestParam == "" {
				requestParam = createdCategories[len(createdCategories)-1].ID.String()
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/categories/"+requestParam, nil)

			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if recorder.Code == http.StatusOK {
				var got model.Category
				if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode category: %v", err)
				}

				if got != createdCategories[len(createdCategories)-1] {
					t.Fatalf(
						"expected %+v, got %+v",
						createdCategories[len(createdCategories)-1],
						got,
					)
				}
			}
		})
	}
}
