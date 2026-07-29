package store

import "github.com/yabu1121/expense-api/internal/model"

var Expenses = []model.Expense{
	{
		ID:       1,
		Title:    "pepper lunch",
		Amount:   900,
		Category: "lunch",
	},
	{
		ID:       2,
		Title:    "cafe latte",
		Amount:   500,
		Category: "coffee"},
	{
		ID:       3,
		Title:    "ts book",
		Amount:   700,
		Category: "book",
	},
}
