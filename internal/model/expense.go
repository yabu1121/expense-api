package model

type Expense struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Amount   int    `json:"amount"`
	Category string `json:"category"`
}
type ExpenseSummary struct {
	Count       int `json:"count"`
	TotalAmount int `json:"total_amount"`
}
