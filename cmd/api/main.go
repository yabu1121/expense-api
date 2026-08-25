package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yabu1121/expense-api/internal/handler"
	"github.com/yabu1121/expense-api/internal/store"
)

func main() {
	dbPath := os.Getenv("DB_PATH")

	if dbPath == "" {
		dbPath = "expenses.db"
	}

	store, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	expenseHandler := handler.NewExpenseHandler(store)
	expenseSummaryHandler := handler.NewExpenseSummaryHandler(store)
	categoryHandler := handler.NewCategoryHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.HealthHandler)
	mux.HandleFunc("GET /version", handler.VersionHandler)

	mux.HandleFunc("GET /expenses", expenseHandler.ListExpenses)
	mux.HandleFunc("GET /expenses/{id}", expenseHandler.GetExpenseByID)
	mux.HandleFunc("POST /expenses", expenseHandler.CreateExpense)
	mux.HandleFunc("PUT /expenses/{id}", expenseHandler.UpdateExpenseByID)
	mux.HandleFunc("DELETE /expenses/{id}", expenseHandler.DeleteExpenseByID)

	mux.HandleFunc("GET /expenses/summary", expenseSummaryHandler.GetExpenseSummary)

	mux.HandleFunc("GET /categories", categoryHandler.ListCategories)
	mux.HandleFunc("POST /categories", categoryHandler.CreateCategory)


	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		fmt.Println("server is running on port 8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	select {}
}
