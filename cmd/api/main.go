package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yabu1121/expense-api/internal/handler"
	"github.com/yabu1121/expense-api/internal/store"
)

func main() {
	expenseStore, err := store.NewSQLiteStore()
	if err != nil {
		log.Fatal(err)
	}
	defer expenseStore.Close()

	expenseHandler := handler.NewExpenseHandler(expenseStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/version", handler.VersionHandler)

	mux.Handle("/expenses", expenseHandler)
	mux.Handle("/expenses/{id}", expenseHandler)

	fmt.Println("server is running on port 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
