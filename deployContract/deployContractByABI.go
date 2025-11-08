package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zkc-web3-project/geth/deployContract/store" //这里引入abigin生成的go绑定代码
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/aab47bf44d3b4988916ae6383cb5d490")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("a430ee17e89cca50cee9ae5f8718875e8461e65a1960627f87f7f4b704a56167")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("部署合约的地址:%s\n", fromAddress.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress) //nonce:交易序号，递增，防止双花和重放攻击
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId) //使用私钥和链id创建交易对象
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) //部署合约时设置ETH数量为0
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	input := "1.0"                                                       //构造构造函数的参数(版本)
	address, tx, instance, err := store.DeployStore(auth, client, input) //调用go绑定代码的部署函数
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("部署的合约地址为:%s\n", address.Hex())

	fmt.Printf("部署合约的交易hash为:%s\n", tx.Hash().Hex())

	_ = instance //后续可以继续调用这个合约实例的函数
}
