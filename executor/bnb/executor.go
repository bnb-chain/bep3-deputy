package bnb

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

	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	sdkMsg "github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	ec "github.com/ethereum/go-ethereum/common"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

var _ common.Executor = &Executor{}

type Executor struct {
	mutex  sync.Mutex
	Config *util.BnbConfig

	NetworkType types.ChainNetwork
	RpcClient   rpc.Client
}

func NewExecutor(networkType types.ChainNetwork, cfg *util.BnbConfig) *Executor {
	keyManager, err := getKeyManager(cfg)
	if err != nil {
		panic(fmt.Sprintf("new key manager err, err=%s", err.Error()))
	}

	rpcClient := rpc.NewRPCClient(cfg.RpcAddr, networkType)
	rpcClient.SetLogger(util.SdkLogger)

	if !bytes.Equal(cfg.DeputyAddr.Bytes(), keyManager.GetAddr().Bytes()) {
		panic(fmt.Sprintf(
			"deputy address supplied in config (%s) does not match mnemonic (%s)",
			cfg.DeputyAddr, keyManager.GetAddr(),
		))
	}
	return &Executor{
		Config:      cfg,
		NetworkType: networkType,
		RpcClient:   rpcClient,
	}
}
func (executor *Executor) GetChain() string {
	return common.ChainBinance
}

func getKeyManager(config *util.BnbConfig) (keys.KeyManager, error) {
	var bnbMnemonic string
	if config.KeyType == util.KeyTypeAWSMnemonic {
		awsMnemonic, err := util.GetSecret(config.AWSSecretName, config.AWSRegion)
		if err != nil {
			return nil, err
		}
		bnbMnemonic = awsMnemonic
	} else {
		bnbMnemonic = config.Mnemonic
	}

	return keys.NewMnemonicKeyManager(bnbMnemonic)
}

