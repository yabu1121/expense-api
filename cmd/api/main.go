package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type Expense struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Amount   int    `json:"amount"`
	Category string `json:"category"`
}

var expenses = []Expense{
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{
		Status: "ok",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := w.Write(jsonData); err != nil {
		log.Println(err)
		return
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("v1.0.0")); err != nil {
		log.Println(err)
		return
	}
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

	var expense Expense

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

	var expense Expense
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
	var updatedExpense Expense
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

func expensesHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	fmt.Println("server is running on port 8080")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/version", versionHandler)

	http.HandleFunc("/expenses", expensesHandler)
	http.HandleFunc("/expenses/{id}", expensesHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
