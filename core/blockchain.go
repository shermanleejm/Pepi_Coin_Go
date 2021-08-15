package core

import (
	"os"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LashHash []byte
	Database *badger.DB
	Pending  []*Transaction
}

// TODO: traverse the blockchain and calculate the  address's available balance
func (bc *BlockChain) GetAvailableBalance() float64 {
	return 100.00
}

func DBExists() bool {
	if _, err := os.Stat(dbPath + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}

// func NewBlockChain(address, nodeID string) *BlockChain {
// 	var lastHash []byte
// 	if !DBExists() {
// 		opts := badger.DefaultOptions(dbPath)
// 		opts.Dir = dbPath
// 		opts.ValueDir = dbPath
// 		opts.Logger = nil
// 		db, err := badger.Open(opts)
// 		ErrorHandler(err)

// 		err = db.Update(func(txn *badger.Txn) error {
// 			genesis := NewBlock()
// 		})
// 	}
// }