func (executor *Executor) GetBlockAndTxs(height int64) (*common.BlockAndTxLogs, error) {
	block, err := executor.RpcClient.Block(&height)
	if err != nil {
		return nil, err
	}

	blockResults, err := executor.RpcClient.BlockResults(&height)
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

		stdTx, err := rpc.ParseTx(tx.Cdc, t)
		if err != nil {
			util.Logger.Errorf("parse tx error, err=%s", err.Error())
			continue
		}

		msgs := stdTx.GetMsgs()
		for _, msg := range msgs {
			switch realMsg := msg.(type) {
			case sdkMsg.HTLTMsg:
				if !realMsg.CrossChain {
					continue
				}

				if len(realMsg.Amount) != 1 {
					continue
				}

				signer := msg.GetSigners()[0]
				randomNumberHash := hex.EncodeToString(realMsg.RandomNumberHash)

				txLog := store.TxLog{
					Chain:  common.ChainBinance,
					TxType: store.TxTypeBEP2HTLT,
					TxHash: txHash,

					SwapId:           strings.ReplaceAll(txResult.Log, "Msg 0: swapID: ", ""),
					SenderAddr:       signer.String(),
					ReceiverAddr:     realMsg.To.String(),
					SenderOtherChain: realMsg.SenderOtherChain,
					OtherChainAddr:   realMsg.RecipientOtherChain,
					InAmount:         realMsg.ExpectedIncome,
					OutAmount:        strconv.FormatInt(realMsg.Amount[0].Amount, 10),
					OutCoin:          realMsg.Amount[0].Denom,
					RandomNumberHash: randomNumberHash,
					ExpireHeight:     realMsg.HeightSpan + height,
					Timestamp:        realMsg.Timestamp,

					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)
			case sdkMsg.ClaimHTLTMsg:
				signer := msg.GetSigners()[0]
				swapId := hex.EncodeToString(realMsg.SwapID)
				randomNum := hex.EncodeToString(realMsg.RandomNumber)

				txLog := store.TxLog{
					Chain:  common.ChainBinance,
					TxType: store.TxTypeBEP2Claim,
					TxHash: txHash,

					SenderAddr:   signer.String(),
					SwapId:       swapId,
					RandomNumber: randomNum,

					Height:    height,
					BlockHash: blockHash,
				}
				txLogs = append(txLogs, &txLog)
			case sdkMsg.RefundHTLTMsg:
				signer := msg.GetSigners()[0]
				swapId := hex.EncodeToString(realMsg.SwapID)

				txLog := store.TxLog{
					Chain:  common.ChainBinance,
					TxType: store.TxTypeBEP2Refund,
					TxHash: txHash,

					SenderAddr: signer.String(),
					SwapId:     swapId,

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
		BlockHash:       block.BlockMeta.BlockID.Hash.String(),
		ParentBlockHash: block.BlockMeta.Header.LastBlockID.Hash.String(),
		BlockTime:       block.Block.Time.Unix(),
		TxLogs:          txLogs,
	}

	return blockAndTxLogs, nil
}

func (executor *Executor) HTLT(randomNumberHash ec.Hash, timestamp int64, heightSpan int64, recipientAddr string, otherChainSenderAddr string, otherChainRecipientAddr string, outAmount *big.Int) (string, *common.Error) {
	executor.mutex.Lock()
	defer executor.mutex.Unlock()

	keyManager, err := getKeyManager(executor.Config)
	if err != nil {
		return "", common.NewError(err, false)
	}
	executor.RpcClient.SetKeyManager(keyManager)
	defer executor.RpcClient.SetKeyManager(nil)

	bep2Addr, err := types.AccAddressFromBech32(recipientAddr)
	if err != nil {
		return "", common.NewError(err, false)
	}

	if !outAmount.IsInt64() {
		return "", common.NewError(
			fmt.Errorf("out amount(%s) is not int64", outAmount),
			false,
		)
	}

	outCoin := types.Coin{
		Denom:  executor.Config.Symbol,
		Amount: outAmount.Int64(),
	}

	res, err := executor.RpcClient.HTLT(bep2Addr, otherChainRecipientAddr, otherChainSenderAddr, randomNumberHash[:],
		timestamp, types.Coins{outCoin}, "", heightSpan, true, rpc.Sync)
	if err != nil {
		return "", common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return "", common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}
	return res.Hash.String(), nil
}

func isInvalidSequenceError(err string) bool {
	return strings.Contains(err, "Invalid sequence")
}

func (executor *Executor) GetFetchInterval() time.Duration {
	return time.Duration(executor.Config.FetchInterval) * time.Second
}

func (executor *Executor) Claim(swapId ec.Hash, randomNumber ec.Hash) (string, *common.Error) {
	executor.mutex.Lock()
	defer executor.mutex.Unlock()

	keyManager, err := getKeyManager(executor.Config)
	if err != nil {
		return "", common.NewError(err, false)
	}
	executor.RpcClient.SetKeyManager(keyManager)
	defer executor.RpcClient.SetKeyManager(nil)

	res, err := executor.RpcClient.ClaimHTLT(swapId[:], randomNumber[:], rpc.Sync)
	if err != nil {
		return "", common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return "", common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}
	return res.Hash.String(), nil
}

func (executor *Executor) Refund(swapId ec.Hash) (string, *common.Error) {
	executor.mutex.Lock()
	defer executor.mutex.Unlock()

	keyManager, err := getKeyManager(executor.Config)
	if err != nil {
		return "", common.NewError(err, false)
	}
	executor.RpcClient.SetKeyManager(keyManager)
	defer executor.RpcClient.SetKeyManager(nil)

	res, err := executor.RpcClient.RefundHTLT(swapId[:], rpc.Sync)
	if err != nil {
		return "", common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return "", common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}
	return res.Hash.String(), nil
}

func (executor *Executor) GetSentTxStatus(hash string) store.TxStatus {
	bz, err := hex.DecodeString(hash)
	if err != nil {
		return store.TxSentStatusNotFound
	}
	txResult, err := executor.RpcClient.Tx(bz, false)
	if err != nil {
		return store.TxSentStatusNotFound
	}
	if txResult.TxResult.Code == 0 {
		return store.TxSentStatusSuccess
	} else {
		return store.TxSentStatusFailed
	}
}

func (executor *Executor) QuerySwap(swapId []byte) (swap types.AtomicSwap, isExist bool, err error) {
	swap, err = executor.RpcClient.GetSwapByID(swapId)
	if err != nil {
		//if strings.Contains(err.Error(), "No matched swap with swapID") {
		if strings.Contains(err.Error(), "zero records") {
			return types.AtomicSwap{}, false, nil
		} else {
			return types.AtomicSwap{}, false, err
		}
	}

	return swap, true, nil
}

func (executor *Executor) HasSwap(swapId ec.Hash) (bool, error) {
	_, isExist, err := executor.QuerySwap(swapId[:])
	return isExist, err
}

func (executor *Executor) GetSwap(swapId ec.Hash) (*common.SwapRequest, error) {
	swap, isExist, err := executor.QuerySwap(swapId[:])
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf("swap does not exist, swapId=%s", swapId.String())
	}
	if len(swap.OutAmount) != 1 {
		return nil, fmt.Errorf("swap request has multi coins, coin_types=%d", swap.OutAmount.Len())
	}

	return &common.SwapRequest{
		Id:                  swapId,
		RandomNumberHash:    ec.BytesToHash(swap.RandomNumberHash),
		ExpireHeight:        swap.ExpireHeight,
		SenderAddress:       swap.From.String(),
		RecipientAddress:    swap.To.String(),
		OutAmount:           big.NewInt(swap.OutAmount[0].Amount),
		RecipientOtherChain: swap.RecipientOtherChain,
	}, nil
}

func (executor *Executor) GetHeight() (int64, error) {
	status, err := executor.RpcClient.Status()
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

func (executor *Executor) Claimable(swapId ec.Hash) (bool, error) {
	swap, isExist, err := executor.QuerySwap(swapId[:])
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, nil
	}

	status, err := executor.RpcClient.Status()
	if err != nil {
		return false, err
	}

	if swap.Status == types.Open && status.SyncInfo.LatestBlockHeight < swap.ExpireHeight {
		return true, nil
	} else {
		return false, nil
	}
}

func (executor *Executor) Refundable(swapId ec.Hash) (bool, error) {
	swap, isExist, err := executor.QuerySwap(swapId[:])
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, nil
	}

	status, err := executor.RpcClient.Status()
	if err != nil {
		return false, err
	}

	if swap.Status == types.Open && status.SyncInfo.LatestBlockHeight >= swap.ExpireHeight {
		return true, nil
	} else {
		return false, nil
	}
}

