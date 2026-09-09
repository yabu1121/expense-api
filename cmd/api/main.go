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
	tagHandler := handler.NewTagHandler(store)
	expenseTagHandler := handler.NewExpenseTagHandler(store)

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
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetCategoryByID)

	mux.HandleFunc("GET /tags", tagHandler.ListTags)
	mux.HandleFunc("POST /tags", tagHandler.CreateTag)

	mux.HandleFunc("POST /expenses/{expenseID}/tags/{tagID}", expenseTagHandler.AddTagToExpense)

	mux.Handle("GET /sandbox/", http.StripPrefix("/sandbox/", http.FileServer(http.Dir("sandbox"))))
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/sandbox/", http.StatusTemporaryRedirect)
	})

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
