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

// func TestKavaToBnbSwap(t *testing.T) {
// 	/*
// 		create kava swap msg
// 		send tx
// 		wait until accepted
// 		wait until relayed
// 		create claim msg
// 		send tx
// 		wait until accepted
// 		wait until relayed
// 		check balances
// 	*/

// 	cdc := kava.MakeCodec()
// 	kavaClient := client.NewKavaClient(cdc, kavaTestUserMnemonic, kava.Bip44CoinType, "tcp://localhost:26657", client.LocalNetwork)

// 	deputyAccAddr, err := sdk.AccAddressFromBech32(kavaDeputyAddr)
// 	require.NoError(t, err)

// 	rndNum, err := bep3.GenerateSecureRandomNumber()
// 	require.NoError(t, err)
// 	timestamp := time.Now().Unix()
// 	hash := bep3.CalculateRandomHash(rndNum, timestamp)

// 	msg := bep3.NewMsgCreateAtomicSwap(
// 		kavaClient.Keybase.GetAddr(),
// 		deputyAccAddr,
// 		bnbTestUserAddr,
// 		bnbDeputyAddr,
// 		hash,
// 		timestamp,
// 		sdk.NewCoins(sdk.NewInt64Coin("bnb", 100_000_000)),
// 		250,
// 	)
// 	res, err := kavaClient.Broadcast(msg, client.Commit) // TODO use client.Sync and polling?
// 	require.NoError(t, err)
// 	require.Equal(t, uint32(0), res.Code, res.Log)
// }

func TestBnbToKavaSwap(t *testing.T) {
	config := util.ParseConfigFromFile("deputy/config.json")
	bnbConfig := config.BnbConfig
	bnbConfig.RpcAddr = "tcp://localhost:26658"
	bnbConfig.Mnemonic = bnbTestUserMnemonic

	bnbExecutor := bnbExe.NewExecutor(types.ProdNetwork, bnbConfig)

	rndNum, err := bep3.GenerateSecureRandomNumber()
	require.NoError(t, err)
	timestamp := time.Now().Unix()
	rndHashSlice := bep3.CalculateRandomHash(rndNum, timestamp)
	var rndHash ec.Hash
	copy(rndHash[:], rndHashSlice)
	txHash, cmnErr := bnbExecutor.HTLT(
		rndHash,
		timestamp,
		20000,
		bnbDeputyAddr,
		kavaDeputyAddr,
		kavaTestUserAddr,
		big.NewInt(100_000_000),
	)
	// Note this cannot use require.NoError as that wraps the commonError in an error interface.
	// This makes err != nil (despite the underlying value being a nil pointer) and err.Error() also panics.
	require.Nil(t, cmnErr)

	err = wait(8*time.Second, func() (bool, error) {
		return bnbExecutor.GetSentTxStatus(txHash) == store.TxSentStatusSuccess, nil
	})
	require.NoError(t, err)
	t.Log("bnb htlt created")

	// kava stuff
	kavaConfig := config.KavaConfig
	kavaConfig.RpcAddr = "tcp://localhost:26657"
	kavaConfig.Mnemonic = kavaDeputyMnemonic // TODO make this the test user and swap out deputy address to make sigs valid?

	kavaExecutor := kavaExe.NewExecutor(client.LocalNetwork, kavaConfig)

	idBytes, err := kavaExecutor.CalcSwapId(rndHash, kavaDeputyAddr, bnbTestUserAddr)
	var id ec.Hash
	copy(id[:], idBytes)

	err = wait(15*time.Second, func() (bool, error) {
		t.Log("waiting...")
		return kavaExecutor.HasSwap(id)
	})
	require.NoError(t, err)
	t.Log("swap created on kava by deputy")

	txHash, cmnErr = kavaExecutor.Claim(id, byteSliceToArray(rndNum))
	if cmnErr != nil { // TODO
		t.Log(cmnErr.Error())
		t.FailNow()
	}
	err = wait(8*time.Second, func() (bool, error) {
		return kavaExecutor.GetSentTxStatus(txHash) == store.TxSentStatusSuccess, nil
	})
	require.NoError(t, err)
	t.Log("kava htlt claimed")

	bz, err := bnbExecutor.CalcSwapId(rndHash, bnbTestUserAddr, kavaDeputyAddr)
	require.NoError(t, err)

	wait(10*time.Second, func() (bool, error) {
		claimable, err := bnbExecutor.Claimable(ec.BytesToHash(bz))
		return !claimable, err
	})
	t.Log("bnb htlt claimed by deputy") // TODO it could also time out sort of

	// TODO check balances n stuff
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

// TODO ec.BytesToHash
func byteSliceToArray(bz []byte) [32]byte {
	var array [32]byte
	copy(array[:], bz)
	return array
}
