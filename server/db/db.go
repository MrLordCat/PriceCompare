package db

import (
	"database/sql"
	"log"
)

func GetAllCategories(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT category FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
func GetProductsByCategoryFromDB(db *sql.DB, category string) ([]Product, error) {
	rows, err := db.Query("SELECT id, title, category, offers, price, link_hv, link_amazon, price_amazon, price_diff, delivery_time, used, price_minus_15, fb_price, fb_link, active  FROM products WHERE category = ?", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Title, &product.Category, &product.Offers, &product.Price, &product.LinkHV, &product.LinkAmazon, &product.PriceAmazon, &product.PriceDiff, &product.DeliveryTime, &product.Used, &product.PriceMinus15, &product.FBPrice, &product.FBLink, &product.Active)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return products, nil
}
func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS products (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"offers" INTEGER,
		"price" INTEGER,
		"link_hv" TEXT UNIQUE,
		"link_amazon" TEXT,
		"price_amazon" INTEGER,
		"price_diff" INTEGER,
		"delivery_time" TEXT,
		"category" TEXT DEFAULT '',
		"used" TEXT DEFAULT 'No',
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"price_minus_15" INTEGER,
		"fb_price" INTEGER,
		"fb_link" TEXT,
		"active" TEXT DEFAULT 'N/A'
	);`
	createCategoriesTableSQL := `CREATE TABLE IF NOT EXISTS categories (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"category" TEXT NOT NULL UNIQUE,
		"url" TEXT NOT NULL
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)

	}
	_, err = db.Exec(createCategoriesTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
func AddCategoryIfNotExists(db *sql.DB, category, url string) error {
	query := `INSERT INTO categories (category, url) VALUES (?, ?) ON CONFLICT(category) DO NOTHING;`
	_, err := db.Exec(query, category, url)
	return err
}
func UpdateProduct(db *sql.DB, product Product) error {
	query := `
		UPDATE products
		SET price_amazon = ?, price_diff = ?, delivery_time = ?, used = ?, price_minus_15 = ?, fb_price = ?, fb_link = ?, active = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, product.PriceAmazon, product.PriceDiff, product.DeliveryTime, product.Used, product.PriceMinus15, product.FBPrice, product.FBLink, product.Active, product.ID)
	return err
}
func InsertOrUpdateProduct(db *sql.DB, product Product) error {
	query := `
		INSERT INTO products (title, offers, price, link_hv, category)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(link_hv) DO UPDATE SET
			title = excluded.title,
			offers = excluded.offers,
			price = excluded.price,
			category = excluded.category`
	_, err := db.Exec(query, product.Title, product.Offers, product.Price, product.LinkHV, product.Category)
	return err
}
func GetCategoryUrl(db *sql.DB, category string) (string, error) {
	var url string
	query := `SELECT url FROM categories WHERE category = ?`
	err := db.QueryRow(query, category).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
