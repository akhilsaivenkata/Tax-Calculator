package main

import (
	"log"
	"os"

	"github.com/akhilsaivenkata/go-tax-calculator/internal/client"
	"github.com/akhilsaivenkata/go-tax-calculator/internal/handler"
	"github.com/akhilsaivenkata/go-tax-calculator/internal/service"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/logger"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {

	logger.Init()

	r := gin.New()
	r.Use(gin.Recovery())             // To recover from panics
	r.Use(middleware.RequestLogger()) // This is for the custom structured logger

	// Init dependencies
	apiBaseURL := os.Getenv("TAX_API_URL")
	if apiBaseURL == "" {
		log.Fatal("TAX_API_URL is not set")
	}

	apiClient := client.NewTaxAPIClient(apiBaseURL)
	//apiClient := client.NewTaxAPIClient("http://localhost:5001") [Uncomment this line if you want to run it locally without env variables]
	taxService := service.NewTaxService()
	taxHandler := handler.NewTaxHandler(apiClient, taxService)

	taxHandler.RegisterRoutes(r)

	logger.Log.Info("Server starting on :8080")
	r.Run(":8080")
}
