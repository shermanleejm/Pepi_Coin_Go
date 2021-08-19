package main

import (
	"fmt"

	"github.com/shermanleejm/pepi_coin/core"
)

//Private Ket: 54b005f0c21eddcee9a3e857237116c35cc4b8205d91cf2a32c88a0adb05bab9
//Public Ket: bdf78ccca0dbc8adb7a984c7dbc2b635d19836822f50979de2dc3ef295d11051912285ecd6ce2ce826aeb8c5ea822619a0cb60d0dfa89e730a078468ce5e8392
func main() {
	var testData = [3]float64{1, 2, 3}
	wFrom := core.NewWallet()
	wTo := core.NewWallet()
	chain := core.NewBlockChain()

	for _, amt := range testData {
		chain.NewTransaction(wFrom, wTo.PublicKey, amt)
	}

	chain.MineBlock(wFrom.PublicKey)
	// var verificationData = append(testData[:], 69)
	iter := chain.Iterator()
	block := iter.Next()
	for block.PrevHash != nil {
		for _, txn := range block.Transactions {
			fmt.Println(txn)
		}
		block = iter.Next()
	}
}