func (executor *Executor) GetBalance(addressString string) (*big.Int, error) {
	address, err := types.AccAddressFromBech32(addressString)
	if err != nil {
		return big.NewInt(0), err
	}
	account, err := executor.RpcClient.GetAccount(address)
	if err != nil {
		return big.NewInt(0), err
	}

	if account == nil || account.GetCoins() == nil {
		return big.NewInt(0), errors.New("get nil account")
	}

	tokenBalance := account.GetCoins().AmountOf(executor.Config.Symbol)
	return big.NewInt(tokenBalance), nil
}

func (executor *Executor) Balance() ([]types.TokenBalance, error) {
	account, err := executor.RpcClient.GetAccount(executor.Config.DeputyAddr)
	if err != nil {
		return nil, err
	}
	coins := account.GetCoins()

	symbols := make([]string, 0, len(coins))
	balances := make([]types.TokenBalance, 0, len(coins))
	for _, coin := range coins {
		symbols = append(symbols, coin.Denom)
		// count locked and frozen coins
		var locked, frozen int64
		acc := account.(types.NamedAccount)
		if acc != nil {
			locked = acc.GetLockedCoins().AmountOf(coin.Denom)
			frozen = acc.GetFrozenCoins().AmountOf(coin.Denom)
		}
		balances = append(balances, types.TokenBalance{
			Symbol: coin.Denom,
			Free:   types.Fixed8(coins.AmountOf(coin.Denom)),
			Locked: types.Fixed8(locked),
			Frozen: types.Fixed8(frozen),
		})
	}
	return balances, nil
}

