package cli

import (
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/shermanleejm/pepi_coin/core"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("createwallet			creates a wallet and saves it in a file")
	fmt.Println("loadwallet				loads your private.pem and public.pem")
	fmt.Println("send -address -amount	provide an address to send to and how much")
	fmt.Println("min					mine the current transactions to your account")
}

func (cli *CommandLine) createWallet() {
	wallet := core.NewWallet()
	wallet.EncodeWalletKeys()
}

func (cli *CommandLine) loadWallet() core.Wallet {
	pri, err := os.ReadFile("private.pem")
	core.ErrorHandler(err)
	pub, err := os.ReadFile("public.pem")
	core.ErrorHandler(err)
	return core.DecodeWalletKeys(string(pri), string(pub))
}

func (cli *CommandLine) send(address string, amount int) {
	blockchain := core.NewBlockChain()
	wallet := cli.loadWallet()
	blockPub, _ := pem.Decode([]byte(address))
	addressBytes := blockPub.Bytes
	blockchain.NewTransaction(&wallet, addressBytes, float64(amount))
}

func (cli *CommandLine) mine() {
	blockchain := core.NewBlockChain()
	wallet := cli.loadWallet()
	blockchain.MineBlock(wallet.PublicKey)
}

func (cli *CommandLine) Run() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendTo := sendCmd.String("address", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createwallet":
		cli.createWallet()
	case "loadwallet":
		cli.loadWallet()
	case "send":
		sendCmd.Parse(os.Args[2:])
	case "mine":
		cli.mine()
	}

	if sendCmd.Parsed() {
		if *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendTo, *sendAmount)
	}

}
