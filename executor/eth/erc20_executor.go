package eth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	bt "github.com/binance-chain/go-sdk/common/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	dc "github.com/binance-chain/bep3-deputy/common"
	da "github.com/binance-chain/bep3-deputy/executor/eth/abi"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

var _ dc.Executor = &Erc20Executor{}

type Erc20Executor struct {
	Abi               abi.ABI
	Provider          string
	Config            *util.Config
	Client            *ethclient.Client
	SwapContractAddr  common.Address
	TokenContractAddr common.Address
}

func NewErc20Executor(provider string, contractAddress common.Address, cfg *util.Config) *Erc20Executor {
	contractAbi, err := abi.JSON(strings.NewReader(da.Erc20SwapABI))
	if err != nil {
		panic("marshal abi error")
	}

	client, err := ethclient.Dial(provider)
	if err != nil {
		panic("new eth client error")
	}

	privKey, err := getPrivateKey(cfg.EthConfig)
	if err != nil {
		panic(fmt.Sprintf("generate private key error, err=%s", err.Error()))
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("get public key error")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	if !bytes.Equal(cfg.EthConfig.DeputyAddr.Bytes(), fromAddress.Bytes()) {
		panic(fmt.Sprintf(
			"deputy address supplied in config (%s) does not match mnemonic (%s)",
			cfg.EthConfig.DeputyAddr, fromAddress,
		))
	}

	return &Erc20Executor{
		Provider: provider,

		Abi:               contractAbi,
		Client:            client,
		SwapContractAddr:  contractAddress,
		TokenContractAddr: cfg.EthConfig.TokenContractAddr,

		Config: cfg,
	}
}

func (executor *Erc20Executor) GetChain() string {
	return dc.ChainEth
}

func (executor *Erc20Executor) GetBlockAndTxs(height int64) (*dc.BlockAndTxLogs, error) {
	header, err := executor.Client.HeaderByNumber(context.Background(), big.NewInt(height))
	if err != nil {
		return nil, err
	}

	txLogs, err := executor.GetLogs(header.Hash())
	if err != nil {
		return nil, err
	}

	return &dc.BlockAndTxLogs{
		Height:          height,
		BlockHash:       header.Hash().String(),
		ParentBlockHash: header.ParentHash.String(),
		BlockTime:       int64(header.Time),
		TxLogs:          txLogs,
	}, nil
}

func (executor *Erc20Executor) GetFetchInterval() time.Duration {
	return time.Duration(executor.Config.EthConfig.FetchInterval) * time.Second
}

func (executor *Erc20Executor) GetHTLTEvent(swapId common.Hash) (*HTLTEvent, error) {
	topics := [][]common.Hash{{HTLTEventHash}, {}, {}, {swapId}}
	logs, err := executor.Client.FilterLogs(context.Background(), ethereum.FilterQuery{
		Topics: topics,
	})
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, fmt.Errorf("swap id does not exist, swap_id=%s", swapId.String())
	}

	event, err := ParseHTLTEvent(&executor.Abi, &logs[0])
	if err != nil {
		util.Logger.Errorf("parse event log error, er=%s", err.Error())
		return nil, err
	}

	htltEvent := event.(HTLTEvent)
	return &htltEvent, nil
}

func (executor *Erc20Executor) GetLogs(blockHash common.Hash) ([]*store.TxLog, error) {
	topics := [][]common.Hash{{ClaimEventHash, HTLTEventHash, RefundEventHash}}

	logs, err := executor.Client.FilterLogs(context.Background(), ethereum.FilterQuery{
		BlockHash: &blockHash,
		Topics:    topics,
		Addresses: []common.Address{executor.SwapContractAddr},
	})
	if err != nil {
		return nil, err
	}

	models := make([]*store.TxLog, 0, len(logs))

	for _, log := range logs {
		event, err := ParseEvent(&executor.Abi, &log)
		if err != nil {
			util.Logger.Errorf("parse event log error, er=%s", err.Error())
			continue
		}
		if event == nil {
			continue
		}

		txLog := event.ToTxLog()
		txLog.Chain = dc.ChainEth
		txLog.ContractAddr = log.Address.Hex()
		txLog.Height = int64(log.BlockNumber)
		txLog.BlockHash = log.BlockHash.Hex()
		txLog.TxHash = log.TxHash.Hex()
		txLog.Status = store.TxStatusInit

		models = append(models, txLog)
	}
	return models, nil
}

func (executor *Erc20Executor) GetHeight() (int64, error) {
	header, err := executor.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, nil
	}
	return header.Number.Int64(), nil
}

