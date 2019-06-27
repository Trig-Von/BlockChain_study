package main

//定义区块连结构(使用数组模拟区块连)
type BlockChain struct {
	Blocks []*Block
}

//创世语
const genesisInfo = "星星之火"

//提供一个创建区块连的方法
func NewBlockChain() *BlockChain {
	//创建BlockChain，同时添加一个创始块
	genesisBlock := NewBlock(genesisInfo, nil)

	bc := BlockChain{
		Blocks: []*Block{genesisBlock},
	}

	return &bc
}

//向区块连添加区块的方法

func (bc *BlockChain)AddBlock(data string)  {

	lastBlock := bc.Blocks[len(bc.Blocks)-1]

	prevHash := lastBlock.Hash

	newBlock := NewBlock(data,prevHash)

	bc.Blocks = append(bc.Blocks,newBlock)
}