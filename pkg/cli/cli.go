package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/zivlakmilos/go-blockchain/pkg/blockchain"
	"github.com/zivlakmilos/go-blockchain/pkg/utils"
)

type CommandLine struct{}

func NewCommandLine() *CommandLine {
	return &CommandLine{}
}

func (c *CommandLine) printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  balance -address ADDRESS - get balance for an address\n")
	fmt.Printf("  create -address ADDRESS - creates a blockchain and sends genesis transaction to address\n")
	fmt.Printf("  print - Prints the blocks in the chain\n")
	fmt.Printf("  send -from FROM -to TO -amount AMOUNT - Send amount of coins\n")
}

func (c *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		c.printUsage()
		runtime.Goexit()
	}
}

func (c *CommandLine) Run() {
	c.validateArgs()

	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	balanceAddress := balanceCmd.String("address", "", "Address")
	createAddress := createCmd.String("address", "", "Address")

	sendFrom := sendCmd.String("from", "", "From")
	sendTo := sendCmd.String("to", "", "To")
	sendAmount := sendCmd.Int("amount", 0, "Amount")

	switch os.Args[1] {
	case "balance":
		err := balanceCmd.Parse(os.Args[2:])
		utils.HandleError(err)
	case "create":
		err := createCmd.Parse(os.Args[2:])
		utils.HandleError(err)
	case "print":
		err := printCmd.Parse(os.Args[2:])
		utils.HandleError(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utils.HandleError(err)
	default:
		c.printUsage()
		runtime.Goexit()
	}

	if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			balanceCmd.Usage()
			runtime.Goexit()
		}
		c.handleBalance(*balanceAddress)
	}

	if createCmd.Parsed() {
		if *createAddress == "" {
			createCmd.Usage()
			runtime.Goexit()
		}
		c.handleCreate(*createAddress)
	}

	if printCmd.Parsed() {
		c.handlePrint()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		c.handleSend(*sendFrom, *sendTo, *sendAmount)
	}
}

func (c *CommandLine) handleBalance(address string) {
	chain := blockchain.OpenBlockChain("")
	defer chain.Database.Close()

	amount := 0

	UTXOs := chain.FindUTXO(address)
	for _, txo := range UTXOs {
		amount += txo.Value
	}

	fmt.Printf("Balance [%s]: %d\n", address, amount)
}

func (c *CommandLine) handleCreate(address string) {
	chain := blockchain.NewBlockChain(address)
	defer chain.Database.Close()
}

func (c *CommandLine) handlePrint() {
	chain := blockchain.OpenBlockChain("")
	defer chain.Database.Close()

	fmt.Printf("%v\n", chain)
}

func (c *CommandLine) handleSend(from, to string, amount int) {
	chain := blockchain.OpenBlockChain("")
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Printf("Success!\n")
}