func (executor *Erc20Executor) HTLT(randomNumberHash common.Hash, timestamp int64, heightSpan int64, recipientAddr string, otherChainSenderAddr string, otherChainRecipientAddr string, outAmount *big.Int) (string, *dc.Error) {
	auth, err := executor.GetTransactor()
	if err != nil {
		return "", dc.NewError(err, false)
	}

	instance, err := da.NewErc20Swap(executor.SwapContractAddr, executor.Client)
	if err != nil {
		return "", dc.NewError(err, false)
	}

	recvAddr := common.HexToAddress(recipientAddr)
	bep2RecipientAddr, err := bt.AccAddressFromBech32(otherChainRecipientAddr)
	if err != nil {
		return "", dc.NewError(err, false)
	}

	var bep2SenderAddr bt.AccAddress
	if otherChainSenderAddr != "" {
		bep2SenderAddr, err = bt.AccAddressFromBech32(otherChainSenderAddr)
		if err != nil {
			return "", dc.NewError(err, false)
		}
	}

	tx, err := instance.Htlt(auth, randomNumberHash, uint64(timestamp), big.NewInt(heightSpan), recvAddr,
		common.BytesToAddress(bep2SenderAddr), common.BytesToAddress(bep2RecipientAddr), outAmount, big.NewInt(0))

	if err != nil {
		return "", dc.NewError(err, true)
	}
	util.Logger.Debugf("init tx sent: %s", tx.Hash().Hex())
	return tx.Hash().String(), nil
}

func (executor *Erc20Executor) Claim(swapId common.Hash, randomNumber common.Hash) (string, *dc.Error) {
	auth, err := executor.GetTransactor()
	if err != nil {
		return "", dc.NewError(err, false)
	}

	instance, err := da.NewErc20Swap(executor.SwapContractAddr, executor.Client)
	if err != nil {
		return "", dc.NewError(err, false)
	}

	tx, err := instance.Claim(auth, swapId, randomNumber)
	if err != nil {
		return "", dc.NewError(err, true)
	}

	return tx.Hash().String(), nil
}

func (executor *Erc20Executor) Refund(swapId common.Hash) (string, *dc.Error) {
	auth, err := executor.GetTransactor()
	if err != nil {
		return "", dc.NewError(err, false)
	}

	instance, err := da.NewErc20Swap(executor.SwapContractAddr, executor.Client)
	if err != nil {
		return "", dc.NewError(err, false)
	}

	tx, err := instance.Refund(auth, swapId)
	if err != nil {
		return "", dc.NewError(err, true)
	}

	return tx.Hash().String(), nil
}

func (executor *Erc20Executor) GetSentTxStatus(hash string) store.TxStatus {
	_, isPending, err := executor.Client.TransactionByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		return store.TxSentStatusNotFound
	}
	if isPending {
		return store.TxSentStatusPending
	}

	txReceipt, err := executor.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return store.TxSentStatusNotFound
	}

	if txReceipt.Status == types.ReceiptStatusFailed {
		return store.TxSentStatusFailed
	} else {
		return store.TxSentStatusSuccess
	}
}

func (executor *Erc20Executor) HasSwap(swapId common.Hash) (bool, error) {
	instance, err := da.NewErc20Swap(executor.SwapContractAddr, executor.Client)
	if err != nil {
		return false, err
	}

	return instance.IsSwapExist(nil, swapId)
}

func (executor *Erc20Executor) GetSwap(swapId common.Hash) (*dc.SwapRequest, error) {
	htltEvent, err := executor.GetHTLTEvent(swapId)
	if err != nil {
		return nil, err
	}

	return &dc.SwapRequest{
		Id:                  swapId,
		RandomNumberHash:    htltEvent.RandomNumberHash,
		ExpireHeight:        htltEvent.ExpireHeight.Int64(),
		SenderAddress:       htltEvent.MsgSender.String(),
		RecipientAddress:    htltEvent.RecipientAddr.String(),
		RecipientOtherChain: bt.AccAddress(htltEvent.Bep2Addr[:]).String(),
		OutAmount:           htltEvent.OutAmount,
	}, nil
}

func (executor *Erc20Executor) Refundable(swapId common.Hash) (bool, error) {
	instance, err := da.NewErc20Swap(executor.SwapContractAddr, executor.Client)
	if err != nil {
		return false, err
	}

	refundable, err := instance.Refundable(nil, swapId)
	return refundable, err
}

func (executor *Erc20Executor) Claimable(swapId common.Hash) (bool, error) {
	instance, err := da.NewErc20Swap(executor.SwapContractAddr, executor.Client)
	if err != nil {
		return false, err
	}

	claimable, err := instance.Claimable(nil, swapId)
	return claimable, err
}

func (executor *Erc20Executor) GetTransactor() (*bind.TransactOpts, error) {
	privateKey, err := getPrivateKey(executor.Config.EthConfig)
	if err != nil {
		return nil, err
	}

	nonce, err := executor.Client.PendingNonceAt(context.Background(), executor.Config.EthConfig.DeputyAddr)
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)                                 // in wei
	auth.GasLimit = uint64(executor.Config.EthConfig.GasLimit) // in units
	auth.GasPrice = executor.Config.EthConfig.GasPrice
	return auth, nil
}

func (executor *Erc20Executor) Allowance() (*big.Int, error) {
	instance, err := da.NewERC20(executor.TokenContractAddr, executor.Client)
	if err != nil {
		return nil, err
	}

	allowance, err := instance.Allowance(nil, executor.Config.EthConfig.DeputyAddr, executor.SwapContractAddr)
	return allowance, err
}

