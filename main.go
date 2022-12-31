package main

import (
	"fmt"

	"github.com/soupersoul/publicchain/BLC"
)

func main() {
	god := BLC.Miner{PubKey: "thisisthegodaddress"}
	BLC.InitBlockchain(god)

	// create a new account
	//blockchain.AddBlock("Send 100 to Num1")
	// blockchain.AddBlock("Send 200 to Num2")
	fmt.Println(BLC.BlockChainIns)
}

/*
func main() {
	hasher := ripemd160.New()
	hasher.Write([]byte("datadatadatadata"))
	bytes := hasher.Sum(nil)
	fmt.Println(fmt.Sprintf("%x", bytes))
}
*/
