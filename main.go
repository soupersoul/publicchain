package main

import (
	"fmt"

	"github.com/soupersoul/publicchain/BLC"
)

func main() {
	god := BLC.Miner{PubKey: "thisisthegodaddress"}
	blockchain := god.NewBlockchain()
	defer blockchain.Close()

	// create a new account
	//blockchain.AddBlock("Send 100 to Num1")
	// blockchain.AddBlock("Send 200 to Num2")
	fmt.Println(blockchain)
}
