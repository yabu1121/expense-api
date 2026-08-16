package model

import (
	"errors"
)

var ErrTitleRequired = errors.New("title must not be empty")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrCategoryRequired = errors.New("category must not be empty")
var ErrExpenseNotFound = errors.New("expense not found")
