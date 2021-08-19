package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, txn := range b.Transactions {
		txHashes = append(txHashes, txn.Serialise())
	}
	tree := NewMerkleTree(txHashes)
	return tree.Root.Data
}

func NewBlock(txns []*Transaction, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txns, prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Init()
	block.Nonce = nonce
	block.Hash = hash[:]
	fmt.Println(block.Transactions[0], "NEW BLOCK <<<<<<<<")
	return block
}

func (b *Block) Serialise() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	ErrorHandler(err)
	return res.Bytes()
}

func DeserialiseBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	ErrorHandler(err)
	return &block
}
