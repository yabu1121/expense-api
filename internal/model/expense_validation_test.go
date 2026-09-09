package model

import (
	"errors"
	"testing"
	"uuid"
)

func TestExpenseNormalize(t *testing.T) {
	spaceOnlyTitleExpense := ExpenseRequest{
		Title:      " ",
		Amount:     500,
		CategoryID: uuid.Nil(),
	}

	spaceOnlyCategoryExpense := ExpenseRequest{
		Title:      "coffee",
		Amount:     500,
		CategoryID: uuid.Nil(),
	}

	tests := []struct {
		name            string
		expense         ExpenseRequest
		expectedExpense ExpenseRequest
	}{
		{
			name:    "space only title",
			expense: spaceOnlyTitleExpense,
			expectedExpense: ExpenseRequest{
				Title:      "",
				Amount:     500,
				CategoryID: uuid.Nil(),
			},
		},
		{
			name:    "space only category",
			expense: spaceOnlyCategoryExpense,
			expectedExpense: ExpenseRequest{
				Title:      "coffee",
				Amount:     500,
				CategoryID: uuid.Nil(),
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
	categoryID := uuid.NewV7()
	validExpense := ExpenseRequest{
		Title:      "coffee",
		Amount:     500,
		CategoryID: categoryID,
	}

	emptyTitleExpense := ExpenseRequest{
		Title:      "",
		Amount:     500,
		CategoryID: categoryID,
	}

	zeroAmountExpense := ExpenseRequest{
		Title:      "coffee",
		Amount:     0,
		CategoryID: categoryID,
	}

	negativeAmountExpense := ExpenseRequest{
		Title:      "coffee",
		Amount:     -1,
		CategoryID: categoryID,
	}

	tests := []struct {
		name    string
		expense ExpenseRequest
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
			name: "empty category id",
			expense: ExpenseRequest{
				Title:  "coffee",
				Amount: 500,
			},
			wantErr: ErrCategoryIDRequired,
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
