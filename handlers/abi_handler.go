package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/blockfuselabs/blockman/models"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gin-gonic/gin"
)

type ABI struct {
	ABI []interface{} `json:"abi"`
}

// UploadABI handles ABI upload and parsing
func UploadABI(c *gin.Context) {
	// Define a struct to match the incoming request format
	var input struct {
		ABI json.RawMessage `json:"abi" binding:"required"`
	}

	// Bind the JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Parse ABI directly using go-ethereum's abi package
	contractABI, err := abi.JSON(bytes.NewReader(input.ABI))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ABI format",
			"details": err.Error(),
		})
		return
	}

	// Save ABI and generate an ID
	abiID := models.SaveABI(contractABI)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "ABI uploaded successfully",
		"abi_id":  abiID,
	})
}

// ListFunctions retrieves and lists all functions in an uploaded ABI
func ListFunctions(c *gin.Context) {
	var input struct {
		ABIID string `json:"abi_id" binding:"required"`
	}

	// Bind and validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Retrieve ABI
	parsedABI, exists := models.GetABI(input.ABIID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "ABI not found",
			"abi_id": input.ABIID,
		})
		return
	}

	// Extract function details
	functions := extractFunctions(parsedABI)

	c.JSON(http.StatusOK, gin.H{
		"functions":       functions,
		"total_functions": len(functions),
	})
}

// ListABIs returns a list of all stored ABIs with their metadata
func ListABIs(c *gin.Context) {
	abis := models.ListABIs()

	// Create a simplified response structure
	response := make(map[string]gin.H, len(abis))
	for id, record := range abis {
		// Get function count
		functions := extractFunctions(record.ABI)

		response[id] = gin.H{
			"created_at":      record.CreatedAt,
			"last_used":       record.LastUsed,
			"function_count":  len(functions),
			"has_constructor": record.ABI.Constructor.Sig != "",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"abis":  response,
		"total": len(abis),
	})
}

// RemoveABI deletes an ABI from storage
func RemoveABI(c *gin.Context) {
	abiID := c.Param("id")
	if abiID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ABI ID is required",
		})
		return
	}

	err := models.RemoveABI(abiID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrABINotFound {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ABI removed successfully",
		"abi_id":  abiID,
	})
}

// FunctionDetail represents a detailed view of a contract function
type FunctionDetail struct {
	Name     string           `json:"name"`
	Inputs   []ArgumentDetail `json:"inputs"`
	Outputs  []ArgumentDetail `json:"outputs"`
	Constant bool             `json:"constant"`
	Payable  bool             `json:"payable"`
	Stateful bool             `json:"stateful"`
}

// ArgumentDetail describes a single argument in a function
type ArgumentDetail struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// extractFunctions converts ABI methods to a more detailed representation
func extractFunctions(parsedABI abi.ABI) []FunctionDetail {
	functions := []FunctionDetail{}

	for name, method := range parsedABI.Methods {
		function := FunctionDetail{
			Name:     name,
			Inputs:   extractArguments(method.Inputs),
			Outputs:  extractArguments(method.Outputs),
			Constant: method.StateMutability == "view" || method.StateMutability == "pure",
			Payable:  method.StateMutability == "payable",
			Stateful: method.StateMutability != "view" && method.StateMutability != "pure",
		}
		functions = append(functions, function)
	}

	return functions
}

// extractArguments converts ABI arguments to a more readable format
func extractArguments(args abi.Arguments) []ArgumentDetail {
	result := []ArgumentDetail{}
	for _, arg := range args {
		result = append(result, ArgumentDetail{
			Name: arg.Name,
			Type: arg.Type.String(),
		})
	}
	return result
}
