package core

import (
	"reflect"
	"testing"
	"time"
)

func TestSerialiseBlock(t *testing.T) {
	w := NewWallet()
	t1 := Transaction{time.Now().Unix(), w.PublicKey, []byte("iron man"), 10, nil}
	var txns []*Transaction
	for i := 0; i < 5; i++ {
		txns = append(txns, &t1)
	}
	block := Block{time.Now().Unix(), []byte{}, txns, []byte{}, 69}
	s := block.Serialise()
	d := DeserialiseBlock(s)
	if !reflect.DeepEqual(d.Transactions[0].Amount, block.Transactions[0].Amount) {
		t.Error("Invalid serialiser")
	}
}
