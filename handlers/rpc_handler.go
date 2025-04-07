package handlers

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"github.com/blockfuselabs/blockman/models"
	"github.com/blockfuselabs/blockman/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// CallFunction calls a function from the uploaded ABI
func CallFunction(c *gin.Context) {
	var input struct {
		ABIID         string        `json:"abi_id"`
		ContractAddr  string        `json:"contract_address"`
		FunctionName  string        `json:"function_name"`
		FunctionInput []interface{} `json:"function_input"`
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

	// Validate contract address
	if input.ContractAddr == "" || !isValidEthereumAddress(input.ContractAddr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid contract address is required"})
		return
	}

	// Get function signature
	method, exists := parsedABI.Methods[input.FunctionName]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function not found in ABI"})
		return
	}

	// Validate input length
	if len(input.FunctionInput) != len(method.Inputs) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Expected %d arguments, got %d", len(method.Inputs), len(input.FunctionInput)),
		})
		return
	}

	// Convert input parameters to the correct types
	args := make([]interface{}, len(input.FunctionInput))
	for i, arg := range input.FunctionInput {
		var err error
		args[i], err = convertArgument(arg, method.Inputs[i].Type.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Failed to convert argument %d: %s", i, err.Error()),
			})
			return
		}
	}

	// Pack function call data
	callData, err := method.Inputs.Pack(args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode function call: " + err.Error(),
		})
		return
	}

	data := append(method.ID, callData...)

	// Determine if the function is a read-only operation or a state-changing transaction
	var result string
	if method.StateMutability == "view" || method.StateMutability == "pure" {
		// Use eth_call for read-only functions
		result, err = utils.EthCall(input.ContractAddr, hex.EncodeToString(data))
	} else {
		// For state-changing functions, inform the user this requires a transaction
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "State-changing functions require transaction signing, which is not yet supported",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Try to decode the result if possible
	decodedResult, err := utils.DecodeCallResult(result, method.Outputs)
	if err != nil {
		// If decoding fails, return the raw result
		c.JSON(http.StatusOK, gin.H{"result": result, "raw": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": decodedResult, "raw": false})
	}
}

// convertArgument converts a generic input value to the correct type based on the ABI type
func convertArgument(input interface{}, abiType string) (interface{}, error) {
	// Handle input as string if it comes as string from JSON
	if strVal, ok := input.(string); ok {
		// Address type
		if strings.HasPrefix(abiType, "address") {
			if !isValidEthereumAddress(strVal) {
				return nil, fmt.Errorf("invalid Ethereum address: %s", strVal)
			}
			return common.HexToAddress(strVal), nil
		}

		// Uint types
		if strings.HasPrefix(abiType, "uint") {
			bigInt, ok := new(big.Int).SetString(strVal, 10)
			if !ok {
				return nil, fmt.Errorf("failed to parse uint: %s", strVal)
			}
			return bigInt, nil
		}

		// Int types
		if strings.HasPrefix(abiType, "int") {
			bigInt, ok := new(big.Int).SetString(strVal, 10)
			if !ok {
				return nil, fmt.Errorf("failed to parse int: %s", strVal)
			}
			return bigInt, nil
		}

		// Bool type
		if abiType == "bool" {
			switch strings.ToLower(strVal) {
			case "true", "1":
				return true, nil
			case "false", "0":
				return false, nil
			default:
				return nil, fmt.Errorf("invalid boolean value: %s", strVal)
			}
		}

		// Bytes and string types
		if abiType == "string" {
			return strVal, nil
		}
		if strings.HasPrefix(abiType, "bytes") {
			// Remove 0x prefix if present
			if strings.HasPrefix(strVal, "0x") {
				strVal = strVal[2:]
			}
			return hex.DecodeString(strVal)
		}
	}

	// Handle input as number if it comes as number from JSON
	if numVal, ok := input.(float64); ok {
		// For boolean values sent as 0 or 1
		if abiType == "bool" {
			return numVal != 0, nil
		}

		// For integer types
		if strings.HasPrefix(abiType, "uint") || strings.HasPrefix(abiType, "int") {
			return big.NewInt(int64(numVal)), nil
		}
	}

	// Handle input as boolean if it comes as boolean from JSON
	if boolVal, ok := input.(bool); ok {
		if abiType == "bool" {
			return boolVal, nil
		}
	}

	return nil, fmt.Errorf("unsupported type conversion from %T to %s", input, abiType)
}

func isValidEthereumAddress(address string) bool {
	// Check if the address starts with '0x' and is 42 characters long
	if !strings.HasPrefix(address, "0x") || len(address) != 42 {
		return false
	}

	// Validate the address contains only hexadecimal characters after '0x'
	matched, _ := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address)
	return matched
}
