package main

import "log"

func main() {
	chain := NewBlockchain()

	chain.Add("My First Block")
	chain.Add("My Second Block")
	chain.Add("My Third Block")

	for _, block := range chain.Blocks {
		log.Println(block.ToString())
	}
}
