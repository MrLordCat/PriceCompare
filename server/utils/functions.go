package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Int64FromPriceStr(priceStr string) int64 {
	priceStr = strings.ReplaceAll(priceStr, "€", "")
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.ReplaceAll(priceStr, ",", ".")
	priceStr = strings.ReplaceAll(priceStr, "..", ".")
	priceFloat, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Printf("Error converting string to float64: %v\n", err)
		return 0
	}
	price := int64(priceFloat)
	return price
}

func SavePage1Content(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func WaitWithCountdown(seconds int) {
	for i := seconds; i > 0; i-- {
		fmt.Printf("\rWaiting for %d seconds...", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("\rWaiting complete!          ")
}

func IntFromStr(s string) int {
	s = strings.TrimSpace(s)
	value, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error converting string to int: %v\n", err)
		return 0
	}
	return value
}

func CleanProductName(name string) string {
	unwantedWords := []string{"Protsessor", "Kõvakettad - SSD", "и т.д."}
	for _, word := range unwantedWords {
		name = strings.ReplaceAll(name, word, "")
	}
	return strings.TrimSpace(name)
}
