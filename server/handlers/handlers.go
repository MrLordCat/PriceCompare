package handlers

import (
	"amazonPriceGet/server/amazon"
	"amazonPriceGet/server/db"
	"amazonPriceGet/server/fb"
	"amazonPriceGet/server/hv"
	"amazonPriceGet/server/utils"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	tmpl, err := template.ParseFiles("static/templates/template.html", "static/templates/head.html", "static/templates/forms.html", "static/templates/table.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "Protsessorid"
	}

	products, err := db.GetProductsByCategoryFromDB(database, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := db.GetAllCategories(database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Products   []db.Product
		Categories []string
		Selected   string
	}{
		Products:   products,
		Categories: categories,
		Selected:   category,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func FetchHvHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	hvUrl := r.FormValue("hvUrl")
	if hvUrl == "" {
		http.Error(w, "HV.ee URL is required", http.StatusBadRequest)
		return
	}

	_, err := hv.GetHvProducts(database, hvUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
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
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/?category=%s", category), http.StatusSeeOther)
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

func AddCustomProductHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	linkHV := r.FormValue("linkHV")
	if linkHV == "" {
		http.Error(w, "HV.ee link is required", http.StatusBadRequest)
		return
	}

	category := r.FormValue("category")
	if category == "" {
		category = "Custom"
	}

	product, err := hv.FetchCustomProductDetails(linkHV)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch product details: %v", err), http.StatusInternalServerError)
		return
	}

	product.Category = sql.NullString{String: category, Valid: true}

	err = db.InsertOrUpdateProduct(database, product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add product: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/?category=%s", category), http.StatusSeeOther)
}
func UpdateFBLinkHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
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

	linkFB := r.FormValue("FBLink")
	var linkFBValue sql.NullString
	fmt.Println(linkFB)
	if linkFB == "" {
		linkFBValue = sql.NullString{String: "", Valid: false}
	} else {
		linkFBValue = sql.NullString{String: linkFB, Valid: true}
	}

	// Выполнение обновления
	result, err := database.Exec("UPDATE products SET fb_link = ? WHERE title = ?", linkFBValue, title)
	if err != nil {
		fmt.Println("Error updating database:", err, result)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/?category=%s", category), http.StatusSeeOther)
}
func UpdateFBPricesHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
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

	for _, product := range products {
		if product.FBLink.Valid {
			productName := utils.CleanProductName(product.Title) // Очистка имени продукта
			fbPrice, activeStatus, err := fb.GetFBPrice("https://www.facebook.com/marketplace/you/selling", productName)
			if err != nil {
				log.Printf("Error fetching FB price for %s: %v", product.Title, err)
				continue
			}
			product.FBPrice = sql.NullInt64{Int64: fbPrice, Valid: true}
			product.Active = activeStatus
			err = db.UpdateProduct(database, product)
			if err != nil {
				log.Printf("Error updating product %s: %v", product.Title, err)
			}
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/?category=%s", category), http.StatusSeeOther)
}
