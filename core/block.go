package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Block struct {
	Hash     []byte
	Data     []*Transaction
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, data, prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Calculate()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Data {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) Serialise() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	ErrorHandler(encoder.Encode(b))

	return res.Bytes()
}

func Deserialise(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	ErrorHandler(decoder.Decode(&block))

	return &block
}
