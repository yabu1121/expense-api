package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
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

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("v1.0.0")); err != nil {
		log.Println(err)
		return
	}
}
