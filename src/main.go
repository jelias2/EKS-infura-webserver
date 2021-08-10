package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jelias2/infra-test/handlers"
	log "github.com/sirupsen/logrus"
)

// Main function
func main() {
	log.SetReportCaller(true)
	log.Info("Beginning Webserver main.go...")
	// Init router
	log.Info("Creating mux router and initalizing mux router")
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/health", handlers.Healthcheck).Methods("GET")
	r.HandleFunc("/", handlers.Healthcheck).Methods("GET")

	r.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	r.HandleFunc("/books", handlers.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	// Start server
	log.Info("Beginning to server traffic on port")
	log.Fatal(http.ListenAndServe(":8000", r))
}
