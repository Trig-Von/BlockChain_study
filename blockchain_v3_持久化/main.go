package main

import (
	"fmt"
)

//打印区块连
func main() {
	//bc := NewBlockChain()

	err := CreateBlockChain()
	fmt.Println("err:",err)

	bc,err := GetBlockChainInstance()
	defer bc.db.Close()

	if err != nil {
		fmt.Println("GetBlockChainInstance err :",err)
		return
	}

	err = bc.AddBlock("hello moto!!")
	if err != nil {
		fmt.Println("AddBlock err :",err)
		return
	}

	err = bc.AddBlock("hello nokia!!!")
	if err != nil {
		fmt.Println("ADdBlock err:",err)
		return
	}

	//调用迭代器，输出blockChain
	it :=bc.NewIterator()
	for  {
		block := it.Next()
		fmt.Printf("\n***************************\n")
		fmt.Printf("Version:%d\n", block.Version)
		fmt.Printf("PrevHash:%x\n", block.PrevHash)
		fmt.Printf("MerkleRoot:%x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp:%d\n", block.TimeStamp)
		fmt.Printf("Bits:%d\n", block.Bits)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Data:%s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid : %v\n",pow.IsValid())

		if block.PrevHash == nil {
			fmt.Println("区块链遍历结束！！")
			break
		}
	}
}
