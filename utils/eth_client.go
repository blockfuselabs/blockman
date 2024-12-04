package utils

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
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
