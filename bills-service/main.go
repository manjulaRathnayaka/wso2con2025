package main

import (
	"bills-service/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create router
	r := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Create bill handler
	billHandler := handlers.NewBillHandler()

	// Define API endpoints
	api := r.Group("/api")
	{
		bills := api.Group("/bills")
		{
			// 1. Process bill image (OCR + parsing)
			bills.POST("/process", billHandler.ProcessBillImage)

			// 2. Store bill information
			bills.POST("", billHandler.StoreBill)

			// 3. Get bill information
			bills.GET("", billHandler.GetBills)
			bills.GET("/:id", billHandler.GetBillByID)
		}
	}

	// Start server with hardcoded port
	log.Printf("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
