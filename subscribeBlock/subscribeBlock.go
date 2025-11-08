package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	//使用webbsockets连接
	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/aab47bf44d3b4988916ae6383cb5d490")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("开始监听Sepolia新区块...")

	headers := make(chan *types.Header) //创建一个通道用于接收新区块的头信息
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("订阅成功，等待新区块...")

	for {
		select { //监听多个通道的数据 区块错误事件和新区块头事件
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Printf("监听到新区块头hash:%s\n", header.Hash().Hex())
			//根据区块头hash获取区块信息
			//这里需要等待区块信息完全同步到节点上，阻塞200秒
			time.Sleep(500 * time.Second)
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
				continue
			}

			fmt.Printf("区块hash:%s\n", block.Hash().Hex())
			fmt.Printf("区块时间:%d秒\n", block.NumberU64())
			fmt.Printf("区块nonce:%d\n", block.Nonce())
			fmt.Printf("区块交易数量:%d\n", len(block.Transactions()))
		}
	}
}
