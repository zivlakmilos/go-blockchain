package blockchain

import "strings"

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*Block{Genesis()},
	}
}

func (c *BlockChain) AddBlock(data string) {
	prevBlock := c.blocks[len(c.blocks)-1]
	block := NewBlock(data, prevBlock.Hash)
	c.blocks = append(c.blocks, block)
}

func (c *BlockChain) String() string {
	var builder strings.Builder

	for _, c := range c.blocks {
		builder.WriteString(c.String())
		builder.WriteString("\n")
	}

	return builder.String()
}
