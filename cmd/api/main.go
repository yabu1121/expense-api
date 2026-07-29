package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yabu1121/expense-api/internal/handlers"
)

func main() {
	fmt.Println("server is running on port 8080")

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/version", handlers.VersionHandler)

	http.HandleFunc("/expenses", handlers.ExpensesHandler)
	http.HandleFunc("/expenses/{id}", handlers.ExpensesHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
