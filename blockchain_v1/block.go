package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

//定义区块结构
//第一阶段
//第二阶段
type Block struct {
	Version uint64
	//前区块哈希
	PrevHash []byte

	MerkleRoot []byte

	TimeStamp uint64

	Bits uint64

	Nonce uint64
	//哈希
	Hash []byte
	//数据
	Data []byte
}

// 创建一个区块(方法)
func NewBlock(data string, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: nil, //
		TimeStamp:  uint64(time.Now().Unix()),
		Bits:       0, //
		Nonce:      0, //
		Hash:       nil,
		Data:       []byte(data),
	}

	//计算哈希值
	b.setHash()

	return &b
}

func (b *Block) setHash() {

	tmp := [][]byte{
		uint2Byte(b.Version),
		b.PrevHash,
		b.MerkleRoot,
		uint2Byte(b.TimeStamp),
		uint2Byte(b.Bits),
		uint2Byte(b.Nonce),
		b.Hash,
		b.Data,
	}

	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}
