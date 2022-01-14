package main

import (
	"os"

	"blockchain/cmd"
)

func main() {
	defer os.Exit(0)

	cmd := cmd.CliCmd{}
	cmd.Run()
}
