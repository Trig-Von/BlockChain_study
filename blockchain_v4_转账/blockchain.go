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

	if isFileExist(blockchainDBFile) {
		fmt.Println("区块链文件已经开始存在!")
		return nil
	}
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

			coinbase := NewCoinbaseTx("钟本聪", genesisInfo)

			txs := []*Transaction{coinbase}

			genesisBlock := NewBlock(txs, nil)
			//创建BlockChain ，同时添加一个创世块
			//genesisBlock := NewBlock(genesisInfo, nil)
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())

			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)

		}
		return nil
	})
	return err
}

func GetBlockChainInstance() (*BlockChain, error) {

	if isFileExist(blockchainDBFile) == false {
		return nil, errors.New("区块链文件不存在,请先创建!!")
	}
	var lastHash []byte //内存中最后一个区块的 哈西值
	//如果区块链不存在，则创建，同时返回blockchain实例
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
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

func (bc *BlockChain) AddBlock(txs []*Transaction) error {

	lastBlockHash := bc.tail
	//创建区块
	NewBlock := NewBlock(txs, lastBlockHash)
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

func (it *Iterator) Next() (block *Block) {

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
		fmt.Println("iterator next err:", err)
		return nil
	}
	return
}

type UTXOInfo struct {
	Txid []byte

	Index int64

	TXOutput
}
//获取指定地址的金额,实现遍历账本的通用函数
//给定一个地址，返回所有的utxo
func (bc *BlockChain) FindMyUTXO(address string) []UTXOInfo {
	//存储所有和目标地址相关的utxo集合
	var utxoInfos []UTXOInfo

	//定义一个存放已经消耗过的所有的utxos的集合(跟指定地址相关的)
	spentUtxos := make(map[string][]int)

	it := bc.NewIterator()
	for {

		block := it.Next()

		for _, tx := range block.Transactions {
		LABEL:
			for outputIndex, output := range tx.TXOutputs {
				fmt.Println("outputIndex :", outputIndex)
				if output.ScriptPubk == address {
					//开始过滤
					//当前交易id
					currentTxid := string(tx.TXID)
					//去spentUtxos中查看
					indexArray := spentUtxos[currentTxid]

					if len(indexArray) != 0 {
						for _, spentIndex := range indexArray {
							if outputIndex == spentIndex {
								continue LABEL
							}

						}
					}
					//utxos = append(utxos, output)
					utxoinfo := UTXOInfo{tx.TXID,int64(outputIndex),output}
					utxoInfos = append(utxoInfos,utxoinfo)
				}
			}

			if tx.isCoinbaseTX() {
				fmt.Println("发现挖矿交易,无需遍历inputs")
				continue
			}
			for _,input := range tx.TXInputs{
				if input.ScriptSig == address {
					spentKey := string(input.Txid)

					spentUtxos[spentKey] = append(spentUtxos[spentKey])
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return utxoInfos
}

func (bc *BlockChain)findNeedUTXO(from string,amount float64) (map[string][]int64,float64) {
	var retMap = make(map[string][]int64)
	var retValue float64
	//遍历账本,找到所有utxo
	utxoInfos := bc.FindMyUTXO(from)
	//遍历utxo,统计当前总额,与amount比较
	for _,utxoinfo := range utxoInfos{
		retValue += utxoinfo.Value

		//统计将要消耗的utxo
		key := string(utxoinfo.Txid)
		retMap[key] = append(retMap[key],utxoinfo.Index)

		//如果大于等于amount直接返回
		if retValue >= amount{
			break
		}
	}
	return retMap, retValue

}