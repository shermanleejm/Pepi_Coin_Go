package core

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST" // check if db exists
	genesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) AddBlock(data []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		ErrorHandler(err)
		lastHash = GetDBValue(*item)
		return nil
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

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func badgerOptions() *badger.Options {
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	opts.Logger = nil
	return &opts
}

func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBExists() {
		fmt.Println("Chain exists")
		runtime.Goexit()
	}

	db, err := badger.Open(*badgerOptions())
	ErrorHandler(err)

	err = db.Update(func(txn *badger.Txn) error {
		coinbaseTxn := CoinbaseTxn(address, genesisData)
		genesis := Genesis(coinbaseTxn)
		err = txn.Set(genesis.Hash, genesis.Serialise())
		ErrorHandler(err)
		err = txn.Set([]byte("lastHash"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	ErrorHandler(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func ContinueBlockChain(address string) *BlockChain {
	if DBExists() == false {
		fmt.Println("No blockchain, create one")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := badger.Open(*badgerOptions())
	ErrorHandler(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lastHash"))
		ErrorHandler(err)
		lastHash = GetDBValue(*item)
		return nil
	})

	chain := BlockChain{lastHash, db}

	return &chain
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

func (chain *BlockChain) FindUnspentTxns(address string) []Transaction {
	var unspentTxns []Transaction
	spentTXOs := make(map[string][]int)
	iter := chain.Iterator()
	for {
		block := iter.Next()

		for _, txn := range block.Data {
			txnID := hex.EncodeToString(txn.ID)

		Outputs: // label for the secondary for loop
			for outputIndex, out := range txn.Outputs {
				if _, ok := spentTXOs[txnID]; ok {
					for _, spentOut := range spentTXOs[txnID] {
						if spentOut == outputIndex {
							continue Outputs
						}
					}
				}

				// get all the txns that are from this address
				if out.CanBeUnlocked(address) {
					unspentTxns = append(unspentTxns, *txn)
				}

				// get all incoming coins to address
				if !txn.IsCoinbase() {
					for _, in := range txn.Inputs {
						if in.CanUnlock(address) {
							inTxnID := hex.EncodeToString(in.ID)
							spentTXOs[inTxnID] = append(spentTXOs[inTxnID], in.OutputIndex)
						}
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}

	}
	return unspentTxns
}

func (chain *BlockChain) FindUnspentTXO(address string) []TxOutput {
	var UnspentTXOs []TxOutput
	unspentTxns := chain.FindUnspentTxns(address)

	for _, tx := range unspentTxns {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UnspentTXOs = append(UnspentTXOs, out)
			}
		}
	}

	return UnspentTXOs
}

// find the amount that the address can spend
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxns := chain.FindUnspentTxns(address)
	accumulated := 0

	Work:
	for _, txn := range unspentTxns {
		txnID := hex.EncodeToString(txn.ID)

		for outIdx, out := range txn.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txnID] = append(unspentOuts[txnID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
