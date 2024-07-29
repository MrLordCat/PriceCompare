package handlers

import (
	"amazonPriceGet/server/amazon"
	"amazonPriceGet/server/db"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var progressMap sync.Map

func UpdateAmazonPricesHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	category := r.FormValue("category")
	if category == "" {
		http.Error(w, "Category not specified", http.StatusBadRequest)
		return
	}

	products, err := db.GetProductsByCategoryFromDB(database, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalProducts := len(products)
	progressKey := fmt.Sprintf("progress-%s", category)
	progressMap.Store(progressKey, 0)

	go func() {
		updatedCount := 0

		for _, product := range products {
			if product.LinkAmazon.Valid {
				priceAmazon, deliveryTime, used, err := amazon.GetAmazonDetails(product.LinkAmazon.String)
				if err != nil {
					log.Printf("Error fetching Amazon details for %s: %v", product.Title, err)
					continue
				}
				product.PriceAmazon = sql.NullInt64{Int64: priceAmazon, Valid: true}
				product.DeliveryTime = sql.NullString{String: deliveryTime, Valid: true}
				product.PriceDiff = sql.NullInt64{Int64: product.Price - priceAmazon, Valid: true}
				product.Used = used
				err = db.UpdateProduct(database, product)
				if err != nil {
					log.Printf("Error updating product %s: %v", product.Title, err)
				}
				updatedCount++
			}
			progress := (updatedCount * 100) / totalProducts
			progressMap.Store(progressKey, progress)
		}
		time.Sleep(1 * time.Second)
		progressMap.Delete(progressKey)
	}()
	http.Redirect(w, r, fmt.Sprintf("/?category=%s", category), http.StatusSeeOther)
}

func GetProgressHandler(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("category")
	if category == "" {
		http.Error(w, "Category not specified", http.StatusBadRequest)
		return
	}

	progressKey := fmt.Sprintf("progress-%s", category)
	progress, ok := progressMap.Load(progressKey)
	if !ok {
		progress = 100
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"progress": %d}`, progress)
}

func UpdateAmazonLinkHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	category := r.FormValue("category")
	if category == "" {
		http.Error(w, "Category not specified", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Invalid product title", http.StatusBadRequest)
		return
	}

	linkAmazon := r.FormValue("linkAmazon")
	var linkAmazonValue sql.NullString

	if linkAmazon == "" {
		linkAmazonValue = sql.NullString{String: "", Valid: false}
	} else {
		linkAmazonValue = sql.NullString{String: linkAmazon, Valid: true}
	}

	// Выполнение обновления
	result, err := database.Exec("UPDATE products SET link_amazon = ? WHERE title = ?", linkAmazonValue, title)
	if err != nil {
		fmt.Println("Error updating database:", err, result)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/?category=%s", category), http.StatusSeeOther)
}
