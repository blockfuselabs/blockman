package handlers

import (
	"encoding/hex"
	"net/http"

	"github.com/blockfuselabs/blockman/models"

	"github.com/blockfuselabs/blockman/utils"

	"github.com/gin-gonic/gin"
)

// CallFunction calls a function from the uploaded ABI
func CallFunction(c *gin.Context) {
	var input struct {
		ABIID         string   `json:"abi_id"`
		ContractAddr  string   `json:"contract_address"`
		FunctionName  string   `json:"function_name"`
		FunctionInput []string `json:"function_input"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Retrieve ABI
	parsedABI, exists := models.GetABI(input.ABIID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ABI not found"})
		return
	}

	// Get function signature
	method, exists := parsedABI.Methods[input.FunctionName]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function not found in ABI"})
		return
	}

	// Pack function call data
	args := make([]interface{}, len(input.FunctionInput))
	for i, arg := range input.FunctionInput {
		args[i] = arg
	}

	callData, err := method.Inputs.Pack(args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode function call"})
		return
	}

	data := append(method.ID, callData...)

	// Perform eth_call
	result, err := utils.EthCall(input.ContractAddr, hex.EncodeToString(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