func (executor *Executor) GetDeputyAddress() string {
	return executor.Config.DeputyAddr.String()
}

func (executor *Executor) GetColdWalletAddress() string {
	return executor.Config.ColdWalletAddr.String()
}

func (executor *Executor) CalcSwapId(randomNumberHash ec.Hash, sender string, senderOtherChain string) ([]byte, error) {
	bep2Addr, err := types.AccAddressFromBech32(sender)
	if err != nil {
		return nil, err
	}
	return sdkMsg.CalculateSwapID(randomNumberHash[:], bep2Addr, senderOtherChain), nil
}

func (executor *Executor) IsSameAddress(addrA string, addrB string) bool {
	return strings.ToLower(addrA) == strings.ToLower(addrB)
}

func (executor *Executor) GetStatus() (interface{}, error) {
	bnbStatus := &common.BnbStatus{}
	bnbBalance, err := executor.Balance()
	if err != nil {
		return nil, err
	}
	bnbStatus.Balance = bnbBalance
	return bnbStatus, nil
}

func (executor *Executor) GetBalanceAlertMsg() (string, error) {
	if executor.Config.BnbBalanceAlertThreshold == 0 && executor.Config.TokenBalanceAlertThreshold == 0 {
		return "", nil
	}

	account, err := executor.RpcClient.GetAccount(executor.Config.DeputyAddr)
	if err != nil {
		return "", err
	}

	if account == nil || account.GetCoins() == nil {
		return "", errors.New("get nil account")
	}

	bnbBalance := account.GetCoins().AmountOf(common.BNBSymbol)
	tokenBalance := account.GetCoins().AmountOf(executor.Config.Symbol)

	alertMsg := ""
	if executor.Config.BnbBalanceAlertThreshold > 0 && bnbBalance < executor.Config.BnbBalanceAlertThreshold {
		alertMsg = alertMsg + fmt.Sprintf("BNB balance(%d) is less than %d\n", bnbBalance, executor.Config.BnbBalanceAlertThreshold)
	}
	if executor.Config.TokenBalanceAlertThreshold > 0 && tokenBalance < executor.Config.TokenBalanceAlertThreshold {
		alertMsg = alertMsg + fmt.Sprintf("%s balance(%d) is less than %d", executor.Config.Symbol,
			tokenBalance, executor.Config.TokenBalanceAlertThreshold)
	}

	return alertMsg, nil
}

func (executor *Executor) SendAmount(address string, amount *big.Int) (string, error) {
	executor.mutex.Lock()
	defer executor.mutex.Unlock()

	keyManager, err := getKeyManager(executor.Config)
	if err != nil {
		return "", common.NewError(err, false)
	}
	executor.RpcClient.SetKeyManager(keyManager)
	defer executor.RpcClient.SetKeyManager(nil)

	recipient, err := types.AccAddressFromBech32(address)
	if err != nil {
		return "", common.NewError(err, false)
	}

	if !amount.IsInt64() {
		return "", common.NewError(
			fmt.Errorf("out amount(%s) is not int64", amount),
			false,
		)
	}
	outCoins := []types.Coin{{
		Denom:  executor.Config.Symbol,
		Amount: amount.Int64(),
	}}

	res, err := executor.RpcClient.SendToken([]sdkMsg.Transfer{{ToAddr: recipient, Coins: outCoins}}, rpc.Sync)
	if err != nil {
		return "", common.NewError(err, isInvalidSequenceError(err.Error()))
	}
	if res.Code != 0 {
		return "", common.NewError(errors.New(res.Log), isInvalidSequenceError(res.Log))
	}
	return res.Hash.String(), nil
}
