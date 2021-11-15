package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"send-to/eth"
	"send-to/global"
)

func init() {
	flag.StringVar(&global.Node, "node", "https://bsc-dataseed.binance.org/", "node url")
	flag.IntVar(&global.MaxGoroutine, "workers", 10, "workers to deal block concurrently")
	flag.StringVar(&global.Address, "address", "0x10ED43C718714eb63d5aA57B78B54704E256024E", "address")
	flag.Uint64Var(&global.From, "from", 12676576, "from block")
	flag.Uint64Var(&global.To, "to", 12676586, "to block")
	flag.Parse()
}

func main() {
	logrus.Infof("get transactions to address %s from block %d to block %d using node %s with %d workers, start!",global.Address, global.From, global.To, global.Node, global.MaxGoroutine)
	eth.Init(global.Node)
	tmp := make([]*eth.Tx,0)
	txn := &tmp
	var err error
	if global.To != 0{
		txn, err = eth.Cli.SendTo(global.MaxGoroutine, common.HexToAddress(global.Address), global.From, global.To)
		if err != nil {
			return
		}
	}else{
		txn, err = eth.Cli.SendTo(global.MaxGoroutine, common.HexToAddress(global.Address), global.From)
		if err != nil {
			return
		}
	}
	for _,t:=range *txn{
		for i,v:= range *t.Transaction{
			fmt.Println(t.BlockNumber,v.Hash().String(),t.Status[i])
		}
	}
}