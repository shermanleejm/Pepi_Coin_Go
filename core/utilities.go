package core

import (
	"log"

	"github.com/dgraph-io/badger"
)

func ErrorHandler(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func GetDBValue(item badger.Item) []byte {
	var res []byte
	err := item.Value(func(val []byte) error {
		res = append([]byte{}, val...)
		return nil
	})
	ErrorHandler(err)
	return res
}
