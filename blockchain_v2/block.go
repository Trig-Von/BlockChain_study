package main

import (
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
	//b.setHash()
	pow := NewProofOfWork(&b)
	hash,nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce
	return &b
}

//func (b *Block) setHash() {
//
//	tmp := [][]byte{
//		uint2Byte(b.Version),
//		b.PrevHash,
//		b.MerkleRoot,
//		uint2Byte(b.TimeStamp),
//		uint2Byte(b.Bits),
//		uint2Byte(b.Nonce),
//		b.Hash,
//		b.Data,
//	}
//
//	data := bytes.Join(tmp, []byte{})
//
//	hash := sha256.Sum256(data)
//	b.Hash = hash[:]
//}
