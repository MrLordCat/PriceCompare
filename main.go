package main

import (
	"amazonPriceGet/server/db"
	"amazonPriceGet/server/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := sql.Open("sqlite3", "server/db/products.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	// Создаем таблицу, если она не существует
	db.CreateTable(database)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllProducts(w, r, database)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateAmazonLinkHandler(w, r, database)
	})
	http.HandleFunc("/fetch-hv", func(w http.ResponseWriter, r *http.Request) {
		handlers.FetchHvHandler(w, r, database)
	})
	http.HandleFunc("/update-amazon-prices", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateAmazonPricesHandler(w, r, database)
	})
	http.HandleFunc("/add-custom", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddCustomProductHandler(w, r, database)
	})
	http.HandleFunc("/update-fb", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateFBLinkHandler(w, r, database)
	})
	http.HandleFunc("/update-fb-prices", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateFBPricesHandler(w, r, database)
	})
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
