package models

import "time"

// Bill represents a bill record
type Bill struct {
    ID           int       `json:"id"`
    MerchantName string    `json:"merchant_name"`
    Amount       float64   `json:"amount"`
    Date         string    `json:"date"` // Format: YYYY-MM-DD
    Category     string    `json:"category"`
    Notes        string    `json:"notes,omitempty"`
    RawText      string    `json:"raw_text,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
}

// CreateBillRequest is used when creating a new bill
type CreateBillRequest struct {
    MerchantName string  `json:"merchant_name" binding:"required"`
    Amount       float64 `json:"amount" binding:"required"`
    Date         string  `json:"date" binding:"required"`
    Category     string  `json:"category" binding:"required"`
    Notes        string  `json:"notes"`
    RawText      string  `json:"raw_text"`
}
