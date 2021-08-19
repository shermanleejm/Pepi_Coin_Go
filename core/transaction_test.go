package core

import (
	"reflect"
	"testing"
	"time"
)

func TestTransactions(t *testing.T) {
	w := NewWallet()
	txn := Transaction{time.Now().Unix(), w.PublicKey, []byte("adam"), 1000, nil}
	txn.Sign(w.PrivateKey)
	verification := txn.Verify()

	if !verification {
		t.Error("Invalid implementation")
	}
}

func TestSerialisation(t *testing.T) {
	w := NewWallet()
	txn := Transaction{time.Now().Unix(), w.PublicKey, []byte("adam"), 1000, nil}
	txn.Sign(w.PrivateKey)

	s := txn.Serialise()
	d := DeserialiseTransaction(s)

	if !reflect.DeepEqual(d, txn) {
		t.Error("Error with serialiser")
	}
}
