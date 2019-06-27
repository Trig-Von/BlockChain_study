package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func uint2Byte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer,binary.LittleEndian,&num)
	if err != nil {
		fmt.Println("binary write ",err)
		return nil
	}

	return buffer.Bytes()
}