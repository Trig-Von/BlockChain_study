package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

//定义区块连结构(使用数组模拟区块连)
type BlockChain struct {
	db   *bolt.DB //用于存储数据
	tail []byte   //最后一个区块的哈西值
}

//创世语
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
const blockchainDBFile = "blockchain.db"
const bucketBlock = "bucketBlock"
const lastBlockHashKey = "lastBlockHashKey"

func CreateBlockChain() error {
	//区块链不存在，创建
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
	if err != nil {
		return nil
	}

	defer db.Close()
	//开始创建
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在
		if bucket == nil {
			//创建bucket
			bucket, err := tx.CreateBucket([]byte(bucketBlock))
			if err != nil {
				return err
			}
			//创建BlockChain ，同时添加一个创世块
			genesisBlock := NewBlock(genesisInfo, nil)
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())

			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)

		}
		return nil
	})
	return err
}

func GetBlockChainInstance() (*BlockChain, error) {

	var lastHash []byte //内存中最后一个区块的 哈西值
	//如果区块链不存在，则创建，同时返回blockchain实例
	db, err := bolt.Open(blockchainDBFile, 0400, nil)
	if err != nil {
		return nil, err
	}

	//如果区块链存在，则直接返回blockchain实例

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		if bucket == nil {
			return errors.New("bucket不应为nil")
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})

	bc := BlockChain{db, lastHash}
	return &bc, nil

}

/*//提供一个创建区块连的方法
func NewBlockChain() *BlockChain {
	//创建BlockChain，同时添加一个创始块
	genesisBlock := NewBlock(genesisInfo, nil)

	bc := BlockChain{
		Blocks: []*Block{genesisBlock},
	}

	return &bc
}*/

//向区块连添加区块的方法

func (bc *BlockChain) AddBlock(data string) error {

	lastBlockHash := bc.tail
	//创建区块
	NewBlock := NewBlock(data, lastBlockHash)
	//写入数据库
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("AddBlock时BUcket 不应为空")
		}

		bucket.Put(NewBlock.Hash, NewBlock.Serialize())
		bucket.Put([]byte(lastBlockHashKey), NewBlock.Hash)

		bc.tail = NewBlock.Hash
		return nil
	})
	return err
}

//定义迭代器
type Iterator struct {
	db          *bolt.DB
	currentHash []byte
}

func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}

	return &it
}

func (it *Iterator)Next() (block *Block)  {

	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("Iterator Next 时bucket 不应为nil")
		}

		blockTmpInfo := bucket.Get(it.currentHash)
		block = Deserialize(blockTmpInfo)
		it.currentHash = block.PrevHash
		return nil
	})

	if err != nil {
		fmt.Println("iterator next err:",err)
		return nil
	}
	return
}
