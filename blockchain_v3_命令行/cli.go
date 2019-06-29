package main

import (
	"fmt"
	"os"
)

//处理用户输入命令，完成具体函数的调用

type CLI struct {
}

const Usage = `
usage :
	./blockchain create "创建区块链"
	./blockchain addBlock <需要写入的数据> "添加区块"
	./blockchain print "打印区块链"
`

func (cli *CLI) RUN() {
	cmds := os.Args

	if len(cmds) <2 {
		fmt.Println("输入参数无效，请检查！！")
		fmt.Println(Usage)
		return
	}

	switch cmds[1] {
	case "create":
		fmt.Println("创建区块被调用！")
		cli.createBlockChain()
	case "addBlock":
		if len(cmds) != 3 {
			fmt.Println("输入参数无效，请检查！！")
			return
		}
		data := cmds[2]
		cli.addBlock(data)
	case "print":
		fmt.Println("打印区块被调用！")
		cli.print()
	default:
		fmt.Println("输入参数无效，请检查！！")
		fmt.Println(Usage)
	}
}
