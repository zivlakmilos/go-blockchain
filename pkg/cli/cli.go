package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/zivlakmilos/go-blockchain/pkg/blockchain"
	"github.com/zivlakmilos/go-blockchain/pkg/utils"
	"github.com/zivlakmilos/go-blockchain/pkg/wallet"
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
	fmt.Printf("  createwallet - Create a new wallet\n")
	fmt.Printf("  listwallets - List all wallet addresses\n")
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
	createWallet := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listWallets := flag.NewFlagSet("listwallet", flag.ExitOnError)

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
	case "createwallet":
		err := createWallet.Parse(os.Args[2:])
		utils.HandleError(err)
	case "listwallets":
		err := listWallets.Parse(os.Args[2:])
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

	if createWallet.Parsed() {
		c.handleCreateWallet()
	}

	if listWallets.Parsed() {
		c.handleListWallts()
	}
}

func (c *CommandLine) handleBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("address is not valid")
	}

	chain := blockchain.OpenBlockChain("")
	defer chain.Database.Close()

	amount := 0

	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	UTXOs := chain.FindUTXO(pubKeyHash)
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
	if !wallet.ValidateAddress(from) {
		log.Panic("address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("address is not valid")
	}

	chain := blockchain.OpenBlockChain("")
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})

	fmt.Printf("Success!\n")
}

func (c *CommandLine) handleListWallts() {
	wallets, _ := wallet.NewWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Printf("%s\n", address)
	}
}

func (c *CommandLine) handleCreateWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}
