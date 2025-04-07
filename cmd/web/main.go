package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/blockfuselabs/blockman/handlers"
	"github.com/blockfuselabs/blockman/models"
	"github.com/blockfuselabs/blockman/utils"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort           = "8080"
	defaultCleanupHours   = 24
	defaultCleanupEnabled = true
)

func main() {
	// Load configuration
	_ = godotenv.Load()

	// Ethereum node configuration
	ethNodeURL := os.Getenv("ETH_NODE_URL")
	if ethNodeURL == "" {
		log.Fatal("ETH_NODE_URL environment variable must be set")
	}

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// ABI cleanup configuration
	cleanupEnabledStr := os.Getenv("CLEANUP_ENABLED")
	cleanupEnabled := defaultCleanupEnabled
	if cleanupEnabledStr != "" {
		var err error
		cleanupEnabled, err = strconv.ParseBool(cleanupEnabledStr)
		if err != nil {
			log.Printf("Warning: Invalid CLEANUP_ENABLED value, using default (%v)", defaultCleanupEnabled)
		}
	}

	cleanupHoursStr := os.Getenv("CLEANUP_HOURS")
	cleanupHours := defaultCleanupHours
	if cleanupHoursStr != "" {
		parsed, err := strconv.Atoi(cleanupHoursStr)
		if err != nil {
			log.Printf("Warning: Invalid CLEANUP_HOURS value, using default (%d)", defaultCleanupHours)
		} else {
			cleanupHours = parsed
		}
	}

	// Initialize Ethereum client
	if err := utils.InitEthClient(ethNodeURL); err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	log.Printf("Connected to Ethereum node: %s", ethNodeURL)

	// Start periodic cleanup if enabled
	if cleanupEnabled {
		cleanupDuration := time.Duration(cleanupHours) * time.Hour
		log.Printf("ABI cleanup enabled, will remove ABIs unused for %v", cleanupDuration)
		go startPeriodicCleanup(cleanupDuration)
	}

	// Set up router
	r := setupRouter()

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Routes
	r.POST("/upload-abi", handlers.UploadABI)
	r.POST("/list-functions", handlers.ListFunctions)
	r.POST("/call-function", handlers.CallFunction)

	// Add new routes
	r.GET("/abis", handlers.ListABIs)
	r.DELETE("/abis/:id", handlers.RemoveABI)

	return r
}

func startPeriodicCleanup(maxAge time.Duration) {
	ticker := time.NewTicker(maxAge / 2) // Run cleanup at half the max age interval
	defer ticker.Stop()

	for range ticker.C {
		count := models.CleanupOldABIs(maxAge)
		if count > 0 {
			log.Printf("Cleaned up %d unused ABIs", count)
		}
	}
}
