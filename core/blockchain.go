package core

import (
	"log"
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

func NewBlockChain(address, nodeID string) *BlockChain {
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Logger = nil
	db, err := badger.Open(opts)
	ErrorHandler(err)

	if !DBExists() {
		err = db.Update(func(txn *badger.Txn) error {
			genesis := NewBlock([]*Transaction{}, nil)
			ErrorHandler(txn.Set(genesis.Hash, genesis.Serialise()))
			ErrorHandler(txn.Set([]byte("lastHash"), genesis.Hash))
			lastHash = genesis.Hash
			return nil
		})
		ErrorHandler(err)
	} else {
		err = db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("lastHash"))
			ErrorHandler(err)
			err = item.Value(func(val []byte) error {
				lastHash = append(lastHash, val...)
				return nil
			})
			return err
		})
		ErrorHandler(err)
	}

	blockchain := BlockChain{lastHash, db, []*Transaction{}}

	return &blockchain
}

func (chain *BlockChain) MineBlock(txns []*Transaction) *Block {
	var lastHash []byte
	for _, t := range txns {
		if !t.Verify() {
			log.Panic("Invalid transaction")
		}
	}

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		ErrorHandler(err)
		lastHash = GetDBValue(*item)
		return nil
	})
	ErrorHandler(err)

	newBlock := NewBlock(txns, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		ErrorHandler(txn.Set(newBlock.Hash, newBlock.Serialise()))
		ErrorHandler(txn.Set([]byte("lastHash"), newBlock.Hash))
		chain.LashHash = newBlock.Hash
		return nil
	})
	ErrorHandler(err)

	return newBlock
}
