package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"math/big"
	"time"
)

type Transaction struct {
	Timestamp int64
	From      []byte
	To        []byte
	Amount    float64
	Signature []byte
}

func (txn Transaction) Serialise() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(txn)
	ErrorHandler(err)
	return encoded.Bytes()
}

func DeserialiseTransaction(data []byte) Transaction {
	var tmp Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tmp)
	ErrorHandler(err)
	return tmp
}

func (txn *Transaction) Sign(privateKey ecdsa.PrivateKey) {
	data := fmt.Sprintf("%d%x%x%f", txn.Timestamp, txn.From, txn.To, txn.Amount)
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(data))
	ErrorHandler(err)
	signature := append(r.Bytes(), s.Bytes()...)
	txn.Signature = signature
}

func (txn *Transaction) Verify() bool {
	data := fmt.Sprintf("%d%x%x%f", txn.Timestamp, txn.From, txn.To, txn.Amount)

	r := big.Int{}
	s := big.Int{}
	sigLen := len(txn.Signature)
	r.SetBytes(txn.Signature[:(sigLen / 2)])
	s.SetBytes(txn.Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(txn.From)
	x.SetBytes(txn.From[:(keyLen / 2)])
	y.SetBytes(txn.From[(keyLen / 2):])
	rawPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	return ecdsa.Verify(&rawPubKey, []byte(data), &r, &s)
}

func (txn *Transaction) IsReward() bool {
	return bytes.Equal(txn.From, []byte("toaa")) && txn.Signature == nil
}

func RewardTransaction(to []byte, reward float64) *Transaction {
	return &Transaction{time.Now().Unix(), []byte("toaa"), to, reward, nil}
}

