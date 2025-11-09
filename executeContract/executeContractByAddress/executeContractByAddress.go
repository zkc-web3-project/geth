package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zkc-web3-project/geth/deployContract/store" //这里引入abigin生成的go绑定代码
)

const (
	contractAddr = "0x1Fc54de16F09Da082f3b546323036Ae86a35F956"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/aab47bf44d3b4988916ae6383cb5d490")
	if err != nil {
		log.Fatal(err)
	}
	//创建合约实例
	storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	//获取私钥
	privateKey, err := crypto.HexToECDSA("a430ee17e89cca50cee9ae5f8718875e8461e65a1960627f87f7f4b704a56167")
	if err != nil {
		log.Fatal(err)
	}
	//组装函数入参
	var key [32]byte
	var value [32]byte

	copy(key[:], []byte("this_is_data_key"))
	copy(value[:], []byte("this_is_data_value"))

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) //sepolia固定链id
	if err != nil {
		log.Fatal(err)
	}
	tx, err := storeContract.SetItem(opt, key, value) //调用合约内的方法
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash:", tx.Hash().Hex())

	//查询交易数据
	callOpt := &bind.CallOpts{Context: context.Background()}
	valueInContract, err := storeContract.Items(callOpt, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("is value saving in contract equals to origin value:", valueInContract == value)
}
