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

	ec "github.com/ethereum/go-ethereum/common"
	sdk "github.com/kava-labs/cosmos-sdk/types"
	"github.com/kava-labs/go-sdk/client"
	"github.com/kava-labs/go-sdk/kava"
	"github.com/kava-labs/go-sdk/kava/bep3"
	tmbytes "github.com/kava-labs/tendermint/libs/bytes"
	"github.com/tendermint/go-amino"

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
	cdc := kava.MakeCodec()

	// Set up Kava HTTP client and set codec
	kavaClient := client.NewKavaClient(cdc, cfg.Mnemonic, kava.Bip44CoinType, cfg.RpcAddr, networkType)
	kavaClient.Keybase.SetCodec(cdc)

	return &Executor{
		Config:        cfg,
		Client:        kavaClient,
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

	blockHash := hex.EncodeToString(block.BlockID.Hash)
	for idx, t := range block.Block.Data.Txs {
		txResult := blockResults.TxsResults[idx]
		if txResult.Code != 0 {
			continue
		}

		txHash := hex.EncodeToString(t.Hash())

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
				if !realMsg.CrossChain {
					continue
				}

				if len(realMsg.Amount) != 1 {
					continue
				}

				signer := msg.GetSigners()[0]
				randomNumberHash := hex.EncodeToString(realMsg.RandomNumberHash)

				// Parse swap ID from create atomic swap event
				var swapID string
				for _, event := range txResult.Events {
					if event.GetType() == "create_atomic_swap" {
						for _, attribute := range event.GetAttributes() {
							if string(attribute.GetKey()) == "atomic_swap_id" {
								swapID = string(attribute.GetValue())
							}
						}
					}
				}
				if len(strings.TrimSpace(swapID)) == 0 {
					util.Logger.Errorf("err='atomic_swap_id' event attribute not found")
					continue
				}

				txLog := store.TxLog{
					Chain:  common.ChainKava,
					TxType: store.TxTypeOtherHTLT,
					TxHash: txHash,

					SwapId:           swapID,
					SenderAddr:       signer.String(),
					ReceiverAddr:     realMsg.To.String(),
					SenderOtherChain: realMsg.SenderOtherChain,
					OtherChainAddr:   realMsg.RecipientOtherChain,
					InAmount:         realMsg.Amount[0].String(),
					OutAmount:        strconv.FormatInt(realMsg.Amount[0].Amount.Int64(), 10),
					OutCoin:          realMsg.Amount[0].Denom,
					RandomNumberHash: randomNumberHash,
					ExpireHeight:     int64(realMsg.HeightSpan) + height,
					Timestamp:        realMsg.Timestamp,
					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)
			case bep3.MsgClaimAtomicSwap:
				signer := msg.GetSigners()[0]
				swapID := hex.EncodeToString(realMsg.SwapID)
				randomNum := hex.EncodeToString(realMsg.RandomNumber)

				txLog := store.TxLog{
					Chain:  common.ChainKava,
					TxType: store.TxTypeOtherClaim,
					TxHash: txHash,

					SenderAddr:   signer.String(),
					SwapId:       swapID,
					RandomNumber: randomNum,

					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)
			case bep3.MsgRefundAtomicSwap:
				signer := msg.GetSigners()[0]
				swapID := hex.EncodeToString(realMsg.SwapID)

				txLog := store.TxLog{
					Chain:  common.ChainKava,
					TxType: store.TxTypeOtherRefund,
					TxHash: txHash,

					SenderAddr: signer.String(),
					SwapId:     swapID,

					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)
			default:
			}
		}
	}

	blockAndTxLogs := &common.BlockAndTxLogs{
		Height:          block.Block.Height,
		BlockHash:       block.BlockID.Hash.String(),
		ParentBlockHash: block.Block.Header.LastBlockID.Hash.String(),
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

	recipient, err := sdk.AccAddressFromBech32(recipientAddr)
	if err != nil {
		return "", common.NewError(err, false)
	}

	if !outAmount.IsInt64() {
		return "", common.NewError(
			fmt.Errorf(fmt.Sprintf("out amount(%s) is not int64", outAmount.String())), false)
	}

	outCoin := sdk.NewCoins(sdk.NewInt64Coin(executor.Config.Symbol, outAmount.Int64()))

	if executor.Client.Keybase == nil {
		return "", common.NewError(errors.New("Err: key missing"), false)
	}

	fromAddr := executor.Client.Keybase.GetAddr()

	createMsg := bep3.NewMsgCreateAtomicSwap(
		fromAddr,
		recipient,
		otherChainRecipientAddr,
		otherChainSenderAddr,
		tmbytes.HexBytes(randomNumberHash.Bytes()),
		timestamp,
		outCoin,
		uint64(heightSpan),
	)

	res, err := executor.Client.Broadcast(createMsg, client.Sync)
	if err != nil {
		return "", common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return "", common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}

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
		tmbytes.HexBytes(swapId.Bytes()),
		tmbytes.HexBytes(trimmedRandomNumber),
	)

	res, err := executor.Client.Broadcast(claimMsg, client.Sync)
	if err != nil {
		return "", common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return "", common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}

	return res.Hash.String(), nil
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
		tmbytes.HexBytes(swapId.Bytes()),
	)

	res, err := executor.Client.Broadcast(refundMsg, client.Sync)
	if err != nil {
		return "", common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return "", common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}

	return res.Hash.String(), nil
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
func (executor *Executor) QuerySwap(swapId []byte) (swap bep3.AtomicSwap, isExist bool, err error) {
	swap, err = executor.Client.GetSwapByID(tmbytes.HexBytes(swapId))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return bep3.AtomicSwap{}, false, nil
		}
		return bep3.AtomicSwap{}, false, err
	}

	return swap, true, nil
}

// HasSwap returns true if an AtomicSwap with this ID exists on kava
func (executor *Executor) HasSwap(swapId ec.Hash) (bool, error) {
	_, isExist, err := executor.QuerySwap(swapId.Bytes())
	return isExist, err
}

// GetSwap gets an AtomicSwap by its ID
func (executor *Executor) GetSwap(swapId ec.Hash) (*common.SwapRequest, error) {
	swap, isExist, err := executor.QuerySwap(swapId.Bytes())
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf("swap does not exist, swapId=%s", swapId.String())
	}
	if len(swap.Amount) != 1 {
		return nil, fmt.Errorf("swap request has multi coins, coin_types=%d", swap.Amount.Len())
	}

	return &common.SwapRequest{
		Id:                  swapId,
		RandomNumberHash:    ec.BytesToHash(swap.RandomNumberHash),
		ExpireHeight:        swap.ExpireHeight,
		SenderAddress:       swap.Sender.String(),
		RecipientAddress:    swap.Recipient.String(),
		OutAmount:           big.NewInt(swap.Amount[0].Amount.Int64()),
		RecipientOtherChain: swap.RecipientOtherChain,
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

	kavaBalance := balances.AmountOf("ukava").Int64()
	tokenBalance := balances.AmountOf(executor.Config.Symbol).Int64()

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

func isInvalidSequenceError(err string) bool {
	return strings.Contains(strings.ToLower(err), "invalid sequence")
}
