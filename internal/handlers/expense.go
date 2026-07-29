package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/yabu1121/expense-api/internal/model"
)

var expenses = []model.Expense{
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

func expenseHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(expenses)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		log.Println(err)
		return
	}
}

func expensePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nextID := 1

	if len(expenses) > 0 {
		nextID = expenses[len(expenses)-1].ID + 1
	}

	var expense model.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	expense.ID = nextID
	expenses = append(expenses, expense)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(expense); err != nil {
		log.Println(err)
		return
	}

}

func getExpenseById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		return
	}
	found := false

	var expense model.Expense
	for _, e := range expenses {
		if e.ID == id {
			expense = e
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "expense not found", http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(expense); err != nil {
		log.Println(err)
		return
	}
}

func deleteExpenseById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		return
	}
	found := false
	var targetIndex int

	for i, e := range expenses {
		if e.ID == id {
			found = true
			targetIndex = i
			break
		}
	}
	if !found {
		http.Error(w, "expense not found", http.StatusNotFound)
		return
	}
	expenses = append(
		expenses[:targetIndex],
		expenses[targetIndex + 1:]...
	)
	w.WriteHeader(http.StatusNoContent)
}

func updateExpenseById (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}
	var updatedExpense model.Expense
	if err := json.NewDecoder(r.Body).Decode(&updatedExpense); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for i, expense := range expenses {
        if expense.ID == id {
            updatedExpense.ID = id
            expenses[i] = updatedExpense
            if err := json.NewEncoder(w).Encode(updatedExpense); err != nil {
                log.Println(err)
            }
            return
        }
    }
    http.Error(w, "expense not found", http.StatusNotFound)
}

func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		if id != "" {
			getExpenseById(w, r)
			return
		}
		expenseHandler(w, r)
		return
	case http.MethodPost:
		if id != "" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		expensePostHandler(w, r)
		return
	case http.MethodDelete:
		if id != "" {
			deleteExpenseById(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	case http.MethodPut:
		if id != "" {
			updateExpenseById(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
