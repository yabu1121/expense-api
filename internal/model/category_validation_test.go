package model

import (
	"errors"
	"testing"
)

func TestCategoryNormalize(t *testing.T) {
	validCategory := Category{
		Name: "food",
	}

	spaceOnlyNameCategory := Category{
		Name: " ",
	}

	expectedSpaceOnlyNameCategory := Category{
		Name: "",
	}

	notLowerNameCategory := Category{
		Name: "Food",
	}

	expectedNotLowerNameCategory := Category{
		Name: "food",
	}

	tests := []struct {
		name             string
		category         Category
		expectedCategory Category
	}{
		{
			name:             "valid",
			category:         validCategory,
			expectedCategory: validCategory,
		},
		{
			name:             "space only name",
			category:         spaceOnlyNameCategory,
			expectedCategory: expectedSpaceOnlyNameCategory,
		},
		{
			name:             "to lower case",
			category:         notLowerNameCategory,
			expectedCategory: expectedNotLowerNameCategory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.category.Normalize()
			if tt.category != tt.expectedCategory {
				t.Fatalf("failed to normalize target")
			}
		})
	}
}

func TestCategoryValidate(t *testing.T) {
	validCategory := Category{
		Name: "food",
	}

	emptyNameCategory := Category{
		Name: "",
	}

	spaceOnlyCategory := Category{
		Name: " ",
	}

	tests := []struct {
		name     string
		category Category
		wantErr  error
	}{
		{
			name:     "success",
			category: validCategory,
			wantErr:  nil,
		},
		{
			name:     "empty name",
			category: emptyNameCategory,
			wantErr:  ErrNameRequired,
		},
		{
			name:     "space only name",
			category: spaceOnlyCategory,
			wantErr:  ErrNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.category.Normalize()
			err := tt.category.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected no error, got: %v", err)
			}
		})
	}
}
