package fb

import (
	"amazonPriceGet/server/utils"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
)

var isLoggedIn = false

func GetFBPrice(url, productName string) (int64, string, error) {
	// Проверка существования файла ok_fb_page.html и его возраста
	const maxFileAge = 15 * time.Hour
	fileInfo, err := os.Stat("Facebook.html")
	if err == nil {
		if time.Since(fileInfo.ModTime()) < maxFileAge {
			html, err := os.ReadFile("Facebook.html")
			if err != nil {
				return 0, "", fmt.Errorf("error reading cached page: %v", err)
			}
			return parsePriceFromHTML(string(html), productName)
		}
	}

	var browser *rod.Browser
	if !isLoggedIn {
		browser, err = EnsureLoggedIn()
		if err != nil {
			fmt.Println(err)
			return 0, "", err
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
			return 0, "", fmt.Errorf("error scrolling page: %v", err)
		}

		count++
		utils.WaitWithCountdown(1)
		fmt.Println("Loading 5/", count)
	}
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	html, err := page.HTML()
	if err != nil {
		return 0, "", fmt.Errorf("error getting page HTML: %v", err)
	}
	utils.WaitWithCountdown(5)

	// Сохраняем страницу на случай неудачи
	err = utils.SavePage1Content("Facebook.html", html)
	if err != nil {
		return 0, "", fmt.Errorf("error saving page HTML: %v", err)
	}

	return parsePriceFromHTML(string(html), productName)
}

func parsePriceFromHTML(html, productName string) (int64, string, error) {
	// Парсинг HTML-контента с помощью goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return 0, "", fmt.Errorf("error parsing HTML: %v", err)
	}
	var activeStatus string
	var priceStr string
	var foundProduct bool
	var prices []string
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), productName) {
			foundProduct = true
			blockHtml, err := getBlockHtml(s)
			if err != nil {
				fmt.Printf("Error getting block HTML: %v\n", err)
			}
			blockDoc, err := goquery.NewDocumentFromReader(strings.NewReader(blockHtml))
			if err != nil {
				return
			}

			blockDoc.Find("span").Each(func(j int, span *goquery.Selection) {
				text := strings.TrimSpace(span.Text())
				if strings.Contains(text, "€") {
					// Найти первую цену и обрезать все после второго символа "€"
					euroIndex := strings.Index(text, "€")
					secondEuroIndex := strings.Index(text[euroIndex+1:], "€")
					if secondEuroIndex != -1 {
						// Увеличить индекс второго "€", чтобы учитывать начальную часть строки
						secondEuroIndex += euroIndex + 1
						text = text[:secondEuroIndex]
					}
					prices = append(prices, text)
				}
				if strings.Contains(text, "Active") {
					activeStatus = "Active"
				}
			})
		}
	})
	priceStr = prices[0]
	if !foundProduct {
		return 0, "", fmt.Errorf("product not found: %s", productName)
	}

	if priceStr == "" {
		return 0, "", fmt.Errorf("price not found for product: %s", productName)
	}
	if activeStatus == "" {
		activeStatus = "Inactive"
	}

	//fmt.Println("FB PRICE (", productName, "):", priceStr)
	//	fmt.Println("Status:", activeStatus)
	price := utils.Int64FromPriceStr(priceStr)

	return price, activeStatus, nil
}

var fileCounter int

func init() {
	fileCounter = 0
}

func getBlockHtml(s *goquery.Selection) (string, error) {
	// Получаем внешний HTML текущего блока
	html, err := s.Html()
	if err != nil {
		return "", err
	}

	// Добавляем больше родительских блоков, чтобы захватить больший размер блока
	parent := s.Parent()
	for i := 0; i < 3 && parent.Length() > 0; i++ { // Захватываем до 3 уровней вверх
		parentHtml, err := parent.Html()
		if err != nil {
			return "", err
		}
		html = "<div>" + parentHtml + "</div>"
		parent = parent.Parent()
	}

	// Проверка условий и сохранение HTML-контента в файл
	if strings.HasPrefix(html, `<div><div aria-label="`) && !strings.HasPrefix(html, `<div><div aria-label="Marketplace`) {
		fileCounter++
		filename := filepath.Join("savedhtmls", fmt.Sprintf("block_%d.html", fileCounter))
		err = os.WriteFile(filename, []byte(html), 0644)
		if err != nil {
			return "", err
		}
		return html, nil
	}

	// Если условие не выполнено, возвращаем пустую строку
	return "", nil
}
