package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func uint2Byte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println("binary write ", err)
		return nil
	}

	return buffer.Bytes()
}

//判断文件是否存在
func isFileExist(filename string) bool {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}
	return true
}
