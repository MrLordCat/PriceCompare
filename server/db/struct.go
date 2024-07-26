package db

import "database/sql"

type Product struct {
	ID           int
	Title        string
	Offers       int
	Price        int64
	LinkHV       string
	LinkAmazon   sql.NullString
	PriceAmazon  sql.NullInt64
	PriceDiff    sql.NullInt64
	DeliveryTime sql.NullString
	Category     sql.NullString
	Used         sql.NullString
	PriceMinus15 sql.NullInt64
	FBPrice      sql.NullInt64
	FBLink       sql.NullString
	Active       string
}
