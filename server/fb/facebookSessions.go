package fb

import (
	"amazonPriceGet/server/utils"
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func EnsureLoggedIn() (*rod.Browser, error) {
	// Указание пути к исполняемому файлу браузера и использование постоянного профиля
	u := launcher.New().
		Set("load-extension", "/home/roppi/amazonPriceGet/SingleFile-master").
		Set("user-data-dir", "/path/to/your/chrome/profile").
		Headless(false).
		MustLaunch()

	// Создание нового браузера с указанием пути
	browser := rod.New().ControlURL(u).MustConnect()

	// Проверка успешного входа
	err := checkLoginSuccess(browser)
	if err != nil {
		fmt.Println("Login check failed, performing manual login")

		// Выполнение ручного входа
		err = performManualLogin(browser)
		if err != nil {
			return nil, fmt.Errorf("error performing manual login: %v", err)
		}

		// Повторная проверка успешного входа
		err = checkLoginSuccess(browser)
		if err != nil {
			return nil, fmt.Errorf("login check failed after manual login: %v", err)
		}
	}

	return browser, nil
}
func checkLoginSuccess(browser *rod.Browser) error {
	page := browser.MustPage("https://www.facebook.com/")
	html, err := page.HTML()
	if err != nil {
		utils.SavePage1Content("check_login_error.html", html)
		return fmt.Errorf("error reading check response: %v", err)
	}

	// Сохранение содержимого body для проверки
	err = utils.SavePage1Content("check_login_body.html", html)
	if err != nil {
		return fmt.Errorf("error saving body content: %v", err)
	}
	utils.WaitWithCountdown(10)
	if !strings.Contains(html, "Chris Vlassenko") {
		utils.SavePage1Content("login_check_failed.html", html)
		return fmt.Errorf("login check failed, 'Chris Vlassenko' not found")
	}

	fmt.Println("Login successful, 'Chris Vlassenko' found")
	return nil
}

func performManualLogin(browser *rod.Browser) error {
	// Открытие страницы входа
	page := browser.MustPage("https://www.facebook.com/login.php")

	// Ожидание ручного ввода логина и пароля
	fmt.Println("Введите логин и пароль вручную, затем нажмите Enter в консоли, чтобы продолжить...")
	fmt.Scanln()

	// Сохранение HTML-контента страницы
	html, err := page.HTML()
	if err != nil {
		return fmt.Errorf("error getting page HTML: %v", err)
	}
	err = utils.SavePage1Content("post_login_page.html", html)
	if err != nil {
		return fmt.Errorf("error saving page HTML: %v", err)
	}

	return nil
}
