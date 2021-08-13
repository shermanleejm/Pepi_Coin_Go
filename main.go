package main

import (
	"github.com/shermanleejm/pepi_coin/core"
)

func main() {
	chain := core.InitBlockChain()
	chain.AddBlock("Exodus")
	chain.AddBlock("Leviticus")
	chain.AddBlock("Numbers")
}
