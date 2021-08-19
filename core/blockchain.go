package core

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
	reward = 69
)

type BlockChain struct {
	LashHash []byte
	Database *badger.DB
	pending  []*Transaction
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

func NewBlockChain() *BlockChain {
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Logger = nil

	if !DBExists() {
		db, err := badger.Open(opts)
		ErrorHandler(err)
		err = db.Update(func(txn *badger.Txn) error {
			genesis := NewBlock([]*Transaction{RewardTransaction([]byte("toaa"), 0)}, nil)
			ErrorHandler(txn.Set(genesis.Hash, genesis.Serialise()))
			ErrorHandler(txn.Set([]byte("lastHash"), genesis.Hash))
			lastHash = genesis.Hash
			return nil
		})
		ErrorHandler(err)

		blockchain := BlockChain{lastHash, db, []*Transaction{}}

		return &blockchain
	} else {
		db, err := badger.Open(opts)
		ErrorHandler(err)
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

		blockchain := BlockChain{lastHash, db, []*Transaction{}}

		return &blockchain
	}
}

func (chain *BlockChain) MineBlock(address []byte) *Block {
	var lastHash []byte
	reward := Transaction{time.Now().Unix(), []byte("toaa"), address, reward, nil}
	chain.pending = append(chain.pending, &reward)

	for _, t := range chain.pending {
		if !t.IsReward() && !t.Verify() {
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

	newBlock := NewBlock(chain.pending, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		ErrorHandler(txn.Set(newBlock.Hash, newBlock.Serialise()))
		ErrorHandler(txn.Set([]byte("lastHash"), newBlock.Hash))
		chain.LashHash = newBlock.Hash
		return nil
	})
	ErrorHandler(err)

	return newBlock
}

type BlockChainIterator struct {
	Current  []byte
	Database *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LashHash, chain.Database}
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	ErrorHandler(iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.Current)
		ErrorHandler(err)
		encoded := GetDBValue(*item)
		block = DeserialiseBlock(encoded)
		fmt.Println(block.Transactions, "desrialised block <========")
		return nil
	}))

	iter.Current = block.PrevHash

	return block
}

func (bc *BlockChain) NewTransaction(wallet *Wallet, to []byte, amount float64) {
	available := bc.GetAvailableBalance()

	if available < amount {
		log.Panic("Not enough funds")
	}
	txn := Transaction{time.Now().Unix(), wallet.PublicKey, to, amount, nil}
	txn.Sign(wallet.PrivateKey)
	bc.pending = append(bc.pending, &txn)
}
