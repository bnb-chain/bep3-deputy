package kava

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	// Cosmos-SDK, Tendermint, Ethereum
	sdk "github.com/cosmos/cosmos-sdk/types"
	ec "github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tendermint/libs/common"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	// Kava
	"github.com/kava-labs/go-sdk/client"
	"github.com/kava-labs/kava/app"
	bep3 "github.com/kava-labs/kava/x/bep3/types"

	// Bep3-deputy
	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

var _ common.Executor = &Executor{}

// Executor implements the common executor interface
type Executor struct {
	mutex  sync.Mutex
	Config *util.KavaConfig

	NetworkType client.ChainNetwork

	Client        *client.KavaClient
	Cdc           *amino.Codec
	DeputyAddress sdk.AccAddress
}

// NewExecutor creates a new Executor
func NewExecutor(rpcAddr string, networkType client.ChainNetwork, cfg *util.KavaConfig) *Executor {
	cdc := app.MakeCodec()

	// Set up Kava HTTP client and set codec
	kava := client.NewKavaClient(cdc, cfg.Mnemonic, cfg.RpcAddr, networkType)
	kava.Keybase.SetCodec(cdc)

	return &Executor{
		Config:        cfg,
		Client:        kava,
		Cdc:           cdc,
		DeputyAddress: cfg.DeputyAddr,
	}
}

// GetChain gets the chain ID
func (executor *Executor) GetChain() string {
	return common.ChainKava
}

// GetBlockAndTxs parses all transactions from a specific block height
func (executor *Executor) GetBlockAndTxs(height int64) (*common.BlockAndTxLogs, error) {
	block, err := executor.Client.HTTP.Block(&height)
	if err != nil {
		return nil, err
	}

	blockResults, err := executor.Client.HTTP.BlockResults(&height)
	if err != nil {
		return nil, err
	}

	txLogs := make([]*store.TxLog, 0)

	blockHash := hex.EncodeToString(block.BlockMeta.BlockID.Hash)
	for idx, t := range block.Block.Data.Txs {
		txResult := blockResults.Results.DeliverTx[idx]
		if txResult.Code != 0 {
			continue
		}

		txHash := hex.EncodeToString(t.Hash())
		// TODO: remove print
		fmt.Println("Witnessed new tx. Tx hash:", txHash)

		var parsedTx sdk.Tx
		err := executor.Cdc.UnmarshalBinaryLengthPrefixed(t, &parsedTx)
		if err != nil {
			return nil, err
		}

		if err != nil {
			util.Logger.Errorf("parse tx error, err=%s", err.Error())
			continue
		}

		msgs := parsedTx.GetMsgs()
		for _, msg := range msgs {
			switch realMsg := msg.(type) {
			case bep3.MsgCreateAtomicSwap:
				// TODO: remove print
				fmt.Println("MsgCreateAtomicSwap case")

				if !realMsg.CrossChain {
					continue
				}

				if len(realMsg.Amount) != 1 {
					continue
				}

				signer := msg.GetSigners()[0]
				randomNumberHash := hex.EncodeToString(realMsg.RandomNumberHash)

				txLog := store.TxLog{
					Chain:  common.ChainKava,
					TxType: store.TxTypeBEP2HTLT,
					TxHash: txHash,

					SwapId:           strings.ReplaceAll(txResult.Log, "Msg 0: swapID: ", ""),
					SenderAddr:       signer.String(),
					ReceiverAddr:     realMsg.To.String(),
					SenderOtherChain: realMsg.SenderOtherChain,
					OtherChainAddr:   realMsg.RecipientOtherChain,
					InAmount:         realMsg.ExpectedIncome,
					OutAmount:        strconv.FormatInt(realMsg.Amount[0].Amount.Int64(), 10),
					OutCoin:          realMsg.Amount[0].Denom,
					RandomNumberHash: randomNumberHash,
					ExpireHeight:     realMsg.HeightSpan + height,
					Timestamp:        realMsg.Timestamp,

					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)

				// TODO: Remove. This is for testing.
				swapIDHashed := ec.HexToHash("e878263268853c56987835218621bda1147f4a04ecec0d092cc84deedce3dc44")
				randomNumberHashed := ec.BytesToHash([]byte("15"))
				executor.Claim(swapIDHashed, randomNumberHashed)

			case bep3.MsgClaimAtomicSwap:
				// TODO: remove print
				fmt.Println("Saw new MsgClaimAtomicSwap")

				signer := msg.GetSigners()[0]
				swapID := hex.EncodeToString(realMsg.SwapID)
				randomNum := hex.EncodeToString(realMsg.RandomNumber)

				txLog := store.TxLog{
					Chain:  common.ChainKava,
					TxType: store.TxTypeBEP2Claim,
					TxHash: txHash,

					SenderAddr:   signer.String(),
					SwapId:       swapID,
					RandomNumber: randomNum,

					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)
			case bep3.MsgRefundAtomicSwap:
				// TODO: remove print
				fmt.Println("Saw new MsgRefundAtomicSwap")

				signer := msg.GetSigners()[0]
				swapID := hex.EncodeToString(realMsg.SwapID)

				txLog := store.TxLog{
					Chain:  common.ChainKava,
					TxType: store.TxTypeBEP2Refund,
					TxHash: txHash,

					SenderAddr: signer.String(),
					SwapId:     swapID,

					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)

				// TODO: Remove. This is for testing.
				//  kvcli q bep3 calc-rnh 15 100 -> "2219edb58f397dee8cdefb5bc05749353e40c47dfcf0654a2b0318912a0dc270"
				executor.HTLT(ec.HexToHash("2219edb58f397dee8cdefb5bc05749353e40c47dfcf0654a2b0318912a0dc270"), 100, 80,
					"kava15qdefkmwswysgg4qxgqpqr35k3m49pkx2jdfnw", "0x9eB05a790e2De0a047a57a22199D8CccEA6d6D5A",
					"0x9eB05a790e2De0a047a57a22199D8CccEA6d6D5A", big.NewInt(1000))
			default:
			}
		}
	}

	blockAndTxLogs := &common.BlockAndTxLogs{
		Height:          block.Block.Height,
		BlockHash:       block.BlockMeta.BlockID.Hash.String(),
		ParentBlockHash: block.BlockMeta.Header.LastBlockID.Hash.String(),
		BlockTime:       block.Block.Time.Unix(),
		TxLogs:          txLogs,
	}

	return blockAndTxLogs, nil
}

