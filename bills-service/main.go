package main

import (
	"bills-service/handlers"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create router
	r := gin.Default()

	// Configure CORS with a more secure configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*.choreo.dev", "https://*.wso2.com", "http://localhost:*", "http://*.svc.cluster.local"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
