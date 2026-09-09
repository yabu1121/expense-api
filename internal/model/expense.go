package model

import (
	"time"
	"uuid"
)

type Expense struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Amount     int       `json:"amount"`
	CategoryID uuid.UUID `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type ExpenseRequest struct {
	Title      string    `json:"title"`
	Amount     int       `json:"amount"`
	CategoryID uuid.UUID `json:"category_id"`
}

type ExpenseSummary struct {
	Count       int `json:"count"`
	TotalAmount int `json:"total_amount"`
}

type ExpenseFilter struct {
	Category string
	Limit    int
	Offset   int
}
