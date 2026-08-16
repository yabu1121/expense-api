package model

import (
	"errors"
	"testing"
)

func TestExpenseNormalize(t *testing.T) {
	spaceOnlyTitleExpense := Expense{
		Title:    " ",
		Amount:   500,
		Category: "food",
	}

	spaceOnlyCategoryExpense := Expense{
		Title:    "coffee",
		Amount:   500,
		Category: " ",
	}

	tests := []struct {
		name            string
		expense         Expense
		expectedExpense Expense
	}{
		{
			name:    "space only title",
			expense: spaceOnlyTitleExpense,
			expectedExpense: Expense{
				Title:    "",
				Amount:   500,
				Category: "food",
			},
		},

		{
			name:    "space only category",
			expense: spaceOnlyCategoryExpense,
			expectedExpense: Expense{
				Title:    "coffee",
				Amount:   500,
				Category: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expense.Normalize()
			if tt.expectedExpense != tt.expense {
				t.Fatalf("failed to normalize target")
			}
		})
	}
}

func TestExpenseValidation(t *testing.T) {
	validExpense := Expense{
		Title:    "coffee",
		Amount:   500,
		Category: "food",
	}

	emptyTitleExpense := Expense{
		Title:    "",
		Amount:   500,
		Category: "food",
	}

	zeroAmountExpense := Expense{
		Title:    "coffee",
		Amount:   0,
		Category: "food",
	}

	negativeAmountExpense := Expense{
		Title:    "coffee",
		Amount:   -1,
		Category: "food",
	}

	emptyCategoryExpense := Expense{
		Title:    "coffee",
		Amount:   500,
		Category: "",
	}

	tests := []struct {
		name    string
		expense Expense
		wantErr error
	}{
		{
			name:    "success",
			expense: validExpense,
			wantErr: nil,
		},
		{
			name:    "empty title",
			expense: emptyTitleExpense,
			wantErr: ErrTitleRequired,
		},
		{
			name:    "zero amount",
			expense: zeroAmountExpense,
			wantErr: ErrAmountMustBePositive,
		},
		{
			name:    "negative amount",
			expense: negativeAmountExpense,
			wantErr: ErrAmountMustBePositive,
		},
		{
			name:    "empty category",
			expense: emptyCategoryExpense,
			wantErr: ErrCategoryRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expense.Normalize()
			err := tt.expense.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected no error, got: %v", err)
			}
		})
	}
}
