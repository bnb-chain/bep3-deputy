// +build integration

package integration_test

import (
	"math/big"
	"testing"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/util"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentBnbToKavaSwaps(t *testing.T) {

	// 1) setup executors

	config := util.ParseConfigFromFile("deputy/config.json")

	var senderExecutors []common.Executor
	for i := range bnbUserMnemonics {
		senderExecutors = append(senderExecutors, setupUserExecutorBnb(*config.BnbConfig, bnbUserMnemonics[i]))

	}
	senderAddrs := bnbUserAddrs

	var receiverExecutors []common.Executor
	for i := range kavaUserMnemonics {
		receiverExecutors = append(receiverExecutors, setupUserExecutorKava(*config.KavaConfig, kavaUserMnemonics[i]))
	}
	receiverAddrs := kavaUserAddrs

	// 2) Send swaps from bnb to kava

	swapAmount := big.NewInt(100_000_000)
	type result struct {
		id  int
		err error
	}
	results := make(chan result)
	for i := range senderExecutors {
		go func(i int) {
			t.Logf("sending swap %d\n", i)
			err := sendCompleteSwap(t, senderExecutors[i], receiverExecutors[i], senderAddrs[i], receiverAddrs[i], swapAmount, bnbDeputyAddr, kavaDeputyAddr, 20000)
			results <- result{i, err}
		}(i)
	}

	// 3) Check results

	for range senderExecutors {
		r := <-results
		t.Logf("swap %d done, err: %v\n", r.id, r.err)
		assert.NoErrorf(t, r.err, "swap %d returned error", r.id)
	}
}
