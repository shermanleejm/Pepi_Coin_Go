package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"time"
)

type Transaction struct {
	timestamp int64
	from      []byte
	to        []byte
	amount    float64
	signature []byte
}

func (txn *Transaction) Serialise() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(txn)
	ErrorHandler(err)
	return encoded.Bytes()
}

func DeserialiseTransaction(data []byte) Transaction {
	var txn Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&txn)
	ErrorHandler(err)
	return txn
}

func (txn *Transaction) Sign(privateKey ecdsa.PrivateKey) {
	data := fmt.Sprintf("%d%x%x%f", txn.timestamp, txn.from, txn.to, txn.amount)
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(data))
	ErrorHandler(err)
	signature := append(r.Bytes(), s.Bytes()...)
	txn.signature = signature
}

func (txn *Transaction) Verify() bool {
	data := fmt.Sprintf("%d%x%x%f", txn.timestamp, txn.from, txn.to, txn.amount)

	r := big.Int{}
	s := big.Int{}
	sigLen := len(txn.signature)
	r.SetBytes(txn.signature[:(sigLen / 2)])
	s.SetBytes(txn.signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(txn.from)
	x.SetBytes(txn.from[:(keyLen / 2)])
	y.SetBytes(txn.from[(keyLen / 2):])
	rawPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	return ecdsa.Verify(&rawPubKey, []byte(data), &r, &s)
}

func (bc *BlockChain) NewTransaction(wallet *Wallet, to []byte, amount float64) Transaction {
	available := bc.GetAvailableBalance()

	if amount < available {
		log.Panic("Not enough funds")
	}

	txn := Transaction{time.Now().Unix(), wallet.PublicKey, to, amount, nil}
	txn.Sign(wallet.PrivateKey)
	return txn
}
