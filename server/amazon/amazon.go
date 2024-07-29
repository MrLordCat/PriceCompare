package amazon

import (
	"amazonPriceGet/server/utils"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"golang.org/x/exp/rand"
)

var (
	userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4894.117 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4855.118 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4892.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4854.191 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4859.153 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36/null",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36,gzip(gfe)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4895.86 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_3_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4860.89 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4885.173 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4864.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4877.207 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML%2C like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.133 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4872.118 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_3_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4876.128 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML%2C like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	}
	randSource = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
)
var captchaStatus = false

func GetAmazonDetails(url string) (int64, string, sql.NullString, error) {
	if captchaStatus {
		return GetAmazonDetailsIncognito(url)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", sql.NullString{}, err
	}

	userAgent := userAgents[randSource.Intn(len(userAgents))]
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return 0, "", sql.NullString{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, "", sql.NullString{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, "", sql.NullString{}, err
	}

	if strings.Contains(strings.ToLower(doc.Text()), "captcha") {
		fmt.Println("Captcha found on the page. Trying incognito.")
		captchaStatus = true
		return GetAmazonDetailsIncognito(url)
	}
	if strings.Contains(strings.ToLower(doc.Text()), "no featured offers available") || strings.Contains(strings.ToLower(doc.Text()), "keine hervorgehobenen angebote verfügbar") {
		fmt.Println("Product is not available.")
		return 0, "N/A", sql.NullString{String: "No", Valid: true}, nil
	}

	// Извлечение цены
	priceWhole := doc.Find("div.a-section.a-spacing-none.aok-align-center.aok-relative span.a-price-whole").First().Text()
	priceFraction := doc.Find("div.a-section.a-spacing-none.aok-align-center.aok-relative span.a-price-fraction").First().Text()

	if priceWhole == "" || priceFraction == "" {
		return 0, "", sql.NullString{}, fmt.Errorf("could not find price on page")
	}

	priceStr := strings.TrimSpace(priceWhole) + "." + strings.TrimSpace(priceFraction)
	price := utils.Int64FromPriceStr(priceStr)

	// Извлечение даты доставки
	deliveryDate := doc.Find("#mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE .a-text-bold").Text()
	deliveryDate = strings.TrimSpace(deliveryDate)
	var used sql.NullString
	if doc.Find("div#usedBuySection span.a-text-bold:contains('Buy used')").Length() > 0 {
		used = sql.NullString{String: "Yes", Valid: true}
	} else {
		used = sql.NullString{String: "No", Valid: true}
	}

	return price, deliveryDate, used, nil
}
func GetAmazonDetailsIncognito(url string) (int64, string, sql.NullString, error) {
	// Настройка контекста для инкогнито-режима
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("incognito", true),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Навигация и извлечение данных
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		//chromedp.Sleep(1*time.Second),
		chromedp.OuterHTML("html", &res),
	)
	if err != nil {
		return 0, "", sql.NullString{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		return 0, "", sql.NullString{}, err
	}

	// Проверка на отсутствие предложений на странице
	if strings.Contains(strings.ToLower(doc.Text()), "no featured offers available") || strings.Contains(strings.ToLower(doc.Text()), "keine hervorgehobenen angebote verfügbar") {
		fmt.Println("Product is not available.")
		return 0, "", sql.NullString{}, err
	}

	// Извлечение цены
	priceWhole := doc.Find("div.a-section.a-spacing-none.aok-align-center.aok-relative span.a-price-whole").First().Text()

	if priceWhole == "" {
		return 0, "", sql.NullString{}, fmt.Errorf("could not find price on page")
	}

	priceStr := strings.TrimSpace(priceWhole) + "."
	price := utils.Int64FromPriceStr(priceStr)

	// Извлечение даты доставки
	deliveryDate := doc.Find("#mir-layout-DELIVERY_BLOCK-slot-PRIMARY_DELIVERY_MESSAGE_LARGE .a-text-bold").Text()
	deliveryDate = strings.TrimSpace(deliveryDate)
	var used sql.NullString
	if doc.Find("div#usedBuySection span.a-text-bold:contains('Buy used')").Length() > 0 {
		used = sql.NullString{String: "Yes", Valid: true}
	} else {
		used = sql.NullString{String: "No", Valid: true}
	}
	return price, deliveryDate, used, nil
}
