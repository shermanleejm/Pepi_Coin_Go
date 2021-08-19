package core

import (
	"bytes"
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

func theOneAboveAll() []byte {
	return []byte("toaa")
}

// TODO: traverse the blockchain and calculate the  address's available balance
func (bc *BlockChain) GetAvailableBalance(address []byte) float64 {
	iter := bc.Iterator()
	block := iter.Next()
	res := 0
	for block.PrevHash != nil {
		txnLength := len(block.Transactions)
		var tmp Transaction
		for i := 0; i < txnLength; i++ {
			tmp = *block.Transactions[i]
			if bytes.Equal(tmp.From, address) {
				res -= int(tmp.Amount)
			}
			if bytes.Equal(tmp.To, address) {
				res += int(tmp.Amount)
			}
		}
	}
	return float64(res)
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
			genesis := NewBlock([]*Transaction{RewardTransaction(theOneAboveAll(), 0)}, nil)
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

		freeMoney := Transaction{time.Now().Unix(), theOneAboveAll(), theOneAboveAll(), 0, nil}
		var freeMonies []*Transaction
		for i := 0; i < 7; i++ {
			freeMonies = append(freeMonies, &freeMoney)
		}
		blockchain := BlockChain{lastHash, db, freeMonies}

		return &blockchain
	}
}

func (chain *BlockChain) MineBlock(address []byte) *Block {
	var lastHash []byte
	reward := Transaction{time.Now().Unix(), theOneAboveAll(), address, reward, nil}
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
	available := bc.GetAvailableBalance(wallet.PublicKey)

	if available < amount {
		log.Panic("Not enough funds")
	}
	txn := Transaction{time.Now().Unix(), wallet.PublicKey, to, amount, nil}
	txn.Sign(wallet.PrivateKey)
	bc.pending = append(bc.pending, &txn)
}
