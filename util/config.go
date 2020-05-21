package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/binance-chain/go-sdk/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/kava-labs/cosmos-sdk/types"

	dc "github.com/binance-chain/bep3-deputy/common"
)

const KeyTypePrivateKey = "private_key"
const KeyTypeMnemonic = "mnemonic"
const KeyTypeAWSMnemonic = "aws_mnemonic"
const KeyTypeAWSPrivateKey = "aws_private_key"

type Config struct {
	DBConfig              *DBConfig              `json:"db_config"`
	AlertConfig           *AlertConfig           `json:"alert_config"`
	ChainConfig           *ChainConfig           `json:"chain_config"`
	LogConfig             *LogConfig             `json:"log_config"`
	InstrumentationConfig *InstrumentationConfig `json:"instrumentation_config"`
	AdminConfig           *AdminConfig           `json:"admin_config"`
	BnbConfig             *BnbConfig             `json:"bnb_config"`
	EthConfig             *EthConfig             `json:"eth_config"`
	KavaConfig            *KavaConfig            `json:"kava_config"`
}

type AlertConfig struct {
	TelegramBotId  string `json:"telegram_bot_id"`
	TelegramChatId string `json:"telegram_chat_id"`

	BnbBlockUpdateTimeOut        int64 `json:"bnb_block_update_time_out"`
	OtherChainBlockUpdateTimeOut int64 `json:"other_chain_block_update_time_out"`

	ReconciliationDiffAmount *big.Float `json:"reconciliation_diff_amount"` // real number without decimal
}

type InstrumentationConfig struct {
	Prometheus           bool   `json:"prometheus"`
	PrometheusListenAddr string `json:"prometheus_listen_addr"`
}

type ChainConfig struct {
	BnbConfirmNum                int64      `json:"bnb_confirm_num"`
	BnbAutoRetryNum              int        `json:"bnb_auto_retry_num"`
	BnbAutoRetryTimeout          int64      `json:"bnb_auto_retry_timeout"`
	BnbExpireHeightSpan          int64      `json:"bnb_expire_height_span"`
	BnbMinAcceptExpireHeightSpan int64      `json:"bnb_min_accept_expire_height_span"`
	BnbMinRemainHeight           int64      `json:"bnb_min_remain_height"`
	BnbMinSwapAmount             *big.Int   `json:"bnb_min_swap_amount"`
	BnbMaxSwapAmount             *big.Int   `json:"bnb_max_swap_amount"`
	BnbMaxDeputyOutAmount        *big.Int   `json:"bnb_max_deputy_out_amount"`
	BnbRatio                     *big.Float `json:"bnb_ratio"`
	BnbFixedFee                  *big.Int   `json:"bnb_fixed_fee"`
	BnbStartHeight               int64      `json:"bnb_start_height"`

	OtherChain                          string     `json:"other_chain"`
	OtherChainConfirmNum                int64      `json:"other_chain_confirm_num"`
	OtherChainDecimal                   int        `json:"other_chain_decimal"`
	OtherChainExpireHeightSpan          int64      `json:"other_chain_expire_height_span"`
	OtherChainAutoRetryNum              int        `json:"other_chain_auto_retry_num"`
	OtherChainAutoRetryTimeout          int64      `json:"other_chain_auto_retry_timeout"`
	OtherChainMinAcceptExpireHeightSpan int64      `json:"other_chain_min_accept_expire_height_span"`
	OtherChainMinRemainHeight           int64      `json:"other_chain_min_remain_height"`
	OtherChainMinSwapAmount             *big.Int   `json:"other_chain_min_swap_amount"`
	OtherChainMaxSwapAmount             *big.Int   `json:"other_chain_max_swap_amount"`
	OtherChainMaxDeputyOutAmount        *big.Int   `json:"other_chain_max_deputy_out_amount"`
	OtherChainRatio                     *big.Float `json:"other_chain_ratio"`
	OtherChainFixedFee                  *big.Int   `json:"other_chain_fixed_fee"`
	OtherChainStartHeight               int64      `json:"other_chain_start_height"`
}

