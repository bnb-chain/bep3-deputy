// +build integration

package integration_test

import (
	"errors"
	"math/big"
	"os"
	"testing"
	"time"

	bnbExe "github.com/binance-chain/bep3-deputy/executor/bnb"
	kavaExe "github.com/binance-chain/bep3-deputy/executor/kava"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
	"github.com/binance-chain/go-sdk/common/types"
	bnbKeys "github.com/binance-chain/go-sdk/keys"
	ec "github.com/ethereum/go-ethereum/common"
	sdk "github.com/kava-labs/cosmos-sdk/types"
	"github.com/kava-labs/go-sdk/client"
	"github.com/kava-labs/go-sdk/kava"
	"github.com/kava-labs/go-sdk/kava/bep3"
	kavaKeys "github.com/kava-labs/go-sdk/keys"
	"github.com/stretchr/testify/require"
)

const (
	// these are the same as the menmonics in the chains and deputy configs
	bnbDeputyMnemonic    = "clinic soap symptom alter mango orient punch table seek among broken bundle best dune hurt predict liquid subject silver once kick metal okay moment"
	bnbTestUserMnemonic  = "then nuclear favorite advance plate glare shallow enhance replace embody list dose quick scale service sentence hover announce advance nephew phrase order useful this"
	kavaDeputyMnemonic   = "equip town gesture square tomorrow volume nephew minute witness beef rich gadget actress egg sing secret pole winter alarm law today check violin uncover"
	kavaTestUserMnemonic = "very health column only surface project output absent outdoor siren reject era legend legal twelve setup roast lion rare tunnel devote style random food"

	bnbHTLTFee = 37500
)

var bnbDeputyAddr, bnbTestUserAddr, kavaDeputyAddr, kavaTestUserAddr string

func TestMain(m *testing.M) {
	kavaConfig := sdk.GetConfig()
	kava.SetBech32AddressPrefixes(kavaConfig)
	kavaConfig.Seal()

	bnbManager, err := bnbKeys.NewMnemonicKeyManager(bnbDeputyMnemonic)
	if err != nil {
		panic(err.Error())
	}
	bnbDeputyAddr = bnbManager.GetAddr().String()
	bnbManager, err = bnbKeys.NewMnemonicKeyManager(bnbTestUserMnemonic)
	if err != nil {
		panic(err.Error())
	}
	bnbTestUserAddr = bnbManager.GetAddr().String()
	kavaManager, err := kavaKeys.NewMnemonicKeyManager(kavaDeputyMnemonic, kava.Bip44CoinType)
	if err != nil {
		panic(err.Error())
	}
	kavaDeputyAddr = kavaManager.GetAddr().String()
	kavaManager, err = kavaKeys.NewMnemonicKeyManager(kavaTestUserMnemonic, kava.Bip44CoinType)
	if err != nil {
		panic(err.Error())
	}
	kavaTestUserAddr = kavaManager.GetAddr().String()

	os.Exit(m.Run())
}

