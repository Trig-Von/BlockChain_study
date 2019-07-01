package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type wallet struct {
	Prikey *ecdsa.PrivateKey

	PubKey []byte
}

func newWalletKeyPair() *wallet {
	curve := elliptic.P256()

	//创建私钥
	priKey ,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		fmt.Println("ecdsa.GenerateKey err:",err)
		return nil
	}

	//获取公钥
	pubKeyRaw :=priKey.PublicKey

	//将公钥X,Y拼接到一起
	pubKey := append(pubKeyRaw.X.Bytes(),pubKeyRaw.Y.Bytes()...)
	//创建wallet结构返回
	wallet := wallet{priKey,pubKey}
	return &wallet
}

func (w *wallet)getAddress() string {
	//公钥
	pubKey := w.PubKey
	hash1 := sha256.Sum256(pubKey)
	//hash160处理
	hasher := ripemd160.New()
	hasher.Write(hash1[:])

	//公钥哈希,锁定output时就是使用该值
	pubKeyHash := hasher.Sum(nil)
	//拼接version和公钥哈希,得到21字节数据
	payload := append([]byte{byte(0x00)},pubKeyHash...)
	//生成4字节的校验码
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	//4字节checksum
	checksum := second[0:4]

	//25字节数据
	payload = append(payload,checksum...)
	address := base58.Encode(payload)
	return address
}











