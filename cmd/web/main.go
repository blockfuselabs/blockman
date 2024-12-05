package main

import (
	"log"
	"os"

	"github.com/blockfuselabs/blockman/handlers"
	"github.com/blockfuselabs/blockman/utils"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	_ = godotenv.Load()
	ethNodeURL := os.Getenv("ETH_NODE_URL")
	if err := utils.InitEthClient(ethNodeURL); err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	r := gin.Default()

	// Routes
	r.POST("/upload-abi", handlers.UploadABI)
	r.POST("/list-functions", handlers.ListFunctions)
	r.POST("/call-function", handlers.CallFunction)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