func (cfg *ChainConfig) Validate() {
	if cfg.BnbConfirmNum <= 0 {
		panic("bnb_confirm_num should be larger than 0")
	}
	if cfg.BnbAutoRetryNum <= 0 {
		panic("bnb_auto_retry_num should be larger than 0")
	}
	if cfg.BnbAutoRetryTimeout <= 0 {
		panic("bnb_auto_retry_timeout should be larger than 0")
	}
	if cfg.BnbExpireHeightSpan <= 0 {
		panic("bnb_expire_height_span should be larger than 0")
	}
	if cfg.BnbMinAcceptExpireHeightSpan <= 0 {
		panic("bnb_min_accept_expire_height should be larger than 0")
	}
	if cfg.BnbMinRemainHeight <= 0 {
		panic("bnb_min_remain_height should be larger than 0")
	}
	if cfg.BnbMinSwapAmount.Cmp(big.NewInt(0)) < 0 {
		panic("bnb_min_swap_amount should be no less than 0")
	}
	if cfg.BnbMaxSwapAmount.Cmp(big.NewInt(0)) <= 0 {
		panic("bnb_max_swap_amount should be larger than 0")
	}
	if cfg.BnbMinSwapAmount.Cmp(cfg.BnbMaxSwapAmount) >= 0 {
		panic("bnb_min_swap_amount should be less than bnb_max_swap_amount")
	}
	if cfg.BnbMaxDeputyOutAmount.Cmp(big.NewInt(0)) <= 0 {
		panic("bnb_max_deputy_out_amount should be larger than 0")
	}
	if cfg.BnbRatio.Cmp(big.NewFloat(0)) <= 0 {
		panic("bnb_ratio should larger than 0")
	}
	if cfg.BnbFixedFee.Cmp(big.NewInt(0)) < 0 {
		panic("bnb_fixed_fee should be no less than 0")
	}

	if cfg.OtherChain != dc.ChainEth && cfg.OtherChain != dc.ChainKava {
		panic(fmt.Sprintf("other chain only supports %s, %s", dc.ChainEth, dc.ChainKava))
	}
	if cfg.OtherChainConfirmNum <= 0 {
		panic("other_chain_confirm_num should be larger than 0")
	}
	if cfg.OtherChainDecimal <= 0 {
		panic("other_chain_decimal should be larger than 0")
	}
	if cfg.OtherChainExpireHeightSpan <= 0 {
		panic("other_chain_expire_height_span should be larger than 0")
	}
	if cfg.OtherChainAutoRetryNum <= 0 {
		panic("other_chain_auto_retry_num should be larger than 0")
	}
	if cfg.OtherChainAutoRetryTimeout <= 0 {
		panic("other_chain_auto_retry_timeout should be larger than 0")
	}
	if cfg.OtherChainMinAcceptExpireHeightSpan <= 0 {
		panic("other_chain_min_accept_expire_height should be larger than 0")
	}
	if cfg.OtherChainMinRemainHeight <= 0 {
		panic("other_chain_min_remain_height should be larger than 0")
	}
	if cfg.OtherChainMinSwapAmount.Cmp(big.NewInt(0)) < 0 {
		panic("other_chain_min_swap_amount should be no less than 0")
	}
	if cfg.OtherChainMaxSwapAmount.Cmp(big.NewInt(0)) <= 0 {
		panic("other_chain_max_swap_amount should be larger than 0")
	}
	if cfg.OtherChainMinSwapAmount.Cmp(cfg.OtherChainMaxSwapAmount) >= 0 {
		panic("other_chain_min_swap_amount should be less than other_chain_max_swap_amount")
	}
	if cfg.OtherChainMaxDeputyOutAmount.Cmp(big.NewInt(0)) <= 0 {
		panic("other_chain_max_deputy_out_amount should be larger than 0")
	}
	if cfg.OtherChainRatio.Cmp(big.NewFloat(0)) <= 0 {
		panic("other_chain_ratio should be larger than 0")
	}
	if cfg.OtherChainFixedFee.Cmp(big.NewInt(0)) < 0 {
		panic("other_chain_fixed_fee should be no less than 0")
	}
}

type DBConfig struct {
	Dialect                 string `json:"dialect"`
	DBPath                  string `json:"db_path"`
	MaxBnbKeptBlockHeight   int64  `json:"max_bnb_kept_block_height"`
	MaxOtherKeptBlockHeight int64  `json:"max_other_kept_block_height"`
}

func (cfg *DBConfig) Validate() {
	if cfg.Dialect != dc.DBDialectMysql && cfg.Dialect != dc.DBDialectSqlite3 {
		panic(fmt.Sprintf("only %s and %s supported", dc.DBDialectMysql, dc.DBDialectSqlite3))
	}
	if cfg.DBPath == "" {
		panic("db path should not be empty")
	}
	if cfg.MaxBnbKeptBlockHeight <= 0 {
		panic(fmt.Sprintf("max_bnb_kept_block_height should be larger than 0"))
	}
	if cfg.MaxOtherKeptBlockHeight <= 0 {
		panic(fmt.Sprintf("max_eth_kept_block_height should be larger than 0"))
	}
}

type AdminConfig struct {
	ListenAddr string `json:"listen_addr"`
}

func (cfg *AdminConfig) Validate() {
	if cfg.ListenAddr == "" {
		panic("listen address should not be empty")
	}
}

