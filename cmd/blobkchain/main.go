package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/zivlakmilos/go-blockchain/pkg/blockchain"
	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func NewCommandLine(blockchain *blockchain.BlockChain) *CommandLine {
	return &CommandLine{
		blockchain: blockchain,
	}
}

func (c *CommandLine) printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  add --block BLOCK_DATA - Add a block to the chain\n")
	fmt.Printf("  print - Prints the blocks in the chain\n")
}

func (c *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		c.printUsage()
		runtime.Goexit()
	}
}

func (c *CommandLine) addBlock(data string) {
	c.blockchain.AddBlock(data)
	fmt.Printf("Added Block!")
}

func (c *CommandLine) printChain() {
	fmt.Printf("%v\n", c.blockchain)
}

func (c *CommandLine) run() {
	c.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		utils.HandleError(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		utils.HandleError(err)
	default:
		c.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		c.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		c.printChain()
	}
}

func main() {
	defer os.Exit(0)

	chain := blockchain.NewBlockChain()
	defer chain.Database.Close()

	cli := NewCommandLine(chain)
	cli.run()
}
