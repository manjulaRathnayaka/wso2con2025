package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"bills-service/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/clientcredentials"
)

// AuthConfig holds authentication configuration for external services
type AuthConfig struct {
	ServiceURL     string
	ConsumerKey    string
	ConsumerSecret string
	TokenURL       string
	APIKey         string
}

// LoadAuthConfig loads service authentication configuration from environment variables
func LoadAuthConfig(serviceType string) AuthConfig {
	switch serviceType {
	case "ocr":
		return AuthConfig{
			ServiceURL:     os.Getenv("CHOREO_OCR_SERVICE_CONN_SERVICEURL"),
			ConsumerKey:    os.Getenv("CHOREO_OCR_SERVICE_CONN_CONSUMERKEY"),
			ConsumerSecret: os.Getenv("CHOREO_OCR_SERVICE_CONN_CONSUMERSECRET"),
			TokenURL:       os.Getenv("CHOREO_OCR_SERVICE_CONN_TOKENURL"),
			APIKey:         os.Getenv("CHOREO_OCR_SERVICE_CONN_APIKEY"),
		}
	case "parser":
		return AuthConfig{
			ServiceURL:     os.Getenv("CHOREO_BILL_PARSER_SERVICE_CONN_SERVICEURL"),
			ConsumerKey:    os.Getenv("CHOREO_BILL_PARSER_SERVICE_CONN_CONSUMERKEY"),
			ConsumerSecret: os.Getenv("CHOREO_BILL_PARSER_SERVICE_CONN_CONSUMERSECRET"),
			TokenURL:       os.Getenv("CHOREO_BILL_PARSER_SERVICE_CONN_TOKENURL"),
			APIKey:         os.Getenv("CHOREO_BILL_PARSER_SERVICE_CONN_APIKEY"),
		}
	default:
		return AuthConfig{}
	}
}

// createAuthenticatedClient creates an HTTP client with authentication configured
// based on the provided AuthConfig. It handles OAuth2 if credentials are available.
func createAuthenticatedClient(config AuthConfig) *http.Client {
	var client *http.Client

	// Use OAuth2 if credentials are available
	if config.ConsumerKey != "" && config.ConsumerSecret != "" && config.TokenURL != "" {
		clientCredsConfig := clientcredentials.Config{
			ClientID:     config.ConsumerKey,
			ClientSecret: config.ConsumerSecret,
			TokenURL:     config.TokenURL,
		}
		client = clientCredsConfig.Client(context.Background())
	} else {
		// Otherwise use standard client
		client = &http.Client{}
	}

	// Set timeout
	client.Timeout = 30 * time.Second
	return client
}

// addAuthHeaders adds authentication headers to the request based on the config
func addAuthHeaders(req *http.Request, config AuthConfig) {
	// Add API key authentication if available
	if config.APIKey != "" {
		req.Header.Add("Choreo-API-Key", config.APIKey)
	}
}

// ProcessBillImageV2 is an enhanced implementation that calls external services
// for OCR and bill parsing using OAuth2 authentication
func ProcessBillImageV2(c *gin.Context) (*models.BillProcessResult, error) {
	// Get file from request
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		return nil, errors.New("no image file provided")
	}
	defer file.Close()

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Get content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Step 1: Extract text using OCR service
	extractedText, err := callOCRService(fileBytes, header.Filename)
	if err != nil {
		return nil, fmt.Errorf("OCR processing failed: %v", err)
	}

	if extractedText == "" {
		return nil, errors.New("no text could be extracted from the image")
	}

	// Step 2: Parse bill information from text
	parsedBill, err := parseBillInformation(extractedText)
	if err != nil {
		// Return partial result if parsing fails
		return &models.BillProcessResult{
			RawText:      extractedText,
			MerchantName: "",
			Amount:       "",
			Date:         "",
			Category:     "",
			Confidence: models.Confidence{
				Merchant: 0,
				Amount:   0,
				Date:     0,
				Category: 0,
			},
			ImageType: contentType,
		}, nil
	}

	// Step 3: Format and return the complete result
	return formatBillResponse(parsedBill, extractedText, contentType), nil
}

