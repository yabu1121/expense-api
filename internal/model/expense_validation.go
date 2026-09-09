package model

import (
	"strings"
	"uuid"
)

func (e *ExpenseRequest) Normalize() {
	trimmedTitle := strings.TrimSpace(e.Title)
	e.Title = trimmedTitle
}

func (e ExpenseRequest) Validate() error {
	if err := e.validateTitle(); err != nil {
		return err
	}

	if err := e.validateAmount(); err != nil {
		return err
	}

	if e.CategoryID == uuid.Nil() {
		return ErrCategoryIDRequired
	}

	return nil
}

func (e ExpenseRequest) validateTitle() error {
	if e.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (e ExpenseRequest) validateAmount() error {
	if e.Amount <= 0 {
		return ErrAmountMustBePositive
	}
	return nil
}
