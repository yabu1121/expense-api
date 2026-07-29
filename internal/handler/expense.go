package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/yabu1121/expense-api/internal/model"
	"github.com/yabu1121/expense-api/internal/store"
)

func expenseHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(store.Expenses)
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

	if len(store.Expenses) > 0 {
		nextID = store.Expenses[len(store.Expenses)-1].ID + 1
	}

	var expense model.Expense

	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	expense.ID = nextID
	store.Expenses = append(store.Expenses, expense)
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
	for _, e := range store.Expenses {
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

	for i, e := range store.Expenses {
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
	store.Expenses = append(
		store.Expenses[:targetIndex],
		store.Expenses[targetIndex + 1:]...
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

	for i, expense := range store.Expenses {
        if expense.ID == id {
            updatedExpense.ID = id
            store.Expenses[i] = updatedExpense
            if err := json.NewEncoder(w).Encode(updatedExpense); err != nil {
                log.Println(err)
            }
            return
        }
    }
    http.Error(w, "expense not found", http.StatusNotFound)
}

func store.ExpensesHandler(w http.ResponseWriter, r *http.Request) {
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
