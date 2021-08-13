package core

import (
	"testing"
)

func TestGenesisToNumbers(t *testing.T) {
	chain := InitBlockChain()

	chain.AddBlock("Exodus")
	chain.AddBlock("Leviticus")
	chain.AddBlock("Numbers")

	if string(chain.Blocks[2].Data) != "Leviticus" {
		t.Errorf("Blockchain Failed %s", string(chain.Blocks[2].Data))
	}
}

func TestProofOfWork(t *testing.T) {
	chain := InitBlockChain()
	chain.AddBlock("Exodus")
	pow := NewProofOfWork(chain.GetLatestBlock())
	
	if pow.Validate() == false {
		t.Error("Proof of work has failed")
	}
}