func TestBnbToKavaSwap(t *testing.T) {

	// 1) setup executors to send txs

	config := util.ParseConfigFromFile("deputy/config.json")
	bnbConfig := config.BnbConfig
	bnbConfig.RpcAddr = "tcp://localhost:26658"
	bnbConfig.Mnemonic = bnbTestUserMnemonic

	bnbExecutor := bnbExe.NewExecutor(types.ProdNetwork, bnbConfig)

	kavaConfig := config.KavaConfig
	kavaConfig.RpcAddr = "tcp://localhost:26657"
	kavaConfig.Mnemonic = kavaDeputyMnemonic

	kavaExecutor := kavaExe.NewExecutor(client.LocalNetwork, kavaConfig)

	// 2) Cache account balances to compare against

	bnbTestUserBalance, err := bnbExecutor.GetBalance(bnbTestUserAddr)
	require.NoError(t, err)
	kavaTestUserBalance, err := kavaExecutor.GetBalance(kavaTestUserAddr)
	require.NoError(t, err)

	swapAmount := big.NewInt(100_000_000)

	// 3) Send bnb swap

	rndNum, err := bep3.GenerateSecureRandomNumber()
	require.NoError(t, err)
	timestamp := time.Now().Unix()
	rndHash := ec.BytesToHash(bep3.CalculateRandomHash(rndNum, timestamp))
	htltTxHash, cmnErr := bnbExecutor.HTLT(
		rndHash,
		timestamp,
		20000,
		bnbDeputyAddr,
		kavaDeputyAddr,
		kavaTestUserAddr,
		swapAmount,
	)
	// Note: this cannot use require.NoError as that wraps the commonError in an error interface.
	// This makes err != nil (despite the underlying value being a nil pointer) and err.Error() also panics.
	require.Nil(t, cmnErr)

	err = wait(8*time.Second, func() (bool, error) {
		return bnbExecutor.GetSentTxStatus(htltTxHash) == store.TxSentStatusSuccess, nil
	})
	require.NoError(t, err)
	t.Log("bnb htlt created")

	// 4) Wait until deputy relays swap to kava

	kavaSwapIDBz, err := kavaExecutor.CalcSwapId(rndHash, kavaDeputyAddr, bnbTestUserAddr)
	require.NoError(t, err)
	kavaSwapID := ec.BytesToHash(kavaSwapIDBz)

	err = wait(15*time.Second, func() (bool, error) {
		t.Log("waiting...")
		return kavaExecutor.HasSwap(kavaSwapID)
	})
	require.NoError(t, err)
	t.Log("swap created on kava by deputy")

	// 5) Send claim on kava

	claimTxHash, cmnErr := kavaExecutor.Claim(kavaSwapID, ec.BytesToHash(rndNum))
	require.Nil(t, cmnErr)

	err = wait(8*time.Second, func() (bool, error) {
		return kavaExecutor.GetSentTxStatus(claimTxHash) == store.TxSentStatusSuccess, nil
	})
	require.NoError(t, err)
	t.Log("kava htlt claimed")

	// 6) Wait until deputy relays claim to bnb

	bnbSwapIDBz, err := bnbExecutor.CalcSwapId(rndHash, bnbTestUserAddr, kavaDeputyAddr)
	require.NoError(t, err)
	bnbSwapID := ec.BytesToHash(bnbSwapIDBz)

	wait(10*time.Second, func() (bool, error) {
		// check the deputy has relayed the claim by checking the status of the swap
		// once claimed it is no longer claimable, if it timesout it will become refundable
		claimable, err := bnbExecutor.Claimable(bnbSwapID)
		if err != nil {
			return false, err
		}
		refundable, err := bnbExecutor.Refundable(bnbSwapID)
		if err != nil {
			return false, err
		}
		return !(claimable || refundable), nil
	})
	t.Log("bnb htlt claimed by deputy")

	// 7) Check balances

	bnbTestUserBalanceFinal, err := bnbExecutor.GetBalance(bnbTestUserAddr)
	require.NoError(t, err)
	kavaTestUserBalanceFinal, err := kavaExecutor.GetBalance(kavaTestUserAddr)
	require.NoError(t, err)

	expectedBnbBalance := big.NewInt(0)
	expectedBnbBalance.Sub(bnbTestUserBalance, swapAmount).Sub(expectedBnbBalance, big.NewInt(bnbHTLTFee))
	require.Zerof(
		t,
		expectedBnbBalance.Cmp(bnbTestUserBalanceFinal),
		"expected: %, actual: %s",
		expectedBnbBalance,
		bnbTestUserBalanceFinal,
	)

	var swapAmountKava = &big.Int{}
	swapAmountKava.Sub(swapAmount, config.ChainConfig.BnbFixedFee)
	require.Zero(t, config.ChainConfig.BnbRatio.Cmp(big.NewFloat(1)), "test does not support ratio conversions other than 1")
	require.Zero(
		t,
		new(big.Int).Add(kavaTestUserBalance, swapAmountKava).Cmp(kavaTestUserBalanceFinal),
	)
}

