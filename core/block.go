package core

import (
	"bytes"
	"encoding/gob"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Calculate()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
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
