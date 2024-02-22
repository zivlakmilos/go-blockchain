package main

import (
	"fmt"

	"github.com/zivlakmilos/go-blockchain/pkg/blockchain"
)

func main() {
	chain := blockchain.NewBlockChain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	fmt.Printf("%v\n", chain)
}
