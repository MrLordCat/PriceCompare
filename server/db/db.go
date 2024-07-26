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
	rows, err := db.Query("SELECT id, title, category, offers, price, link_hv, link_amazon, price_amazon, price_diff, delivery_time, used, price_minus_15, fb_price, fb_link  FROM products WHERE category = ?", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Title, &product.Category, &product.Offers, &product.Price, &product.LinkHV, &product.LinkAmazon, &product.PriceAmazon, &product.PriceDiff, &product.DeliveryTime, &product.Used, &product.PriceMinus15, &product.FBPrice, &product.FBLink)
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
		"fb_link" TEXT
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateProduct(db *sql.DB, product Product) error {
	query := `
		UPDATE products
		SET price_amazon = ?, price_diff = ?, delivery_time = ?, used = ?, price_minus_15 = ?, fb_price = ?, fb_link = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, product.PriceAmazon, product.PriceDiff, product.DeliveryTime, product.Used, product.PriceMinus15, product.FBPrice, product.FBLink, product.ID)
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