type BnbConfig struct {
	KeyType                    string           `json:"key_type"`
	AWSRegion                  string           `json:"aws_region"`
	AWSSecretName              string           `json:"aws_secret_name"`
	Mnemonic                   string           `json:"mnemonic"`
	RpcAddr                    string           `json:"rpc_addr"`
	Symbol                     string           `json:"symbol"`
	FetchInterval              int64            `json:"fetch_interval"`
	TokenBalanceAlertThreshold int64            `json:"token_balance_alert_threshold"`
	BnbBalanceAlertThreshold   int64            `json:"bnb_balance_alert_threshold"`
	DeputyAddr                 types.AccAddress `json:"deputy_addr"`
}

func (cfg *BnbConfig) Validate() {
	if cfg.KeyType == "" {
		panic(fmt.Sprintf("key_type of binance chain should not be empty"))
	}
	if cfg.KeyType != KeyTypeMnemonic && cfg.KeyType != KeyTypeAWSMnemonic {
		panic(fmt.Sprintf("key_type of binance chain only supports %s and %s", KeyTypeMnemonic, KeyTypeAWSMnemonic))
	}
	if cfg.KeyType == KeyTypeAWSMnemonic && cfg.AWSRegion == "" {
		panic(fmt.Sprintf("aws_region of binance chain should not be empty"))
	}
	if cfg.KeyType == KeyTypeAWSMnemonic && cfg.AWSSecretName == "" {
		panic(fmt.Sprintf("aws_secret_name of binance chain should not be empty"))
	}
	if cfg.RpcAddr == "" {
		panic(fmt.Sprintf("rpc address of binance chain should not be empty"))
	}
	if cfg.Symbol == "" {
		panic(fmt.Sprintf("symbol of binance chain should not be empty"))
	}
	if len(cfg.DeputyAddr) != types.AddrLen {
		panic(fmt.Sprintf("length of deputy address should be %d", types.AddrLen))
	}
	if cfg.FetchInterval <= 0 {
		panic(fmt.Sprintf("fetch_interval of binance chain should be larger than 0"))
	}
}

type EthConfig struct {
	SwapType                       string         `json:"swap_type"`
	KeyType                        string         `json:"key_type"`
	AWSRegion                      string         `json:"aws_region"`
	AWSSecretName                  string         `json:"aws_secret_name"`
	PrivateKey                     string         `json:"private_key"`
	Provider                       string         `json:"provider"`
	SwapContractAddr               common.Address `json:"swap_contract_addr"`
	TokenContractAddr              common.Address `json:"token_contract_addr"`
	DeputyAddr                     common.Address `json:"deputy_addr"`
	TokenBalanceAlertThreshold     *big.Int       `json:"token_balance_alert_threshold"`
	EthBalanceAlertThreshold       *big.Int       `json:"eth_balance_alert_threshold"`
	AllowanceBalanceAlertThreshold *big.Int       `json:"allowance_balance_alert_threshold"`
	FetchInterval                  int64          `json:"fetch_interval"`
	GasLimit                       int64          `json:"gas_limit"`
	GasPrice                       *big.Int       `json:"gas_price"`
}

func (cfg *EthConfig) Validate() {
	if cfg.SwapType == "" {
		panic("swap_type of ethereum should not be empty")
	}
	if cfg.SwapType != dc.EthSwapTypeEth && cfg.SwapType != dc.EthSwapTypeErc20 {
		panic(fmt.Sprintf("swap_type of ethereum only support %s and %s", dc.EthSwapTypeEth, dc.EthSwapTypeErc20))
	}
	if cfg.Provider == "" {
		panic(fmt.Sprintf("provider address of ethereum should not be empty"))
	}

	if cfg.KeyType == "" {
		panic(fmt.Sprintf("key_type ethereum should not be empty"))
	}
	if cfg.KeyType != KeyTypePrivateKey && cfg.KeyType != KeyTypeAWSPrivateKey {
		panic(fmt.Sprintf("key_type of ethereum only supports %s and %s", KeyTypePrivateKey, KeyTypeAWSPrivateKey))
	}
	if cfg.KeyType == KeyTypeAWSPrivateKey && cfg.AWSRegion == "" {
		panic(fmt.Sprintf("aws_region of ethereum should not be empty"))
	}
	if cfg.KeyType == KeyTypeAWSPrivateKey && cfg.AWSSecretName == "" {
		panic(fmt.Sprintf("aws_secret_name of ethereum should not be empty"))
	}

	var emptyAddr common.Address
	if cfg.SwapContractAddr.String() == emptyAddr.String() {
		panic(fmt.Sprintf("swap_contract_addrs of ethereum should not be empty"))
	}
	if cfg.SwapType == dc.EthSwapTypeErc20 && cfg.TokenContractAddr.String() == emptyAddr.String() {
		panic(fmt.Sprintf("token_contract_addr of ethereum should not be empty"))
	}
	if cfg.DeputyAddr.String() == emptyAddr.String() {
		panic(fmt.Sprintf("deputy_addr of ethereum should not be empty"))
	}
	if cfg.GasLimit <= 0 {
		panic(fmt.Sprintf("gas_limit of ethereum should be larger than 0"))
	}
	if cfg.FetchInterval <= 0 {
		panic(fmt.Sprintf("fetch_interval of ethereum should be larger than 0"))
	}
	if cfg.GasPrice.Cmp(big.NewInt(0)) <= 0 {
		panic("gas_price should be larger than 0")
	}
}

