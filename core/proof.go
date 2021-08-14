package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// TODO: should increment over time
const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	ErrorHandler(binary.Write(buff, binary.BigEndian, num))
	return buff.Bytes()
}

func (pow *ProofOfWork) Init(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevHash,
		pow.Block.HashTransactions(),
		ToHex(int64(nonce)),
		ToHex(int64(Difficulty)),
	}, []byte{})
	return data
}

func (pow *ProofOfWork) Calculate() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.Init(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.Init(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
