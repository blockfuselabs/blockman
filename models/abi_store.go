package models

import (
	"errors"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/google/uuid"
)

// ABIRecord represents a stored ABI with metadata
type ABIRecord struct {
	ABI       abi.ABI   `json:"abi"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

var (
	abiStore = make(map[string]ABIRecord)
	storeMu  = sync.RWMutex{}

	// Common errors
	ErrABINotFound = errors.New("ABI not found")
)

// SaveABI stores an ABI and returns its ID
func SaveABI(parsedABI abi.ABI) string {
	storeMu.Lock()
	defer storeMu.Unlock()

	abiID := generateUniqueID()
	now := time.Now()

	record := ABIRecord{
		ABI:       parsedABI,
		CreatedAt: now,
		LastUsed:  now,
	}

	abiStore[abiID] = record
	return abiID
}

// GetABI retrieves an ABI by ID
func GetABI(abiID string) (abi.ABI, bool) {
	storeMu.RLock()
	record, exists := abiStore[abiID]
	storeMu.RUnlock()

	// If the ABI exists, update its last used time
	if exists {
		storeMu.Lock()
		// Check again in case it was removed between the read and write locks
		if _, stillExists := abiStore[abiID]; stillExists {
			record.LastUsed = time.Now()
			abiStore[abiID] = record
		}
		storeMu.Unlock()
	}

	return record.ABI, exists
}

// ListABIs returns all stored ABIs with their IDs and metadata
func ListABIs() map[string]ABIRecord {
	storeMu.RLock()
	defer storeMu.RUnlock()

	// Create a copy to avoid external code modifying the internal store
	result := make(map[string]ABIRecord, len(abiStore))
	for id, record := range abiStore {
		result[id] = record
	}

	return result
}

// RemoveABI deletes an ABI from the store
func RemoveABI(abiID string) error {
	storeMu.Lock()
	defer storeMu.Unlock()

	if _, exists := abiStore[abiID]; !exists {
		return ErrABINotFound
	}

	delete(abiStore, abiID)
	return nil
}

// Clean up old ABIs that haven't been used for a while
func CleanupOldABIs(maxAge time.Duration) int {
	storeMu.Lock()
	defer storeMu.Unlock()

	now := time.Now()
	count := 0

	for id, record := range abiStore {
		if now.Sub(record.LastUsed) > maxAge {
			delete(abiStore, id)
			count++
		}
	}

	return count
}

func generateUniqueID() string {
	return uuid.New().String()
}