type KavaConfig struct {
	KeyType                    string         `json:"key_type"`
	AWSRegion                  string         `json:"aws_region"`
	AWSSecretName              string         `json:"aws_secret_name"`
	Mnemonic                   string         `json:"mnemonic"`
	RpcAddr                    string         `json:"rpc_addr"`
	Symbol                     string         `json:"symbol"`
	FetchInterval              int64          `json:"fetch_interval"`
	TokenBalanceAlertThreshold int64          `json:"token_balance_alert_threshold"`
	KavaBalanceAlertThreshold  int64          `json:"kava_balance_alert_threshold"`
	DeputyAddr                 sdk.AccAddress `json:"deputy_addr"`
}

func (cfg *KavaConfig) Validate() {
	if cfg.KeyType == "" {
		panic(fmt.Sprintf("key_type of kava chain should not be empty"))
	}
	if cfg.KeyType != KeyTypeMnemonic && cfg.KeyType != KeyTypeAWSMnemonic {
		panic(fmt.Sprintf("key_type of kava chain only supports %s and %s", KeyTypeMnemonic, KeyTypeAWSMnemonic))
	}
	if cfg.KeyType == KeyTypeAWSMnemonic && cfg.AWSRegion == "" {
		panic(fmt.Sprintf("aws_region of kava chain should not be empty"))
	}
	if cfg.KeyType == KeyTypeAWSMnemonic && cfg.AWSSecretName == "" {
		panic(fmt.Sprintf("aws_secret_name of kava chain should not be empty"))
	}
	if cfg.RpcAddr == "" {
		panic(fmt.Sprintf("rpc address of kava chain should not be empty"))
	}
	if cfg.Symbol == "" {
		panic(fmt.Sprintf("symbol of kava chain should not be empty"))
	}
	if len(cfg.DeputyAddr) != types.AddrLen {
		panic(fmt.Sprintf("length of deputy address should be %d", types.AddrLen))
	}
	if cfg.FetchInterval <= 0 {
		panic(fmt.Sprintf("fetch_interval of kava chain should be larger than 0"))
	}
}

type LogConfig struct {
	Level                        string `json:"level"`
	Filename                     string `json:"filename"`
	MaxFileSizeInMB              int    `json:"max_file_size_in_mb"`
	MaxBackupsOfLogFiles         int    `json:"max_backups_of_log_files"`
	MaxAgeToRetainLogFilesInDays int    `json:"max_age_to_retain_log_files_in_days"`
	UseConsoleLogger             bool   `json:"use_console_logger"`
	UseFileLogger                bool   `json:"use_file_logger"`
	Compress                     bool   `json:"compress"`
}

func (cfg *LogConfig) Validate() {
	if cfg.UseFileLogger {
		if cfg.Filename == "" {
			panic("filename should not be empty if use file logger")
		}
		if cfg.MaxFileSizeInMB <= 0 {
			panic("max_file_size_in_mb should be larger than 0 if use file logger")
		}
		if cfg.MaxBackupsOfLogFiles <= 0 {
			panic("max_backups_off_log_files should be larger than 0 if use file logger")
		}
	}
}

func (cfg *Config) Validate() {
	cfg.DBConfig.Validate()
	cfg.ChainConfig.Validate()
	cfg.BnbConfig.Validate()
	cfg.AdminConfig.Validate()
	// Validate the secondary chain's config
	switch cfg.ChainConfig.OtherChain {
	case dc.ChainEth:
		cfg.EthConfig.Validate()
	case dc.ChainKava:
		cfg.KavaConfig.Validate()
	}
	cfg.LogConfig.Validate()
}

func ParseConfigFromJson(content string) *Config {
	var config Config
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		panic(err)
	}
	return &config
}

func ParseConfigFromFile(filePath string) *Config {
	bz, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.Unmarshal(bz, &config); err != nil {
		panic(err)
	}
	return &config
}
