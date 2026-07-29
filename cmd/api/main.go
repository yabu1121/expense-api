package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yabu1121/expense-api/internal/handler"
	"github.com/yabu1121/expense-api/internal/store"
)

func main() {
	fmt.Println("server is running on port 8080")

	if err := store.InitDB(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", handler.HealthHandler)
	http.HandleFunc("/version", handler.VersionHandler)

	http.HandleFunc("/expenses", handler.ExpensesHandler)
	http.HandleFunc("/expenses/{id}", handler.ExpensesHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
