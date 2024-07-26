package fb

import (
	"amazonPriceGet/server/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
)

var isLoggedIn = false

func GetFBPrice(url, productName string) (int64, error) {
	// Проверка существования файла ok_fb_page.html и его возраста
	const maxFileAge = 15 * time.Hour
	fileInfo, err := os.Stat("Facebook.html")
	if err == nil {
		if time.Since(fileInfo.ModTime()) < maxFileAge {
			html, err := os.ReadFile("Facebook.html")
			if err != nil {
				return 0, fmt.Errorf("error reading cached page: %v", err)
			}
			return parsePriceFromHTML(string(html), productName)
		}
	}

	var browser *rod.Browser
	if !isLoggedIn {
		browser, err = EnsureLoggedIn()
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		isLoggedIn = true
		defer browser.MustClose()
	}

	page := browser.MustPage(url)
	utils.WaitWithCountdown(5)
	count := 0
	for i := 0; i < 5; i++ {
		_, err := page.Eval(`() => { window.scrollBy(0, window.innerHeight); }`)
		if err != nil {
			return 0, fmt.Errorf("error scrolling page: %v", err)
		}

		count++
		utils.WaitWithCountdown(1)
		fmt.Println("Loading 5/", count)
	}
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	html, err := page.HTML()
	if err != nil {
		return 0, fmt.Errorf("error getting page HTML: %v", err)
	}
	utils.WaitWithCountdown(5)

	// Сохраняем страницу на случай неудачи
	err = utils.SavePage1Content("Facebook.html", html)
	if err != nil {
		return 0, fmt.Errorf("error saving page HTML: %v", err)
	}

	return parsePriceFromHTML(string(html), productName)
}

func parsePriceFromHTML(html, productName string) (int64, error) {
	// Парсинг HTML-контента с помощью goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0, fmt.Errorf("error parsing HTML: %v", err)
	}

	var priceStr string
	var foundProduct bool
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), productName) {
			foundProduct = true
			// Поиск цены в дочерних элементах <span>
			s.Find("span").Each(func(j int, span *goquery.Selection) {
				text := strings.TrimSpace(span.Text())
				if strings.Contains(text, "€") {
					priceStr = text
					return
				}
			})
			if priceStr != "" {
				return
			}
		}
	})

	if !foundProduct {
		return 0, fmt.Errorf("product not found: %s", productName)
	}

	if priceStr == "" {
		return 0, fmt.Errorf("price not found for product: %s", productName)
	}

	fmt.Println("FB PRICE (cached):", priceStr)
	price := utils.Int64FromPriceStr(priceStr)

	return price, nil
}
