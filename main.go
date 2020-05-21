package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/binance-chain/go-sdk/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	sdk "github.com/kava-labs/cosmos-sdk/types"
	"github.com/kava-labs/go-sdk/client"
	app "github.com/kava-labs/go-sdk/kava"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/binance-chain/bep3-deputy/admin"
	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/deputy"
	"github.com/binance-chain/bep3-deputy/executor/bnb"
	"github.com/binance-chain/bep3-deputy/executor/eth"
	"github.com/binance-chain/bep3-deputy/executor/kava"
	"github.com/binance-chain/bep3-deputy/observer"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

const flagConfigType = "config-type"
const flagConfigAwsRegion = "aws-region"
const flagConfigAwsSecretKey = "aws-secret-key"
const flagConfigPath = "config-path"
const flagBnbNetwork = "bnb-network"
const flagKavaNetwork = "kava-network"

const ConfigTypeLocal = "local"
const ConfigTypeAws = "aws"

func printUsage() {
	fmt.Print("usage: ./deputy --bnb-network [0 for testnet, 1 for mainnet] --config-type [aws or local] --config-path config_file_path --aws-region region --aws-secret-key secret key\n")
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("could not create directory %v. %v", dir, err)
		}
	}
	return nil
}

func initFlags() {
	flag.String(flagConfigPath, "", "config file path")
	flag.String(flagConfigType, "", "config type, local or aws")
	flag.String(flagConfigAwsRegion, "", "aws s3 region")
	flag.String(flagConfigAwsSecretKey, "", "aws s3 secret key")
	flag.Int(flagBnbNetwork, int(types.TestNetwork), "binance chain network type")
	flag.Int(flagKavaNetwork, int(types.TestNetwork), "kava network type")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	initFlags()

	bnbNetwork := viper.GetInt(flagBnbNetwork)
	if bnbNetwork != int(types.TestNetwork) && bnbNetwork != int(types.ProdNetwork) {
		printUsage()
		return
	}
	// we set binance chain network type first because we need to parse binance chain address from config
	types.Network = types.ChainNetwork(bnbNetwork)

	// we set kava address prefixes first because we need to parse kava chain address from config
	kavaConfig := sdk.GetConfig()
	app.SetBech32AddressPrefixes(kavaConfig)
	kavaConfig.Seal()

	configType := viper.GetString(flagConfigType)
	if configType == "" {
		printUsage()
		return
	}

	if configType != ConfigTypeAws && configType != ConfigTypeLocal {
		printUsage()
		return
	}

	var config *util.Config

	// get config from aws s3
	if configType == ConfigTypeAws {
		awsSecretKey := viper.GetString(flagConfigAwsSecretKey)
		if awsSecretKey == "" {
			printUsage()
			return
		}

		awsRegion := viper.GetString(flagConfigAwsRegion)
		if awsRegion == "" {
			printUsage()
			return
		}

		configContent, err := util.GetSecret(awsSecretKey, awsRegion)
		if err != nil {
			fmt.Printf("get aws config error, err=%s", err.Error())
			return
		}
		config = util.ParseConfigFromJson(configContent)
	} else {
		configFilePath := viper.GetString(flagConfigPath)
		if configFilePath == "" {
			printUsage()
			return
		}
		config = util.ParseConfigFromFile(configFilePath)
	}
	config.Validate()

	util.InitLogger(*config.LogConfig)

	if config.InstrumentationConfig.Prometheus && config.InstrumentationConfig.PrometheusListenAddr != "" {
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			if err := http.ListenAndServe(config.InstrumentationConfig.PrometheusListenAddr, nil); err != http.ErrServerClosed {
				util.Logger.Error("Prometheus HTTP server ListenAndServe, err: %v", err)
			}
		}()
		util.MustRegisterMetrics()
	}

	if config.DBConfig.Dialect == common.DBDialectSqlite3 {
		err := ensureDir(filepath.Dir(config.DBConfig.DBPath))
		if err != nil {
			panic(err.Error())
		}
	}

	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.DBPath)
	if err != nil {
		panic(fmt.Sprintf("open db error, err=%s", err.Error()))
	}
	defer db.Close()

	// init db if tables do not exist
	store.InitTables(db)

	bnbExecutor := bnb.NewExecutor(config.BnbConfig.RpcAddr, types.ChainNetwork(bnbNetwork), config.BnbConfig)

	var otherExecutor common.Executor
	switch config.ChainConfig.OtherChain {
	case common.ChainEth:
		if config.EthConfig.SwapType == common.EthSwapTypeEth {
			otherExecutor = eth.NewEthExecutor(config.EthConfig.Provider, config.EthConfig.SwapContractAddr, config)
		} else {
			otherExecutor = eth.NewErc20Executor(config.EthConfig.Provider, config.EthConfig.SwapContractAddr, config)
		}
	case common.ChainKava:
		kavaNetwork := viper.GetInt(flagKavaNetwork)
		switch kavaNetwork {
		case int(client.LocalNetwork), int(client.TestNetwork), int(client.ProdNetwork):
			break
		default:
			printUsage()
			return
		}
		otherExecutor = kava.NewExecutor(config.KavaConfig.RpcAddr, client.ChainNetwork(kavaNetwork), config.KavaConfig)
	}

	dp := deputy.NewDeputy(db, config, bnbExecutor, otherExecutor)
	dp.Start()

	ob := observer.NewObserver(db, config, bnbExecutor, otherExecutor)
	ob.Start()

	adm := admin.NewAdmin(config, dp)
	go adm.Serve()

	select {}
}