func (executor *Erc20Executor) GetBalance(addressString string) (*big.Int, error) {
	address := common.HexToAddress(addressString)
	return executor.Erc20Balance(address)
}

func (executor *Erc20Executor) Erc20Balance(address common.Address) (*big.Int, error) {
	instance, err := da.NewERC20(executor.TokenContractAddr, executor.Client)
	if err != nil {
		return nil, err
	}

	balance, err := instance.BalanceOf(nil, address)
	return balance, err
}

func (executor *Erc20Executor) EthBalance(address common.Address) (*big.Int, error) {
	return executor.Client.BalanceAt(context.Background(), address, nil)
}

func (executor *Erc20Executor) GetDeputyAddress() string {
	return executor.Config.EthConfig.DeputyAddr.String()
}

func (executor *Erc20Executor) GetColdWalletAddress() string {
	return executor.Config.EthConfig.ColdWalletAddr.String()
}

func (executor *Erc20Executor) CalcSwapId(randomNumberHash common.Hash, sender string, senderOtherChain string) ([]byte, error) {
	var bep2Addr = bt.AccAddress{}
	if senderOtherChain != "" {
		parsedAddr, err := bt.AccAddressFromBech32(senderOtherChain)
		if err != nil {
			return nil, err
		}

		bep2Addr = parsedAddr
	}

	return CalculateSwapID(randomNumberHash[:], common.FromHex(sender), bep2Addr[:]), nil
}

func (executor *Erc20Executor) IsSameAddress(addrA string, addrB string) bool {
	return bytes.Equal(common.FromHex(addrA), common.FromHex(addrB))
}

func (executor *Erc20Executor) GetStatus() (interface{}, error) {
	ethStatus := &dc.EthStatus{}

	allowance, err := executor.Allowance()
	if err != nil {
		return nil, err
	}
	ethStatus.Allowance = util.QuoBigInt(allowance, util.GetBigIntForDecimal(executor.Config.ChainConfig.OtherChainDecimal)).String()

	balance, err := executor.Erc20Balance(executor.Config.EthConfig.DeputyAddr)
	if err != nil {
		return nil, err
	}
	ethStatus.Erc20Balance = util.QuoBigInt(balance, util.GetBigIntForDecimal(executor.Config.ChainConfig.OtherChainDecimal)).String()

	ethBalance, err := executor.EthBalance(executor.Config.EthConfig.DeputyAddr)
	if err != nil {
		return nil, err
	}
	ethStatus.EthBalance = util.QuoBigInt(ethBalance, util.GetBigIntForDecimal(18)).String()

	return ethStatus, nil
}

func (executor *Erc20Executor) GetBalanceAlertMsg() (string, error) {
	if executor.Config.EthConfig.EthBalanceAlertThreshold.Cmp(big.NewInt(0)) == 0 &&
		executor.Config.EthConfig.TokenBalanceAlertThreshold.Cmp(big.NewInt(0)) == 0 &&
		executor.Config.EthConfig.AllowanceBalanceAlertThreshold.Cmp(big.NewInt(0)) == 0 {
		return "", nil
	}

	alertMsg := ""
	if executor.Config.EthConfig.EthBalanceAlertThreshold.Cmp(big.NewInt(0)) > 0 {
		ethBalance, err := executor.EthBalance(executor.Config.EthConfig.DeputyAddr)
		if err != nil {
			return "", err
		}

		if ethBalance.Cmp(executor.Config.EthConfig.EthBalanceAlertThreshold) < 0 {
			alertMsg = alertMsg + fmt.Sprintf("eth balance(%s) is less than %s\n",
				ethBalance.String(), executor.Config.EthConfig.EthBalanceAlertThreshold.String())
		}
	}

	if executor.Config.EthConfig.AllowanceBalanceAlertThreshold.Cmp(big.NewInt(0)) > 0 {
		allowance, err := executor.Allowance()
		if err != nil {
			return "", err
		}
		if allowance.Cmp(executor.Config.EthConfig.AllowanceBalanceAlertThreshold) < 0 {
			alertMsg = alertMsg + fmt.Sprintf("token allowance balance(%s) is less than %s\n",
				allowance.String(), executor.Config.EthConfig.AllowanceBalanceAlertThreshold.String())
		}
	}

	if executor.Config.EthConfig.TokenBalanceAlertThreshold.Cmp(big.NewInt(0)) > 0 {
		tokenBalance, err := executor.Erc20Balance(executor.Config.EthConfig.DeputyAddr)
		if err != nil {
			return "", err
		}
		if tokenBalance.Cmp(executor.Config.EthConfig.TokenBalanceAlertThreshold) < 0 {
			alertMsg = alertMsg + fmt.Sprintf("token balance(%s) is less than %s",
				tokenBalance.String(), executor.Config.EthConfig.TokenBalanceAlertThreshold.String())
		}
	}
	return alertMsg, nil
}

func (executor *Erc20Executor) SendAmount(address string, amount *big.Int) (string, error) {
	return "", fmt.Errorf("not implemented") // TODO
}
