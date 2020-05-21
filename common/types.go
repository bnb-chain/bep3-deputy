package common

import (
	"encoding/json"
	"math"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/kava-labs/cosmos-sdk/types"

	"github.com/binance-chain/go-sdk/common/types"

	"github.com/binance-chain/bep3-deputy/store"
)

// chains supported now, if you want to add other chain, pls add chain name here
const ChainEth = "ETH"
const ChainBinance = "BNB"
const ChainKava = "KAVA"

const (
	DBDialectMysql   = "mysql"
	DBDialectSqlite3 = "sqlite3"
)

const BNBSymbol = "BNB"
const BEP2Decimal = 8

const KAVASymbol = "KAVA"

const EthSwapTypeEth = "eth_swap"
const EthSwapTypeErc20 = "erc20_swap"

type DeputyMode int

const (
	DeputyModeNormal       DeputyMode = 0
	DeputyModeStopSendHTLT DeputyMode = 1
)

func (mode DeputyMode) String() string {
	if mode == DeputyModeNormal {
		return "NormalMode"
	} else if mode == DeputyModeStopSendHTLT {
		return "StopSendHTLTMode"
	} else {
		return "UNKNOWN"
	}
}

func (mode DeputyMode) MarshalJSON() ([]byte, error) {
	var s = mode.String()
	return json.Marshal(s)
}

var Fixed8Decimals = big.NewInt(int64(math.Pow10(8)))

type BlockAndTxLogs struct {
	Height          int64
	BlockHash       string
	ParentBlockHash string
	BlockTime       int64
	TxLogs          []*store.TxLog
}

type DeputyStatus struct {
	Mode                         DeputyMode `json:"mode"`
	BnbChainHeight               int64      `json:"bnb_chain_height"`
	BnbSyncHeight                int64      `json:"bnb_sync_height"`
	OtherChainHeight             int64      `json:"other_chain_height"`
	OtherChainSyncHeight         int64      `json:"other_chain_sync_height"`
	BnbChainLastBlockFetchedAt   time.Time  `json:"bnb_chain_last_block_fetched_at"`
	OtherChainLastBlockFetchedAt time.Time  `json:"other_chain_last_block_fetched_at"`

	BnbStatus        interface{} `json:"bnb_status"`
	OtherChainStatus interface{} `json:"other_chain_status"`
}

type BnbStatus struct {
	Balance []types.TokenBalance `json:"balance"`
}

type KavaStatus struct {
	Balance sdk.Coins `json:"balance"`
}

type EthStatus struct {
	Allowance    string `json:"allowance"`
	Erc20Balance string `json:"erc20_balance"`
	EthBalance   string `json:"eth_balance"`
}

type SwapStatus struct {
	Id               int64            `json:"id"`
	Type             store.SwapType   `json:"type"`
	SenderAddr       string           `json:"sender_addr"`
	ReceiverAddr     string           `json:"receiver_addr"`
	OtherChainAddr   string           `json:"other_chain_addr"`
	InAmount         string           `json:"in_amount"`
	OutAmount        string           `json:"out_amount"`
	RandomNumberHash string           `json:"random_number_hash"`
	ExpireHeight     int64            `json:"expire_height"`
	Height           int64            `json:"height"`
	Timestamp        int64            `json:"timestamp"`
	RandomNumber     string           `json:"random_number"`
	Status           store.SwapStatus `json:"status"`
	TxsSent          []*store.TxSent  `json:"txs_sent"`
}

type FailedSwaps struct {
	TotalCount int           `json:"total_count"`
	CurPage    int           `json:"cur_page"`
	NumPerPage int           `json:"num_per_page"`
	Swaps      []*SwapStatus `json:"swaps"`
}

type ReconciliationStatus struct {
	Bep2TokenBalance          *big.Float
	Bep2TokenOutPending       *big.Float
	OtherChainTokenBalance    *big.Float
	OtherChainTokenOutPending *big.Float
}

type SwapRequest struct {
	Id                  common.Hash
	RandomNumberHash    common.Hash
	ExpireHeight        int64
	SenderAddress       string
	RecipientAddress    string
	RecipientOtherChain string
	OutAmount           *big.Int
}
