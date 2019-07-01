package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

//定义区块结构
//第一阶段
//第二阶段
type Block struct {
	//版本号
	Version uint64
	//前区块哈希
	PrevHash []byte
	//根哈希
	MerkleRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值
	Bits uint64
	//随机数
	Nonce uint64
	//哈希
	Hash []byte
	//数据
	//Data []byte
	Transactions []*Transaction
}

// 创建一个区块(方法)
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: nil, //
		TimeStamp:  uint64(time.Now().Unix()),
		Bits:       0, //
		Nonce:      0, //
		Hash:       nil,
		///Data:       []byte(data),
		Transactions: txs,
	}

	//填充梅克尔根
	b.HashTransactionMerkleRoot()
	fmt.Printf("merkleroot:%x\n",b.MerkleRoot)

	//计算哈希值
	//b.setHash()
	pow := NewProofOfWork(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce
	return &b
}

//Serialize方法,gob编码
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(b)
	if err != nil {
		fmt.Println("Encode err:", err)
		return nil
	}

	return buffer.Bytes()
}

func Deserialize(src []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(src))

	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("decode err:", err)
		return nil
	}
	return &block
}

func (block *Block) HashTransactionMerkleRoot() {
	var info [][]byte

	for _, tx := range block.Transactions {
		txHashValue := tx.TXID
		info = append(info, txHashValue)
	}
	value := bytes.Join(info, []byte{})
	hash := sha256.Sum256(value)

	block.MerkleRoot = hash[:]
}
