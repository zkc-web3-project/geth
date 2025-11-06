package main

import (
    "context"
    "crypto/ecdsa"
    "fmt"
    "log"
    "math/big"

    "golang.org/x/crypto/sha3"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    client, err := ethclient.Dial("https://sepolia.infura.io/v3/aab47bf44d3b4988916ae6383cb5d490")
    if err != nil {
        log.Fatal(err)
    }

    privateKey, err := crypto.HexToECDSA("a430ee17e89cca50cee9ae5f8718875e8461e65a1960627f87f7f4b704a56167") //测试账户私钥
    if err != nil {
        log.Fatal(err)
    }

    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//打印公钥地址
	fmt.Println("fromAddress-公钥地址：", crypto.PubkeyToAddress(*publicKeyECDSA).Hex())
    if !ok {
        log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	fmt.Println("交易序号nonce:", nonce)
    if err != nil {
        log.Fatal(err)
    }

    value := big.NewInt(0) //当转账代币时转账金额设置为0
    gasPrice, err := client.SuggestGasPrice(context.Background())
	fmt.Println("获取建议的gas价格gasPrice:", gasPrice)
    if err != nil {
        log.Fatal(err)
    }

    // toAddress := common.HexToAddress("0xb1b29850e895add42661f51cf8cda44280404a3b") //接收方账户地址
    toAddress := common.HexToAddress("0x31dc153670ce03ad2628d23555cc161ac3f62b4e") //接收方账户地址
    tokenAddress := common.HexToAddress("0x654CBcFAA14608b8b428B2C108C8b281A7acEf24") //代币合约地址

	//构造转账数据
    transferFnSignature := []byte("transfer(address,uint256)") //转账函数的方法签名
    hash := sha3.NewLegacyKeccak256()
    hash.Write(transferFnSignature)
    methodID := hash.Sum(nil)[:4]
    fmt.Println(hexutil.Encode(methodID))
    paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32) //将地址和金额填充为32字节(ABI编码规范)
    fmt.Println(hexutil.Encode(paddedAddress))
    amount := new(big.Int)
    amount.SetString("5000000000000000000", 10) //转账的代币数量(这里表示5个代币，18位的精度)，第二个参数10代表采用十进制
    paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
    fmt.Println(hexutil.Encode(paddedAmount))
    var data []byte
    data = append(data, methodID...)  //函数选择器  注意！！！！这里的三个append都不能交换顺序
    data = append(data, paddedAddress...) //接收方地址
    data = append(data, paddedAmount...) //代币数量

    gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
        To:   &toAddress,
        Data: data,
    })
	//打印预估的gasLimit
	fmt.Println("预估的gasLimit:", gasLimit)
    if err != nil {
        log.Fatal(err)
    }
	gasLimit = 300000 //手动重置，预估的不准，可能造成交易失败
    tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    //打印链ID和交易签名
	fmt.Println("链ID:", chainID)
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println("交易签名:", signedTx)
    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("获取到交易hash:", signedTx.Hash().Hex())
}