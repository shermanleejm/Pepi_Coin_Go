package core

type BlockChain struct {
	Blocks []*Block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{CreateBlock("Genesis", []byte{})}}
}

func (bc *BlockChain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}
