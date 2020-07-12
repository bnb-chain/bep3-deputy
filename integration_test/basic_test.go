// +build integration

package integration_test

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/binance-chain/go-sdk/common/types"
	bnbKeys "github.com/binance-chain/go-sdk/keys"
	ec "github.com/ethereum/go-ethereum/common"
	sdk "github.com/kava-labs/cosmos-sdk/types"
	"github.com/kava-labs/go-sdk/client"
	"github.com/kava-labs/go-sdk/kava"
	"github.com/kava-labs/go-sdk/kava/bep3"
	kavaKeys "github.com/kava-labs/go-sdk/keys"
	"github.com/stretchr/testify/require"

	"github.com/binance-chain/bep3-deputy/common"
	bnbExe "github.com/binance-chain/bep3-deputy/executor/bnb"
	kavaExe "github.com/binance-chain/bep3-deputy/executor/kava"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

const bnbHTLTFee = 37500

var (
	// these are the same as the menmonics in the chains and deputy configs
	bnbDeputyMnemonic  = "clinic soap symptom alter mango orient punch table seek among broken bundle best dune hurt predict liquid subject silver once kick metal okay moment"
	kavaDeputyMnemonic = "equip town gesture square tomorrow volume nephew minute witness beef rich gadget actress egg sing secret pole winter alarm law today check violin uncover"

	kavaUserMnemonics = []string{
		"very health column only surface project output absent outdoor siren reject era legend legal twelve setup roast lion rare tunnel devote style random food",
		"curtain camp spoil tiny vehicle pottery deer corn truly banner salmon lift yard throw open move state lamp van sign glow glue shrug faith",
		"desert october mammal tuition illness album engine solid enjoy harvest symptom rely camera unable okay avocado actual oppose remember lady dove canal argue cave",
		"profit law bounce grunt earth ice share skill valve awful around shoot include kite lecture also smooth ball vintage snake embark brief ill gather",
		"census museum crew rude tower vapor mule rib weasel faith page cushion rain inherit much cram that blanket occur region track hub zero topple",
		"flavor print loyal canyon expand salmon century field say frequent human dinosaur frame claim bridge affair web way direct win become merry crash frequent",
	}
	bnbUserMnemonics = []string{
		"then nuclear favorite advance plate glare shallow enhance replace embody list dose quick scale service sentence hover announce advance nephew phrase order useful this",
		"almost design doctor exist destroy candy zebra insane client grocery govern idea library degree two rebuild coffee hat scene deal average fresh measure potato",
		"welcome bean crystal pave chapter process bless tribe inside bottom exhaust hollow display envelope rally moral admit round hidden junk silly afraid awesome muffin",
		"end bicycle walnut empty bus silly camera lift fancy symptom office pluck detail unable cry sense scrap tuition relax amateur hold win debate hat",
		"cloud deal hurdle sound scout merit carpet identify fossil brass ancient keep disorder save lobster whisper course intact winter bullet flame mother upgrade install",
		"mutual duck begin remind release brave patrol squeeze abandon pact valid close fragile plastic disorder saddle bring inspire corn kitten reduce candy side honey",
	}
	bnbDeputyAddr, kavaDeputyAddr string
	bnbUserAddrs, kavaUserAddrs   []string
)

func TestMain(m *testing.M) {
	kavaConfig := sdk.GetConfig()
	kava.SetBech32AddressPrefixes(kavaConfig)
	kavaConfig.Seal()

	bnbDeputyAddr = bnbAddressFromMnemonic(bnbDeputyMnemonic)
	for _, m := range bnbUserMnemonics {
		bnbUserAddrs = append(bnbUserAddrs, bnbAddressFromMnemonic(m))
	}
	kavaDeputyAddr = kavaAddressFromMnemonic(kavaDeputyMnemonic)
	for _, m := range kavaUserMnemonics {
		kavaUserAddrs = append(kavaUserAddrs, kavaAddressFromMnemonic(m))
	}

	os.Exit(m.Run())
}

type logger interface {
	Log(args ...interface{})
}