func TestKavaToBnbSwap(t *testing.T) {

	// 1) setup executors to send txs

	config := util.ParseConfigFromFile("deputy/config.json")

	kavaConfig := config.KavaConfig
	kavaConfig.RpcAddr = "tcp://localhost:26657"
	kavaConfig.Mnemonic = kavaTestUserMnemonic
	kavaExecutor := kavaExe.NewExecutor(client.LocalNetwork, kavaConfig)

	bnbConfig := config.BnbConfig
	bnbConfig.RpcAddr = "tcp://localhost:26658"
	bnbConfig.Mnemonic = bnbDeputyMnemonic
	bnbExecutor := bnbExe.NewExecutor(types.ProdNetwork, bnbConfig)

	// 2) Cache account balances to compare against

	bnbTestUserBalance, err := bnbExecutor.GetBalance(bnbTestUserAddr)
	require.NoError(t, err)
	kavaTestUserBalance, err := kavaExecutor.GetBalance(kavaTestUserAddr)
	require.NoError(t, err)

	swapAmount := big.NewInt(99_000_000)

	// 3) Send kava swap

	rndNum, err := bep3.GenerateSecureRandomNumber()
	require.NoError(t, err)
	timestamp := time.Now().Unix()
	rndHash := ec.BytesToHash(bep3.CalculateRandomHash(rndNum, timestamp))
	htltTxHash, cmnErr := kavaExecutor.HTLT(
		rndHash,
		timestamp,
		230,
		kavaDeputyAddr,
		bnbDeputyAddr,
		bnbTestUserAddr,
		swapAmount,
	)
	// Note: this cannot use require.NoError as that wraps the commonError in an error interface.
	// This makes err != nil (despite the underlying value being a nil pointer) and err.Error() also panics.
	require.Nil(t, cmnErr)

	err = wait(8*time.Second, func() (bool, error) {
		return kavaExecutor.GetSentTxStatus(htltTxHash) == store.TxSentStatusSuccess, nil
	})
	require.NoError(t, err)
	t.Log("kava htlt created")

	// 4) Wait until deputy relays swap to bnb

	bnbSwapIDBz, err := bnbExecutor.CalcSwapId(rndHash, bnbDeputyAddr, kavaTestUserAddr)
	require.NoError(t, err)
	bnbSwapID := ec.BytesToHash(bnbSwapIDBz)

	err = wait(15*time.Second, func() (bool, error) {
		t.Log("waiting...")
		return bnbExecutor.HasSwap(bnbSwapID)
	})
	require.NoError(t, err)
	t.Log("swap created on bnb by deputy")

	// 5) Send claim on bnb

	claimTxHash, cmnErr := bnbExecutor.Claim(bnbSwapID, ec.BytesToHash(rndNum))
	require.Nil(t, cmnErr)

	err = wait(8*time.Second, func() (bool, error) {
		return bnbExecutor.GetSentTxStatus(claimTxHash) == store.TxSentStatusSuccess, nil
	})
	require.NoError(t, err)
	t.Log("bnb htlt claimed")

	// 6) Wait until deputy relays claim to kava

	kavaSwapIDBz, err := kavaExecutor.CalcSwapId(rndHash, kavaTestUserAddr, bnbDeputyAddr)
	require.NoError(t, err)
	kavaSwapID := ec.BytesToHash(kavaSwapIDBz)

	wait(10*time.Second, func() (bool, error) {
		// check the deputy has relayed the claim by checking the status of the swap
		// once claimed it is no longer claimable, if it timesout it will become refundable
		claimable, err := kavaExecutor.Claimable(kavaSwapID)
		if err != nil {
			return false, err
		}
		refundable, err := kavaExecutor.Refundable(kavaSwapID)
		if err != nil {
			return false, err
		}
		return !(claimable || refundable), nil
	})
	t.Log("kava htlt claimed by deputy")

	// 7) Check balances

	bnbTestUserBalanceFinal, err := bnbExecutor.GetBalance(bnbTestUserAddr)
	require.NoError(t, err)
	kavaTestUserBalanceFinal, err := kavaExecutor.GetBalance(kavaTestUserAddr)
	require.NoError(t, err)

	expectedKavaBalance := big.NewInt(0)
	expectedKavaBalance.Sub(kavaTestUserBalance, swapAmount)
	require.Zerof(
		t,
		expectedKavaBalance.Cmp(kavaTestUserBalanceFinal),
		"expected: %, actual: %s",
		expectedKavaBalance,
		kavaTestUserBalanceFinal,
	)

	var swapAmountBnb = &big.Int{}
	swapAmountBnb.Sub(swapAmount, config.ChainConfig.OtherChainFixedFee)
	require.Zero(t, config.ChainConfig.OtherChainRatio.Cmp(big.NewFloat(1)), "test does not support ratio conversions other than 1")
	require.Zero(
		t,
		new(big.Int).Add(bnbTestUserBalance, swapAmountBnb).Cmp(bnbTestUserBalanceFinal),
	)
}

func wait(timeout time.Duration, shouldStop func() (bool, error)) error {
	endTime := time.Now().Add(timeout)

	for {
		stop, err := shouldStop()
		switch {
		case err != nil || stop:
			return err
		case time.Now().After(endTime):
			return errors.New("waiting timed out")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
