package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value     int
	PublicKey string
}

type TxInput struct {
	ID          []byte
	OutputIndex int
	Signature   string
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulated, ValidOutputs := chain.FindSpendableOutputs(from, amount)

	if accumulated < amount {
		log.Panic("Not enough funds")
	}

	for id, outs := range ValidOutputs {
		txnID, err := hex.DecodeString(id)
		ErrorHandler(err)

		for _, out := range outs {
			inputs = append(inputs, TxInput{txnID, out, from})
		}
	}

	outputs = append(outputs, TxOutput{amount, to})
	// if enough money the nminus
	if amount <= accumulated {
		outputs = append(outputs, TxOutput{accumulated - amount, from})
	}

	txn := Transaction{nil, inputs, outputs}
	txn.SetID()
	
	return &txn 
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	ErrorHandler(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTxn(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Paying %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{69, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].OutputIndex == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PublicKey == data
}