func sendCompleteSwap(logger logger, senderExecutor, receiverExecutor common.Executor, senderAddr, receiverAddr string, swapAmount *big.Int, senderChainDeputyAddr string, heightSpan int64) error {

	// 1) Send initial swap

	rndNum, err := bep3.GenerateSecureRandomNumber()
	if err != nil {
		return fmt.Errorf("couldn't generate random number: %w", err)
	}
	timestamp := time.Now().Unix()
	rndHash := ec.BytesToHash(bep3.CalculateRandomHash(rndNum, timestamp))
	htltTxHash, cmnErr := senderExecutor.HTLT(
		rndHash,
		timestamp,
		heightSpan,
		senderChainDeputyAddr, // TODO
		receiverExecutor.GetDeputyAddress(),
		receiverAddr,
		swapAmount,
	)
	if cmnErr != nil {
		return fmt.Errorf("couldn't send htlt tx: %w", cmnErr)
	}

	err = wait(8*time.Second, func() (bool, error) {
		s := senderExecutor.GetSentTxStatus(htltTxHash)
		return s == store.TxSentStatusSuccess, nil
	})
	if err != nil {
		return fmt.Errorf("couldn't submit htlt: %w", err)
	}
	logger.Log("sender htlt created")

	// 4) Wait until deputy relays swap to receiver chain

	receiverSwapIDBz, err := receiverExecutor.CalcSwapId(rndHash, receiverExecutor.GetDeputyAddress(), senderAddr) // TODO senderAddr == senderExe.DeputyAddr
	if err != nil {
		return fmt.Errorf("couldn't calculate swap id: %w", err)
	}
	receiverSwapID := ec.BytesToHash(receiverSwapIDBz)

	err = wait(20*time.Second, func() (bool, error) {
		return receiverExecutor.HasSwap(receiverSwapID)
	})
	if err != nil {
		return fmt.Errorf("deputy did not relay swap: %w", err)
	}

	logger.Log("swap created on receiver by deputy")

	// 5) Send claim on receiver

	claimTxHash, cmnErr := receiverExecutor.Claim(receiverSwapID, ec.BytesToHash(rndNum))
	if cmnErr != nil {
		return fmt.Errorf("claim couldn't be submitted: %w", cmnErr)
	}

	err = wait(8*time.Second, func() (bool, error) {
		return receiverExecutor.GetSentTxStatus(claimTxHash) == store.TxSentStatusSuccess, nil
	})
	if err != nil {
		return fmt.Errorf("claim was not submitted: %w", err)
	}

	logger.Log("receiver htlt claimed")

	// 6) Wait until deputy relays claim to sender chian

	senderSwapIDBz, err := senderExecutor.CalcSwapId(rndHash, senderAddr, receiverExecutor.GetDeputyAddress()) // TODO senderAddr == senderExe.DeputyAddr
	if err != nil {
		return fmt.Errorf("couldn't calculate swap id: %w", err)
	}
	senderSwapID := ec.BytesToHash(senderSwapIDBz)

	wait(10*time.Second, func() (bool, error) {
		// check the deputy has relayed the claim by checking the status of the swap
		// once claimed it is no longer claimable, if it timesout it will become refundable
		claimable, err := senderExecutor.Claimable(senderSwapID)
		if err != nil {
			return false, err
		}
		refundable, err := senderExecutor.Refundable(senderSwapID)
		if err != nil {
			return false, err
		}
		return !(claimable || refundable), nil
	})
	logger.Log("sender htlt claimed by deputy")

	return nil
}
func TestBnbToKavaSwap(t *testing.T) {

	// 1) setup executors

	config := util.ParseConfigFromFile("deputy/config.json")
	bnbConfig := config.BnbConfig
	bnbConfig.RpcAddr = "tcp://localhost:26658"
	bnbConfig.Mnemonic = bnbUserMnemonics[0]

	senderExecutor := bnbExe.NewExecutor(types.ProdNetwork, bnbConfig)

	kavaConfig := config.KavaConfig
	kavaConfig.RpcAddr = "tcp://localhost:26657"
	kavaConfig.Mnemonic = kavaDeputyMnemonic

	receiverExecutor := kavaExe.NewExecutor(client.LocalNetwork, kavaConfig)

	senderAddr := bnbUserAddrs[0]
	receiverAddr := kavaUserAddrs[0]

	// 2) Cache account balances

	senderBalance, err := senderExecutor.GetBalance(senderAddr)
	require.NoError(t, err)
	receiverBalance, err := receiverExecutor.GetBalance(receiverAddr)
	require.NoError(t, err)

	// 3) Send swap

	swapAmount := big.NewInt(100_000_000)
	err = sendCompleteSwap(t, senderExecutor, receiverExecutor, senderAddr, receiverAddr, swapAmount, bnbDeputyAddr, 20000)
	require.NoError(t, err)

	// 4) Check balances

	senderBalanceFinal, err := senderExecutor.GetBalance(senderAddr)
	require.NoError(t, err)

	expectedSenderBalance := new(big.Int)
	expectedSenderBalance.Sub(senderBalance, swapAmount).Sub(expectedSenderBalance, big.NewInt(bnbHTLTFee))
	require.Zerof(t,
		expectedSenderBalance.Cmp(senderBalanceFinal),
		"expected: %s, actual: %s", expectedSenderBalance, senderBalanceFinal,
	)

	receiverBalanceFinal, err := receiverExecutor.GetBalance(receiverAddr)
	require.NoError(t, err)

	swapAmountReceiver := new(big.Int)
	swapAmountReceiver.Sub(swapAmount, config.ChainConfig.BnbFixedFee)
	expectedReceiverBalance := new(big.Int).Add(receiverBalance, swapAmountReceiver)
	require.Zero(t, config.ChainConfig.BnbRatio.Cmp(big.NewFloat(1)), "test does not support ratio conversions other than 1")
	require.Zerof(t,
		expectedReceiverBalance.Cmp(receiverBalanceFinal),
		"expected: %s, actual: %s", expectedReceiverBalance, receiverBalanceFinal,
	)

}

