package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var books []Book

// コントローラー
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get prams from url
	// Loop through books and find with id
	for _, book := range books {

		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func crateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)  // 戻り値は使わない。Bookインスタンスに送られてきたデータをはめ込む
	book.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID ダブる可能性が一応ある
	books = append(books, book)                // ここでsave
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book) // 戻り値は使わない。Bookインスタンスに送られてきたデータをはめ込む
			book.ID = params["id"]                    // Mock ID ダブる可能性が一応ある
			books = append(books, book)               // ここでsave
			json.NewEncoder(w).Encode(book)
		}
	}

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
}

func createMockBook() []Book {
	books = append(books,
		Book{
			ID:     "1",
			Isbn:   "44326",
			Title:  "book one",
			Author: &Author{FirstName: "jon", LastName: "deo"},
		},
	)

	books = append(books,
		Book{
			ID:     "2",
			Isbn:   "44332",
			Title:  "book two",
			Author: &Author{FirstName: "yam", LastName: "momo"},
		},
	)
	return books
}

func main() {
	books = createMockBook()
	// 初期化処理
	r := mux.NewRouter()

	// ルーティング
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", crateBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
