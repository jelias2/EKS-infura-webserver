package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"time"

	"github.com/gorilla/mux"
	"github.com/jelias2/infra-test/apis"
	log "github.com/sirupsen/logrus"
)

// Healthcheck will display test response to make sure the server is running
func Healthcheck(w http.ResponseWriter, r *http.Request) {

	log.Info("Healthcheck ")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis.Healthcheck{
		Status:   200,
		Message:  "Healthcheck response",
		Datetime: time.Now().String(),
	})
}

// Get all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	// Hardcoded data - @todo: add database

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis.Book{})
}

// Get single book
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params

	var books []apis.Book
	// Loop through books and find one with the id from the params
	books = append(books, apis.Book{ID: "1", Isbn: "438227", Title: "apis.Book One", Author: &apis.Author{Firstname: "John", Lastname: "Doe"}})
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&apis.Book{})
}

// Add new book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book apis.Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe

	var books []apis.Book
	// Loop through books and find one with the id from the params
	books = append(books, apis.Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &apis.Author{Firstname: "John", Lastname: "Doe"}})
	json.NewEncoder(w).Encode(book)
}

// Update book
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var books []apis.Book
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book apis.Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete book
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var books []apis.Book
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}
