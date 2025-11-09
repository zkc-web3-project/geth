package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zkc-web3-project/geth/deployContract/store" //这里引入abigin生成的go绑定代码
)

const (
	contractAddr = "0x2d28a07a518DdFfB8819be5D2b7fE4303E050BB8"
)

func main() {
	//使用已部署的合约地址加载合约
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/aab47bf44d3b4988916ae6383cb5d490")
	if err != nil {
		log.Fatal(err)
	}
	storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	_ = storeContract
	fmt.Println("合约加载成功")
}
