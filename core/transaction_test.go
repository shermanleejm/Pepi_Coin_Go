package core

import (
	"testing"
)

func TestMakeNewTransaction(t *testing.T) {
	txn := Transaction{}

	if txn.fromAddress != "" {
		t.Error("Transaction constructor is invalid")
	}
}

func TestTransactionHash(t *testing.T) {
	txn := Transaction{}
	if getHash(txn) == "" {
		t.Error("Hashing failed")
	}
}