// callOCRService uses the OCR service to extract text from an image
func callOCRService(imageBytes []byte, filename string) (string, error) {
	// Create multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	filePart.Write(imageBytes)
	writer.Close()

	// Load OCR service authentication config
	config := LoadAuthConfig("ocr")

	if config.ServiceURL == "" {
		return "", fmt.Errorf("OCR service URL not configured")
	}

	// Build full URL with resource path
	// The base URL may or may not have a trailing slash, so handle both cases
	serviceURL := config.ServiceURL
	if !strings.HasSuffix(serviceURL, "/") {
		serviceURL = serviceURL + "/ocr"
	} else {
		serviceURL = serviceURL + "ocr"
	}

	// Create request with the full URL including resource path
	req, err := http.NewRequest("POST", serviceURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create OCR request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	// Add authentication headers
	addAuthHeaders(req, config)

	// Get authenticated client
	client := createAuthenticatedClient(config)

	// Execute OCR request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("OCR service error: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OCR service returned error (status %d): %s",
			resp.StatusCode, string(errBody))
	}

	// Parse OCR response
	var ocrResult struct {
		Text     string `json:"text"`
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ocrResult); err != nil {
		return "", fmt.Errorf("failed to parse OCR response: %v", err)
	}

	return ocrResult.Text, nil
}

// parseBillInformation calls the bill parser service to extract structured data from text
func parseBillInformation(text string) (map[string]interface{}, error) {
	// Load bill parser service authentication config
	config := LoadAuthConfig("parser")

	if config.ServiceURL == "" {
		return nil, fmt.Errorf("Bill parser service URL not configured")
	}

	// Build full URL with resource path
	// The base URL may or may not have a trailing slash, so handle both cases
	serviceURL := config.ServiceURL
	if !strings.HasSuffix(serviceURL, "/") {
		serviceURL = serviceURL + "/process_bill"
	} else {
		serviceURL = serviceURL + "process_bill"
	}

	// Prepare request body
	reqBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to prepare parser request: %v", err)
	}

	// Create request with the full URL including resource path
	req, err := http.NewRequest("POST", serviceURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create parser request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication headers
	addAuthHeaders(req, config)

	// Get authenticated client
	client := createAuthenticatedClient(config)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bill parser service error: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bill parser service returned error (status %d): %s",
			resp.StatusCode, string(errBody))
	}

	// Parse response body
	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read parser response: %v", err)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(resultBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse bill parser response: %v", err)
	}

	return result, nil
}

// formatBillResponse formats the parsed bill information for client response
func formatBillResponse(parsedBill map[string]interface{}, rawText string, contentType string) *models.BillProcessResult {
	// Extract values with type safety
	merchantName, _ := parsedBill["merchant_name"].(string)
	totalAmount, _ := parsedBill["total_amount"].(string)
	date, _ := parsedBill["date"].(string)
	category, _ := parsedBill["category"].(string)
	rawTextFromParser, _ := parsedBill["raw_text"].(string)

	// Use raw text from either parser or OCR
	finalRawText := rawTextFromParser
	if finalRawText == "" {
		finalRawText = rawText
	}

	// Extract confidence scores if they exist
	var merchantConf, amountConf, dateConf, categoryConf float64

	if mc, ok := parsedBill["merchant_confidence"].(float64); ok {
		merchantConf = mc
	}
	if ac, ok := parsedBill["amount_confidence"].(float64); ok {
		amountConf = ac
	}
	if dc, ok := parsedBill["date_confidence"].(float64); ok {
		dateConf = dc
	}
	if cc, ok := parsedBill["category_confidence"].(float64); ok {
		categoryConf = cc
	}

	return &models.BillProcessResult{
		RawText:      finalRawText,
		MerchantName: merchantName,
		Amount:       totalAmount,
		Date:         date,
		Category:     category,
		Confidence: models.Confidence{
			Merchant: merchantConf,
			Amount:   amountConf,
			Date:     dateConf,
			Category: categoryConf,
		},
		ImageType: contentType,
	}
}
