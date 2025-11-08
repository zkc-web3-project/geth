package main

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zkc-web3-project/geth/queryERC20Balance/token" //这里导入abigen 生成的包
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/aab47bf44d3b4988916ae6383cb5d490")
	if err != nil {
		log.Fatal(err)
	}
	// Golem (GNT) Address
	tokenAddress := common.HexToAddress("0x654CBcFAA14608b8b428B2C108C8b281A7acEf24") //代币合约地址 MDT
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	// name, err := instance.Name(&bind.CallOpts{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// symbol, err := instance.Symbol(&bind.CallOpts{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// decimals, err := instance.Decimals(&bind.CallOpts{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	// fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	// fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"
	fmt.Printf("wei: %s\n", bal)
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(18))))
	fmt.Printf("balance: %f", value)
}
