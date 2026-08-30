package model

import "errors"

var ErrNameRequired = errors.New("name must not be empty")
var ErrCategoryAlreadyExists = errors.New("this category name already exists")
var ErrCategoryNotFound = errors.New("category is not found")