func TestKavaToBnbSwap(t *testing.T) {

	// 1) setup executors

	config := util.ParseConfigFromFile("deputy/config.json")

	kavaConfig := config.KavaConfig
	kavaConfig.RpcAddr = "tcp://localhost:26657"
	kavaConfig.Mnemonic = kavaUserMnemonics[0]
	senderExecutor := kavaExe.NewExecutor(client.LocalNetwork, kavaConfig)

	bnbConfig := config.BnbConfig
	bnbConfig.RpcAddr = "tcp://localhost:26658"
	bnbConfig.Mnemonic = bnbDeputyMnemonic
	receiverExecutor := bnbExe.NewExecutor(types.ProdNetwork, bnbConfig)

	senderAddr := kavaUserAddrs[0]
	receiverAddr := bnbUserAddrs[0]

	// 2) Cache account balances

	receiverBalance, err := receiverExecutor.GetBalance(receiverAddr)
	require.NoError(t, err)
	senderBalance, err := senderExecutor.GetBalance(senderAddr)
	require.NoError(t, err)

	// 3) Send swap

	swapAmount := big.NewInt(99_000_000)
	err = sendCompleteSwap(t, senderExecutor, receiverExecutor, senderAddr, receiverAddr, swapAmount, kavaDeputyAddr, 250)
	require.NoError(t, err)

	// 4) Check balances

	senderBalanceFinal, err := senderExecutor.GetBalance(senderAddr)
	require.NoError(t, err)

	expectedSenderBalance := new(big.Int)
	expectedSenderBalance.Sub(senderBalance, swapAmount) // no bnb tx fee when sending from kava
	require.Zerof(t,
		expectedSenderBalance.Cmp(senderBalanceFinal),
		"expected: %s, actual: %s", expectedSenderBalance, senderBalanceFinal,
	)

	receiverBalanceFinal, err := receiverExecutor.GetBalance(receiverAddr)
	require.NoError(t, err)

	swapAmountReceiver := new(big.Int)
	swapAmountReceiver.Sub(swapAmount, config.ChainConfig.OtherChainFixedFee)
	expectedReceiverBalance := new(big.Int).Add(receiverBalance, swapAmountReceiver)
	require.Zero(t, config.ChainConfig.OtherChainRatio.Cmp(big.NewFloat(1)), "test does not support ratio conversions other than 1")
	require.Zerof(t,
		expectedReceiverBalance.Cmp(receiverBalanceFinal),
		"expected: %s, actual: %s", expectedReceiverBalance, receiverBalanceFinal,
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

func bnbAddressFromMnemonic(mnemonic string) string {
	manager, err := bnbKeys.NewMnemonicKeyManager(mnemonic)
	if err != nil {
		panic(err.Error())
	}
	return manager.GetAddr().String()
}

func kavaAddressFromMnemonic(mnemonic string) string {
	manager, err := kavaKeys.NewMnemonicKeyManager(mnemonic, kava.Bip44CoinType)
	if err != nil {
		panic(err.Error())
	}
	return manager.GetAddr().String()
}
