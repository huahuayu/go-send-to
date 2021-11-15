package eth

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"sync"
	"time"
)

var (
	Cli *Client
)

type (
	Client struct {
		EthClient *ethclient.Client
	}

	Tx struct {
		BlockNumber uint64
		Transaction *[]*types.Transaction
		Status []bool
	}
)

func Init(url string) {
	var err error
	Cli, err = NewClient(url)
	if err != nil {
		panic(err)
	}
}

func NewClient(url string) (*Client, error) {
	client := new(Client)
	if ethClient, err := ethclient.Dial(url); err != nil {
		panic(err)
		return nil, err
	} else {
		client.EthClient = ethClient
	}
	return client, nil
}

func (client *Client) SendTo(concurrentBlock int, address common.Address, from uint64, to ...uint64) (*[]*Tx, error) {
	var toBlock uint64
	ctx := context.Background()
	if len(to) > 0 {
		toBlock = to[0]
	} else {
		number, err := client.EthClient.BlockNumber(ctx)
		if err != nil {
			return nil, err
		}
		toBlock = number
	}
	if from > toBlock {
		return nil, errors.New("from should less or equal to to")
	}
	var(
		txs = make([]*Tx, 0)
		mu = sync.RWMutex{}
		wg = &sync.WaitGroup{}
		maxGoroutines = concurrentBlock
		guard = make(chan struct{}, maxGoroutines)
	)
	t0 := time.Now()
	totalBlocks := int(toBlock - from + 1)
	wg.Add(totalBlocks)
	for i := from; i <= toBlock; i++ {
		guard <- struct{}{}
		txns := make([]*types.Transaction,0)
		status := make([]bool,0)
		go func(blockNumber uint64, txns []*types.Transaction, status []bool) {
			defer wg.Done()
			block, err := client.EthClient.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
			if err != nil {
				logrus.Error("blockByNumber:", err)
			}
			if block == nil{
				return
			}
			transactions := block.Transactions()
			txMu := sync.RWMutex{}
			txWg := &sync.WaitGroup{}
			txWg.Add(len(transactions))
			for _, t := range transactions {
				go func(transaction *types.Transaction, txns *[]*types.Transaction, status *[]bool) {
					defer txWg.Done()
					if transaction.To() != nil {
						if strings.EqualFold(transaction.To().String(), address.String()) {
							txMu.Lock()
							defer txMu.Unlock()
							*txns = append(*txns, transaction)
							txStatus, err := client.TxStatus(transaction.Hash())
							if err != nil {
								logrus.Error("TxStatus: ", err)
								return
							}
							*status = append(*status,txStatus)
						}
					}
				}(t, &txns, &status)
			}
			txWg.Wait()
			if len(txns) > 0{
				mu.Lock()
				txs = append(txs, &Tx{
					BlockNumber: blockNumber,
					Transaction: &txns,
					Status: status,
				})
				//fmt.Println(blockNumber,len(txns),txns[0].Hash().String())
				defer mu.Unlock()
			}
			<-guard
		}(i,txns,status)
	}
	wg.Wait()
	duration := time.Since(t0)
	logrus.Infof("range %d blocks in %v, in average %dms per block",totalBlocks, duration, duration.Milliseconds()/int64(totalBlocks))
	return &txs, nil
}

func (client *Client) TxStatus(txHash common.Hash)(success bool,err error){
	receipt, err := client.EthClient.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return false, err
	}
	if receipt.Status == types.ReceiptStatusSuccessful{
		return true, nil
	}else{
		return false, nil
	}
}
