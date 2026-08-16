package model

import (
	"strings"
)

func (e *Expense) Normalize() {
	trimmedTitle := strings.TrimSpace(e.Title)
	trimmedCategory := strings.TrimSpace(e.Category)
	e.Title = trimmedTitle
	e.Category = trimmedCategory
}

func (e Expense) Validate() error {
	if err := e.validateTitle(); err != nil {
		return err
	}

	if err := e.validateAmount(); err != nil {
		return err
	}

	if err := e.validateCategory(); err != nil {
		return err
	}

	return nil
}

func (e Expense) validateTitle() error {
	if e.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (e Expense) validateAmount() error {
	if e.Amount <= 0 {
		return ErrAmountMustBePositive
	}
	return nil
}

func (e Expense) validateCategory() error {
	if e.Category == "" {
		return ErrCategoryRequired
	}
	return nil
}
