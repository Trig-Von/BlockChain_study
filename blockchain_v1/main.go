package main

import (
	"fmt"
	"time"
)




//打印区块连
func main() {
	bc := NewBlockChain()

	time.Sleep(1*time.Second)

	bc.AddBlock("可以燎原")
	bc.AddBlock("!!!!!!")

	for i, block := range bc.Blocks {
		fmt.Printf("当前区块高度:%d\n", i)
		fmt.Printf("Version:%d\n", block.Version)
		fmt.Printf("PrevHash:%x\n", block.PrevHash)
		fmt.Printf("MerkleRoot:%x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp:%d\n", block.TimeStamp)
		fmt.Printf("Bits:%d\n", block.Bits)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)
	}
}
