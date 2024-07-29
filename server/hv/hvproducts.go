package hv

import (
	"amazonPriceGet/server/db"
	"amazonPriceGet/server/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetHvProducts(database *sql.DB, url string) ([]db.Product, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var products []db.Product

	// Извлечение названия категории
	category := doc.Find("div.header.svelte-uvmab2 h1.svelte-uvmab2").Text()
	category = strings.TrimSpace(category)

	// Добавление категории в базу данных, если она еще не была записана
	err = db.AddCategoryIfNotExists(database, category, url)
	if err != nil {
		return nil, fmt.Errorf("error adding category to database: %v", err)
	}

	doc.Find("tr.svelte-1gwx8vp").Each(func(i int, s *goquery.Selection) {
		// Извлечение названия продукта
		title := s.Find("a.product-name.subtitle-main.svelte-1gwx8vp").Text()
		title = strings.TrimSpace(title)

		// Извлечение ссылки на продукт
		link, exists := s.Find("a.product-name.subtitle-main.svelte-1gwx8vp").Attr("href")
		if !exists {
			log.Printf("URL not found for product: %s", title)
			return
		}
		link = "https://www.hinnavaatlus.ee" + link

		// Извлечение количества предложений
		offersStr := s.Find("td.offers-cell.svelte-1gwx8vp").Text()
		offersStr = strings.TrimSpace(offersStr)
		offers := 0
		if len(offersStr) > 0 {
			offers = utils.IntFromStr(offersStr)
		}

		// Извлечение цены
		priceStr := s.Find("td.price-cell.svelte-1gwx8vp .price.svelte-1gwx8vp").Text()
		priceStr = strings.TrimSpace(priceStr)
		price := utils.Int64FromPriceStr(priceStr)
		product := db.Product{
			Title:    title,
			Offers:   offers,
			Price:    price,
			LinkHV:   link,
			Category: sql.NullString{String: category, Valid: true},
		}

		// Сохранение продукта в базе данных
		err := db.InsertOrUpdateProduct(database, product)
		if err != nil {
			log.Printf("Error inserting or updating product %s: %v", title, err)
		}

		products = append(products, product)
	})

	return products, nil
}

func FetchCustomProductDetails(link string) (db.Product, error) {
	res, err := http.Get(link)
	if err != nil {
		return db.Product{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return db.Product{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return db.Product{}, err
	}

	title := doc.Find("title").Text()
	title = strings.TrimSpace(strings.Split(title, "|")[0])

	offerText := doc.Find("div.data.svelte-1p4umvb button.btn.svelte-1h48h55 span.svelte-1h48h55").First().Text()
	offerText = strings.TrimSpace(offerText)
	offerParts := strings.Fields(offerText)
	if len(offerParts) < 4 {
		return db.Product{}, fmt.Errorf("unexpected format for offer and price: %s", offerText)
	}
	offers := utils.IntFromStr(offerParts[0])
	priceStr := offerParts[len(offerParts)-2] + " " + offerParts[len(offerParts)-1]
	price := utils.Int64FromPriceStr(priceStr)

	return db.Product{
		Title:  title,
		Offers: offers,
		Price:  price,
		LinkHV: link,
	}, nil
}
