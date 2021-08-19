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
	timestamp int64
	from      []byte
	to        []byte
	Amount    float64
	signature []byte
}

type encodableTransaction struct {
	Timestamp int64
	From      []byte
	To        []byte
	Amount    float64
	Signature []byte
}

func (txn Transaction) Serialise() []byte {
	tmp := encodableTransaction{txn.timestamp, txn.from, txn.to, txn.Amount, txn.signature}
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tmp)
	ErrorHandler(err)
	return encoded.Bytes()
}

func DeserialiseTransaction(data []byte) Transaction {
	var tmp encodableTransaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tmp)
	ErrorHandler(err)
	txn := Transaction{tmp.Timestamp, tmp.From, tmp.To, tmp.Amount, tmp.Signature}
	return txn
}

func (txn *Transaction) Sign(privateKey ecdsa.PrivateKey) {
	data := fmt.Sprintf("%d%x%x%f", txn.timestamp, txn.from, txn.to, txn.Amount)
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(data))
	ErrorHandler(err)
	signature := append(r.Bytes(), s.Bytes()...)
	txn.signature = signature
}

func (txn *Transaction) Verify() bool {
	data := fmt.Sprintf("%d%x%x%f", txn.timestamp, txn.from, txn.to, txn.Amount)

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

func (txn *Transaction) IsReward() bool {
	return bytes.Equal(txn.from, []byte("toaa")) && txn.signature == nil
}

func RewardTransaction(to []byte, reward float64) *Transaction {
	return &Transaction{time.Now().Unix(), []byte("toaa"), to, reward, nil}
}

