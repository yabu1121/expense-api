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
		category           model.Category
		expectedStatus     int
	}{
		{
			name: "success",
			category: model.Category{
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
			category: model.Category{
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
				_, err := categoryStore.CreateCategory(c)
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
