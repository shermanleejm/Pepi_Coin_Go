package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty)) // locality sensitive hashing to reduce dimentionality of big int
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) CalculateHash(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevHash,
		pow.Block.HashTransactions(),
		ToHex(int64(nonce)),
		ToHex(int64(Difficulty)),
	}, []byte{})
	return data
}

func ToHex(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	ErrorHandler(err)
	return buffer.Bytes()
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.CalculateHash(pow.Block.Nonce)
	tmp := sha256.Sum256(data)
	intHash.SetBytes(tmp[:])
	return intHash.Cmp(pow.Target) == -1 // check if the target is correct
}

func (pow *ProofOfWork) Init() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0

	for nonce < math.MaxInt64 {
		hash = sha256.Sum256(pow.CalculateHash(nonce))
		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}
	fmt.Println()
	return nonce, hash[:]
}
