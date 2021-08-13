package core

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		ErrorHandler(err)
		lastHash = GetDBValue(*item)
		return err
	})
	ErrorHandler(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialise())
		ErrorHandler(err)
		err = txn.Set([]byte("lastHash"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	ErrorHandler(err)
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Logger = nil
	db, err := badger.Open(opts)
	ErrorHandler(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lastHash")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialise())
			ErrorHandler(err)
			err = txn.Set([]byte("lastHash"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lastHash"))
			ErrorHandler(err)
			lastHash = GetDBValue(*item)
			return err
		}
	})
	ErrorHandler(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		ErrorHandler(err)
		block = Deserialise(GetDBValue(*item))
		return err
	})
	ErrorHandler(err)

	iter.CurrentHash = block.PrevHash

	return block
}
