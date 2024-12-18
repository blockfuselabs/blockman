package models

import (
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/google/uuid"
)

var (
	abiStore = make(map[string]abi.ABI)
	storeMu  = sync.Mutex{}
)

// SaveABI stores an ABI and returns its ID
func SaveABI(parsedABI abi.ABI) string {
	storeMu.Lock()
	defer storeMu.Unlock()

	abiID := generateUniqueID()
	abiStore[abiID] = parsedABI
	return abiID
}

// GetABI retrieves an ABI by ID
func GetABI(abiID string) (abi.ABI, bool) {
	storeMu.Lock()
	defer storeMu.Unlock()

	parsedABI, exists := abiStore[abiID]
	return parsedABI, exists
}

func generateUniqueID() string {
	return uuid.New().String()
}
