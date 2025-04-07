package utils

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

var ethClient *rpc.Client

// Initialize Ethereum client
func InitEthClient(nodeURL string) error {
	client, err := rpc.Dial(nodeURL)
	if err != nil {
		return err
	}
	ethClient = client
	return nil
}

// EthCall simulates a function call to a smart contract
func EthCall(contractAddress string, data string) (string, error) {
	address := common.HexToAddress(contractAddress)
	msg := ethereum.CallMsg{
		To:   &address,
		Data: common.FromHex(data),
	}

	var result []byte
	err := ethClient.CallContext(context.Background(), &result, "eth_call", msg, big.NewInt(0))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(result), nil
}

// DecodeCallResult attempts to decode the hex result returned from an eth_call
// based on the ABI output types
func DecodeCallResult(hexResult string, outputs abi.Arguments) (interface{}, error) {
	// If the result is empty or just 0x, return an empty result
	if hexResult == "" || hexResult == "0x" {
		return nil, nil
	}

	// Convert hex string to bytes
	resultBytes, err := hex.DecodeString(hexResult)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex result: %w", err)
	}

	// If there are no outputs defined, return the raw bytes
	if len(outputs) == 0 {
		return resultBytes, nil
	}

	// Use the ABI to unpack the result
	result, err := outputs.Unpack(resultBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	// If there's only one return value, return it directly
	if len(result) == 1 {
		return result[0], nil
	}

	// Create a map of named outputs
	namedOutputs := make(map[string]interface{})
	for i, output := range outputs {
		if i < len(result) {
			// Use the output name if available, otherwise use the index
			name := output.Name
			if name == "" {
				name = fmt.Sprintf("output%d", i)
			}
			namedOutputs[name] = result[i]
		}
	}

	return namedOutputs, nil
}
