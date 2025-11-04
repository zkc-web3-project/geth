package main

import (
    "crypto/ecdsa"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
    "golang.org/x/crypto/sha3"
)

func main() {
    privateKey, err := crypto.GenerateKey() //生成随机私钥，32字节
    if err != nil {
        log.Fatal(err)
    }

    privateKeyBytes := crypto.FromECDSA(privateKey) //椭圆曲线数字签名算法
    fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 去掉'0x'
    publicKey := privateKey.Public()//获取私钥对应的公钥（公钥又私钥派生）
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
    fmt.Println("from pubKey:", hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04' 
    address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex() //从公钥计算地址
    fmt.Println(address)
    hash := sha3.NewLegacyKeccak256()  //创建一个Keccak256哈希实例
    hash.Write(publicKeyBytes[1:])
    fmt.Println("full:", hexutil.Encode(hash.Sum(nil)[:]))
    fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 原长32位，截去12位，保留后20位
}