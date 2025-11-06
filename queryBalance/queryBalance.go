package main

import (
    "context"
    "fmt"
    "log"
    "math"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    client, err := ethclient.Dial("https://sepolia.infura.io/v3/aab47bf44d3b4988916ae6383cb5d490")
    if err != nil {
        log.Fatal(err)
    }

    account := common.HexToAddress("0xb1b29850e895add42661f51cf8cda44280404a3b")
    balance, err := client.BalanceAt(context.Background(), account, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("当前账户余额:", balance)
    blockNumber := big.NewInt(9415486)
    balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("9415486块的账户余额:",balanceAt)
    fbalance := new(big.Float)
    fbalance.SetString(balanceAt.String())
    ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
    fmt.Println(ethValue)
    pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
    fmt.Println(pendingBalance)
}