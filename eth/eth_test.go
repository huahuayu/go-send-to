package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

var url = "https://bsc-dataseed2.binance.org"

func TestClient_SendTo(t *testing.T) {
	Init(url)
	txs, err := Cli.SendTo(100, common.HexToAddress("0x10ED43C718714eb63d5aA57B78B54704E256024E"), 12676566, 12676576)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(*txs))
}
