package main

import (
	"log"

	"github.com/blockfuselabs/blockman/handlers"
	"github.com/blockfuselabs/blockman/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := utils.InitEthClient(""); err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	r := gin.Default()

	// Routes
	r.POST("/upload-abi", handlers.UploadABI)
	r.POST("/list-functions", handlers.ListFunctions)
	r.POST("/call-function", handlers.CallFunction)

	// Start server
	log.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
