package main

import (
	"fmt"
)

func main() {
	chain := NewBlockchain()

	chain.Add("My First Block")
	chain.Add("My Second Block")
	chain.Add("My Third Block")

	for _, block := range chain.Blocks {
		fmt.Println(block.ToString())
	}
}