// HTLT sends a transaction containing a MsgCreateAtomicSwap to kava
func (executor *Executor) HTLT(randomNumberHash ec.Hash, timestamp int64, heightSpan int64, recipientAddr string,
	otherChainSenderAddr string, otherChainRecipientAddr string, outAmount *big.Int) (string, *common.Error) {

	executor.mutex.Lock()
	defer executor.mutex.Unlock()

	// TODO: why is Binance setting here instead of globally? this applies to claim/refund as well
	// executor.RpcClient.SetKeyManager(keyManager)
	// defer executor.RpcClient.SetKeyManager(nil)

	recipient, err := sdk.AccAddressFromBech32(recipientAddr)
	if err != nil {
		return "", common.NewError(err, false)
	}

	if !outAmount.IsInt64() {
		return "", common.NewError(
			fmt.Errorf(fmt.Sprintf("out amount(%s) is not int64", outAmount.String())), false)
	}

	coinSymbol := "btc" // TODO: use executor.Config.Symbol
	outCoin := sdk.NewCoins(sdk.NewInt64Coin(coinSymbol, outAmount.Int64()))

	if executor.Client.Keybase == nil {
		return "", common.NewError(errors.New("Err: key missing"), false)
	}

	fromAddr := executor.Client.Keybase.GetAddr()

	createMsg := bep3.NewMsgCreateAtomicSwap(
		fromAddr,
		recipient,
		otherChainRecipientAddr,
		otherChainSenderAddr,
		cmn.HexBytes(randomNumberHash.Bytes()),
		timestamp,
		outCoin,
		fmt.Sprintf("%d%s", outAmount.Int64(), coinSymbol),
		heightSpan,
		true,
	)

	// TODO: are 'options...' required?
	res, err := executor.Broadcast(createMsg, client.Sync)
	if err != nil {
		return "", common.NewError(err, false) // TODO: 'true'?
	}

	// TODO: remove print
	fmt.Println("Result of msg broadcast:", res)
	fmt.Println("Tx hash:", res.Hash.String())

	return res.Hash.String(), nil
}

// GetFetchInterval gets the duration between fetches
func (executor *Executor) GetFetchInterval() time.Duration {
	return time.Duration(executor.Config.FetchInterval) * time.Second
}

