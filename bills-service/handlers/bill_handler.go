package handlers

import (
    "bills-service/models"
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
)

type BillHandler struct {
    // This would normally have a database connection
    bills []models.Bill // In-memory storage for simplicity
    nextID int
}

func NewBillHandler() *BillHandler {
    return &BillHandler{
        bills:  make([]models.Bill, 0),
        nextID: 1,
    }
}

// ProcessBillImage handles uploaded bill images, extracts text with OCR and parses information
func (h *BillHandler) ProcessBillImage(c *gin.Context) {
    // Get file from request (not implemented)
    file, _, err := c.Request.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
        return
    }
    defer file.Close()

    // Dummy implementation (would call OCR service and parser in real implementation)
    // Returns dummy data for demonstration
    c.JSON(http.StatusOK, gin.H{
        "raw_text":      "Walmart Receipt\nDate: 2025-06-15\nTotal: $45.67\nThank you!",
        "merchant_name": "Walmart",
        "amount":        "45.67",
        "date":          "2025-06-15",
        "category":      "Groceries",
        "confidence": gin.H{
            "merchant": 0.95,
            "amount":   0.98,
            "date":     0.96,
            "category": 0.85,
        },
    })
}

// StoreBill stores bill information in the database
func (h *BillHandler) StoreBill(c *gin.Context) {
    var req models.CreateBillRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create a new bill with dummy ID
    bill := models.Bill{
        ID:           h.nextID,
        MerchantName: req.MerchantName,
        Amount:       req.Amount,
        Date:         req.Date,
        Category:     req.Category,
        Notes:        req.Notes,
        RawText:      req.RawText,
        CreatedAt:    time.Now(),
    }

    // Store bill (in memory for this simple implementation)
    h.bills = append(h.bills, bill)
    h.nextID++

    c.JSON(http.StatusCreated, gin.H{
        "message": "Bill created successfully",
        "id":      bill.ID,
    })
}

// GetBills returns all bills
func (h *BillHandler) GetBills(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "data": h.bills,
    })
}

// GetBillByID returns a specific bill by ID
func (h *BillHandler) GetBillByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    // Find the bill
    for _, bill := range h.bills {
        if bill.ID == id {
            c.JSON(http.StatusOK, bill)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
}
