package main

import (
	"os"

	"github.com/zivlakmilos/go-blockchain/pkg/cli"
)

func main() {
	defer os.Exit(0)

	cli := cli.NewCommandLine()
	cli.Run()
}