// Claim sends a MsgClaimAtomicSwap to kava
func (executor *Executor) Claim(swapId ec.Hash, randomNumber ec.Hash) (string, *common.Error) {
	executor.mutex.Lock()
	defer executor.mutex.Unlock()

	if executor.Client.Keybase == nil {
		return "", common.NewError(errors.New("Err: key missing"), false)
	}

	trimmedRandomNumber := bytes.Trim(randomNumber.Bytes(), "\x00")

	claimMsg := bep3.NewMsgClaimAtomicSwap(
		executor.DeputyAddress,
		cmn.HexBytes(swapId.Bytes()),
		cmn.HexBytes(trimmedRandomNumber),
	)

	res, err := executor.Broadcast(claimMsg, client.Sync)
	if err != nil {
		return "", common.NewError(err, false)
	}

	// TODO: remove print
	fmt.Println("Result of msg broadcast:", res)
	fmt.Println("Tx hash:", res.Hash.String())

	return "", nil
}

// Refund sends a MsgRefundAtomicSwap to kava
func (executor *Executor) Refund(swapId ec.Hash) (string, *common.Error) {
	executor.mutex.Lock()
	defer executor.mutex.Unlock()

	if executor.Client.Keybase == nil {
		return "", common.NewError(errors.New("Err: key missing"), false)
	}

	refundMsg := bep3.NewMsgRefundAtomicSwap(
		executor.DeputyAddress,
		cmn.HexBytes(swapId.Bytes()),
	)

	res, err := executor.Broadcast(refundMsg, client.Sync)
	if err != nil {
		return "", common.NewError(err, false)
	}

	// TODO: remove print
	fmt.Println("Result of msg broadcast:", res)
	fmt.Println("Tx hash:", res.Hash.String())

	return "", nil
}

// GetSentTxStatus gets a sent transaction's status
func (executor *Executor) GetSentTxStatus(hash string) store.TxStatus {
	bz, err := hex.DecodeString(hash)
	if err != nil {
		return store.TxSentStatusNotFound
	}
	txResult, err := executor.Client.HTTP.Tx(bz, false)
	if err != nil {
		return store.TxSentStatusNotFound
	}
	if txResult.TxResult.Code == 0 {
		return store.TxSentStatusSuccess
	}
	return store.TxSentStatusFailed
}

// QuerySwap queries kava for an AtomicSwap
func (executor *Executor) QuerySwap(swapID []byte) (swap bep3.AtomicSwap, isExist bool, err error) {
	swap, err = executor.Client.GetSwapByID(swapID)
	if err != nil {
		//if strings.Contains(err.Error(), "No matched swap with swapID") {
		if strings.Contains(err.Error(), "zero records") {
			return bep3.AtomicSwap{}, false, nil
		}
		return bep3.AtomicSwap{}, false, err
	}

	return swap, true, nil
}

// HasSwap returns true if an AtomicSwap with this ID exists on kava
func (executor *Executor) HasSwap(swapID ec.Hash) (bool, error) {
	_, isExist, err := executor.QuerySwap(swapID[:])
	return isExist, err
}

// GetSwap gets an AtomicSwap by its ID
func (executor *Executor) GetSwap(swapID ec.Hash) (*common.SwapRequest, error) {
	swap, isExist, err := executor.QuerySwap(swapID[:])
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf("swap does not exist, swapId=%s", swapID.String())
	}
	if len(swap.Amount) != 1 {
		return nil, fmt.Errorf("swap request has multi coins, coin_types=%d", swap.Amount.Len())
	}

	return &common.SwapRequest{
		Id:                  swapID,
		RandomNumberHash:    ec.BytesToHash(swap.RandomNumberHash),
		ExpireHeight:        swap.ExpireHeight,
		SenderAddress:       swap.Sender.String(),
		RecipientAddress:    swap.Recipient.String(),
		OutAmount:           big.NewInt(swap.Amount[0].Amount.Int64()),
		RecipientOtherChain: swap.SenderOtherChain,
		//TODO: RecipientOtherChain: swap.RecipientOtherChain,
	}, nil
}

// GetHeight gets the current block height of the kava blockchain
func (executor *Executor) GetHeight() (int64, error) {
	status, err := executor.Client.HTTP.Status()
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// Claimable returns true is an AtomicSwap is currently claimable
func (executor *Executor) Claimable(swapId ec.Hash) (bool, error) {
	swap, isExist, err := executor.QuerySwap(swapId[:])
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, nil
	}

	status, err := executor.Client.HTTP.Status()
	if err != nil {
		return false, err
	}

	if swap.Status == bep3.Open && status.SyncInfo.LatestBlockHeight < swap.ExpireHeight {
		return true, nil
	}
	return false, nil
}

