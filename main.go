package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pump_finder/helper"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProcessTokenList(input string) []string {
	// Pecah berdasarkan newlines (\n atau \r\n)
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	// Trim setiap baris untuk menghindari whitespace berlebih
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}
func main() {
	// Inisialisasi router Gin
	router := gin.Default()

	// Route untuk halaman utama
	router.GET("/", func(c *gin.Context) {
		c.File("templates/index.html") // Render halaman HTML
	})

	// Route untuk menerima POST di /find
	router.POST("/find", func(c *gin.Context) {
		tokenList := c.PostForm("userInput")
		minBuy := c.PostForm("minbuy")
		feetFloat, _ := strconv.ParseFloat(strings.TrimSpace(minBuy), 64)
		minBuyStr := helper.SolToLamports(feetFloat)
		tokenListArr := ProcessTokenList(tokenList)
		firstResult, err := fetchAndHandleError(tokenListArr[0], minBuyStr)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println()
		for _, token := range tokenListArr[1:] {
			nextResult, err := fetchAndHandleError(token, minBuyStr)
			if err != nil {
				continue
			}
			fmt.Println()
			firstResult = helper.FindMatchingUsers(firstResult, nextResult)
		}
		finalResult := removeDuplicateUsers(firstResult)
		// Konversi hasil ke JSON
		sendJSONResponse(c, finalResult)
	})

	// Jalankan server di port 8080
	router.Run(":8080")
}

func removeDuplicateUsers(results []helper.Result) []helper.Result {
	seen := make(map[string]bool)
	var uniqueResults []helper.Result

	for _, result := range results {
		if !seen[result.UserAddress] {
			seen[result.UserAddress] = true
			uniqueResults = append(uniqueResults, result)
		}
	}

	return uniqueResults
}
func fetchAndHandleError(token, minBuyStr string) ([]helper.Result, error) {
	result, err := helper.MainFetch(token, minBuyStr)
	if err != nil {
		fmt.Println("Error fetching data for token:", token, err)
		return nil, err
	}
	return result, nil
}

// sendJSONResponse: Mengirim respons JSON ke client.
func sendJSONResponse(c *gin.Context, data []helper.Result) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	fmt.Println(string(jsonData))
	c.String(http.StatusOK, string(jsonData))
}