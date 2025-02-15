package merge

import (
	"fmt"
	"sync"

	"github.com/aclindsa/ofxgo"
)

type transactionHash string

type transactionSet struct {
	mu           sync.RWMutex
	transactions map[transactionHash]int
}

func NewTransactionSet() *transactionSet {
	return &transactionSet{
		transactions: make(map[transactionHash]int),
	}
}

// isDuplicate checks if a transaction has already been processed at a different index.
// It returns true if the transaction exists at a different index, false if it's new
// or exists at the same index.
//
// The function uses a transaction hash as a unique identifier and stores the index
// where the transaction was first seen. If the transaction is found at the same
// index, it's not considered a duplicate to handle cases where the same transaction
// appears multiple times in the same ofx file.
//
// Parameters:
//   - idx: the current position in the transaction list
//   - transaction: the transaction to check for duplicates
//
// Returns:
//   - bool: true if duplicate at different index, false if new or same index
func (t *transactionSet) isDuplicate(idx int, transaction ofxgo.Transaction) bool {
	hash := t.transactionHash(transaction)
	t.mu.RLock()
	i, ok := t.transactions[hash]
	t.mu.RUnlock()

	if !ok {
		t.mu.Lock()
		t.transactions[hash] = idx
		t.mu.Unlock()
		return false
	}
	if ok && i == idx {
		return false
	}
	return true
}

func (t *transactionSet) transactionHash(transaction ofxgo.Transaction) transactionHash {
	return transactionHash(fmt.Sprintf("%s%s%s%d",
		transaction.TrnType,
		transaction.DtPosted.Format("2006-01-02"),
		transaction.Name,
		transaction.TrnAmt.String()))
}