// Refundable returns true is an AtomicSwap is currently refundable
func (executor *Executor) Refundable(swapId ec.Hash) (bool, error) {
	swap, isExist, err := executor.QuerySwap(swapId[:])
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, nil
	}

	status, err := executor.Client.HTTP.Status()
	if err != nil {
		return false, err
	}

	if swap.Status == bep3.Open && status.SyncInfo.LatestBlockHeight >= swap.ExpireHeight {
		return true, nil
	}
	return false, nil
}

// GetBalance gets the deputy's current kava balance
func (executor *Executor) GetBalance() (*big.Int, error) {
	deputy, err := executor.Client.GetAccount(executor.DeputyAddress)
	if err != nil {
		return big.NewInt(0), err
	}

	if deputy.Address.Empty() {
		return big.NewInt(0), errors.New("invalid deputy address")
	}

	for _, coin := range deputy.Coins {
		// TODO: Confirm that executor.Config.Symbol is lowercase
		if coin.Denom == executor.Config.Symbol {
			return big.NewInt(coin.Amount.Int64()), nil
		}
	}

	return big.NewInt(0), fmt.Errorf(fmt.Sprintf("deputy doesn't have any %s", executor.Config.Symbol))
}

// GetDeputyAddress gets the deputy's address from the config
func (executor *Executor) GetDeputyAddress() string {
	return executor.Config.DeputyAddr.String()
}

// CalcSwapId calculates the swap ID for a given random number hash, sender, and sender other chain
func (executor *Executor) CalcSwapId(randomNumberHash ec.Hash, sender string, senderOtherChain string) ([]byte, error) {
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return nil, err
	}
	return bep3.CalculateSwapID(randomNumberHash[:], senderAddr, senderOtherChain), nil
}

// IsSameAddress checks for equality between two addresses
func (executor *Executor) IsSameAddress(addrA string, addrB string) bool {
	return strings.ToLower(addrA) == strings.ToLower(addrB)
}

// GetStatus gets the total coin balances of the deputy
func (executor *Executor) GetStatus() (interface{}, error) {
	kavaStatus := &common.KavaStatus{}

	deputy, err := executor.Client.GetAccount(executor.DeputyAddress)
	if err != nil {
		return nil, err
	}

	if deputy.Address.Empty() {
		return big.NewInt(0), errors.New("invalid deputy address")
	}

	kavaStatus.Balance = sdk.NewCoins(deputy.Coins...)
	return kavaStatus, nil
}

// GetBalanceAlertMsg constructs an alert message if the deputy's balance is low
func (executor *Executor) GetBalanceAlertMsg() (string, error) {
	if executor.Config.KavaBalanceAlertThreshold == 0 && executor.Config.TokenBalanceAlertThreshold == 0 {
		return "", nil
	}

	deputy, err := executor.Client.GetAccount(executor.DeputyAddress)
	if err != nil {
		return "", err
	}

	if deputy.Address.Empty() {
		return "", errors.New("invalid deputy address")
	}

	balances := sdk.NewCoins(deputy.Coins...)

	// TODO: token symbols are both "KAVA", which doesn't work
	// kavaBalance := balances.AmountOf(common.KAVASymbol).Int64()
	// tokenBalance := balances.AmountOf(executor.Config.Symbol).Int64()
	kavaBalance := balances.AmountOf("ukava").Int64()
	tokenBalance := balances.AmountOf("btc").Int64()

	alertMsg := ""
	if executor.Config.KavaBalanceAlertThreshold > 0 && kavaBalance < executor.Config.KavaBalanceAlertThreshold {
		alertMsg = alertMsg + fmt.Sprintf("Kava balance(%d) is less than %d\n", kavaBalance, executor.Config.KavaBalanceAlertThreshold)
	}
	if executor.Config.TokenBalanceAlertThreshold > 0 && tokenBalance < executor.Config.TokenBalanceAlertThreshold {
		alertMsg = alertMsg + fmt.Sprintf("%s balance(%d) is less than %d", executor.Config.Symbol,
			tokenBalance, executor.Config.TokenBalanceAlertThreshold)
	}

	return alertMsg, nil
}

// Broadcast sends a transaction to Kava containing the given msg
func (executor *Executor) Broadcast(m sdk.Msg, syncType client.SyncType) (*ctypes.ResultBroadcastTx, *common.Error) {
	res, err := executor.Client.Broadcast(m, syncType)
	if err != nil {
		return &ctypes.ResultBroadcastTx{}, common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return &ctypes.ResultBroadcastTx{}, common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}
	return res, nil
}

func isInvalidSequenceError(err string) bool {
	return strings.Contains(err, "Invalid sequence")
}
