package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TextRequest struct {
	Text string `json:"text"`
}

type ClassificationResponse struct {
	Service    string  `json:"service"`
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Merchant   string  `json:"merchant"`
	Status     string  `json:"status"`
}

type Category struct {
	Keywords   []string
	Confidence float64
}

var categories = map[string][]string{
	"Groceries":     {"grocery", "supermarket", "food", "vegetables", "fruits", "milk"},
	"Restaurant":    {"restaurant", "cafe", "dining", "menu", "takeaway"},
	"Transport":     {"uber", "lyft", "taxi", "bus", "train", "fare"},
	"Utilities":     {"electricity", "water", "gas", "internet", "phone"},
	"Entertainment": {"movie", "theatre", "concert", "game", "netflix"},
}

func classifyText(text string) (string, float64) {
	text = strings.ToLower(text)
	maxMatches := 0
	bestCategory := "Other"

	for category, keywords := range categories {
		matches := 0
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				matches++
			}
		}
		if matches > maxMatches {
			maxMatches = matches
			bestCategory = category
		}
	}

	confidence := 0.3
	if maxMatches > 0 {
		confidence = float64(maxMatches) * 0.2
		if confidence > 0.9 {
			confidence = 0.9
		}
	}

	return bestCategory, confidence
}

func extractAmount(text string) float64 {
	// Look for currency patterns like $XX.XX or XX.XX
	re := regexp.MustCompile(`\$?\d+\.\d{2}`)
	matches := re.FindAllString(text, -1)

	if len(matches) > 0 {
		// Convert first match to float
		amountStr := strings.TrimPrefix(matches[0], "$")
		amount, _ := strconv.ParseFloat(amountStr, 64)
		return amount
	}
	return 0
}

func extractDate(text string) string {
	// Look for common date formats
	dateFormats := []string{
		`\d{2}/\d{2}/\d{4}`,
		`\d{2}-\d{2}-\d{4}`,
		`\d{4}-\d{2}-\d{2}`,
	}

	for _, format := range dateFormats {
		re := regexp.MustCompile(format)
		if match := re.FindString(text); match != "" {
			return match
		}
	}
	return ""
}

func main() {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/classify", func(c *gin.Context) {
		var request TextRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category, confidence := classifyText(request.Text)
		amount := extractAmount(request.Text)
		date := extractDate(request.Text)

		response := ClassificationResponse{
			Service:    "classifier",
			Category:   category,
			Confidence: confidence,
			Amount:     amount,
			Date:       date,
			Status:     "success",
		}

		c.JSON(http.StatusOK, response)
	})

	log.Fatal(r.Run(":8001"))
}
