package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//实现挖矿功能 pow

// 字段：
// ​	区块：block
// ​	目标值：target
// 方法：
// ​	run计算
// ​	功能：找到nonce，从而满足哈希币目标值小

type ProofOfWork struct {
	block *Block
	//目标值，与生成哈希值作比较
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork  {

	pow := ProofOfWork{
		block:block,
	}

	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"

	tmpBigInt := new(big.Int)

	tmpBigInt.SetString(targetStr,16)

	pow.target = tmpBigInt

	return &pow
}

func (pow *ProofOfWork) Run() ([]byte,uint64) {
	var nonce uint64
	var hash [32]byte
	fmt.Println("开始挖矿。。。")

	for  {
		fmt.Printf("%x\r",hash[:])
		//拼接字符串  + nonce
		data := pow.PrepareData(nonce)
		//哈希值
		hash = sha256.Sum256(data)
		//哈希->bigint
		tmpInt := new(big.Int)
		tmpInt.SetBytes(hash[:])

		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//当前计算的哈希.Cmp(难度值)
		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功,hash :%x,nouce:%d\n",hash[:],nonce)
			break
		}else {
			nonce++
		}

	}
	return hash[:], nonce
}

func (pow *ProofOfWork)PrepareData(nonce uint64) []byte  {
	b := pow.block

	tmp := [][]byte{
		uint2Byte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		uint2Byte(b.TimeStamp),
		uint2Byte(b.Bits),
		uint2Byte(nonce),
		b.Hash,
		b.Data,
	}

	data := bytes.Join(tmp, []byte{})

	return data
}
