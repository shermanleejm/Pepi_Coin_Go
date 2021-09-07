package main

import (
	"os"

	"github.com/shermanleejm/pepi_coin/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
