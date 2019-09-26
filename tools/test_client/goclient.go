package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"time"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	dc "github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/executor/bnb"
	"github.com/binance-chain/bep3-deputy/executor/eth"
	"github.com/binance-chain/bep3-deputy/util"
)

const RandomNumberLength = 32

func GetRandomBytes() []byte {
	randBytes := make([]byte, RandomNumberLength)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(fmt.Sprintf("get random bytes err=%s", err.Error()))
	}
	return randBytes
}

func GetRandomBEPAddr() types.AccAddress {
	randBytes := make([]byte, RandomNumberLength)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(fmt.Sprintf("get random bep address err=%s", err.Error()))
	}
	addr := make([]byte, 20, 20)
	copy(addr[:], randBytes[:20])
	return addr
}

func CalculateRandomHash(randomNumber []byte, timestamp int64) common.Hash {
	randomNumberAndTimestamp := make([]byte, RandomNumberLength+8)
	copy(randomNumberAndTimestamp[:RandomNumberLength], randomNumber)
	binary.BigEndian.PutUint64(randomNumberAndTimestamp[RandomNumberLength:], uint64(timestamp))
	res := sha256.Sum256(randomNumberAndTimestamp)
	return res
}

func printUsage() {
	fmt.Print("usage: ./goclient --bnb-network [0 for testnet, 1 for mainnet] --config config_file_path\n")
}

func initFlags() {
	flag.String("config", "", "config file path")
	flag.Int("bnb-network", int(types.TestNetwork), "binance chain network type")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	initFlags()

	configFile := viper.GetString("config")
	if configFile == "" {
		printUsage()
		return
	}

	bnbNetwork := viper.GetInt("bnb-network")
	if bnbNetwork != int(types.TestNetwork) && bnbNetwork != int(types.ProdNetwork) {
		printUsage()
		return
	}
	types.Network = types.ChainNetwork(bnbNetwork)

	go func() {
		for {
			bep2(bnbNetwork, configFile)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			erc20(bnbNetwork, configFile)
			time.Sleep(2 * time.Second)
		}
	}()

	select {}
}

func getAddressFromPrivateKey(privateKey string) common.Address {
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(fmt.Sprintf("generate private key error, err=%s", err.Error()))
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("get public key error")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress
}

func getAddressFromMnemonic(mnemonic string) types.AccAddress {
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	if err != nil {
		panic(fmt.Sprintf("new key manager err, err=%s", err.Error()))
	}

	return keyManager.GetAddr()
}

func bep2(bnbNetwork int, configFile string) {
	types.Network = types.ChainNetwork(bnbNetwork)

	config := util.ParseConfigFromFile(configFile)

	var otherExecutor dc.Executor
	if config.EthConfig.SwapType == dc.EthSwapTypeEth {
		otherExecutor = eth.NewEthExecutor(config.EthConfig.Provider, config.EthConfig.SwapContractAddr, config)
	} else {
		otherExecutor = eth.NewErc20Executor(config.EthConfig.Provider, config.EthConfig.SwapContractAddr, config)
	}

	bnbExecutor := bnb.NewExecutor(config.BnbConfig.RpcAddr, types.ChainNetwork(bnbNetwork), config.BnbConfig)

	timestamp := time.Now().Unix()
	randomBytes := GetRandomBytes()
	randomKey := common.BytesToHash(randomBytes)
	randomHash := CalculateRandomHash(randomBytes, timestamp)
	address := getAddressFromPrivateKey(config.EthConfig.PrivateKey)

	swapId, err := otherExecutor.CalcSwapId(randomHash, address.String(), "")

	util.Logger.Infof("eth swap timestamp: %d", timestamp)
	util.Logger.Infof("eth swap randomNum: %s", hex.EncodeToString(randomBytes))
	util.Logger.Infof("eth swap randomHash: %s", hex.EncodeToString(randomHash[:]))
	util.Logger.Infof("eth swap swap id: %s", hex.EncodeToString(swapId))

	bepAddr := GetRandomBEPAddr()
	txHash, cmnErr := otherExecutor.HTLT(randomHash, timestamp, 1000, config.EthConfig.DeputyAddr.String(), "", bepAddr.String(), big.NewInt(100000000))
	if cmnErr != nil {
		util.Logger.Infof("init eth swap error, err=%s", cmnErr.Error())
		return
	}

	util.Logger.Infof("init tx sent: %s", txHash)

	util.Logger.Infof("init eth swap, swap_id=%s, randomHash=%s, randomKey=%s", hex.EncodeToString(swapId), randomHash.String(), randomKey.String())

	time.Sleep(1 * time.Minute)

	bnbSwapId := msg.CalculateSwapID(randomHash[:], config.BnbConfig.DeputyAddr, address.String())
	util.Logger.Infof("bnb swap id: %s", hex.EncodeToString(bnbSwapId))
	for {
		_, isExist, _ := bnbExecutor.QuerySwap(bnbSwapId[:])

		util.Logger.Infof("query bep2 exist, swap_id=%s", hex.EncodeToString(bnbSwapId))

		if isExist {
			tx, _ := bnbExecutor.Claim(common.BytesToHash(bnbSwapId), common.BytesToHash(randomBytes))
			if err != nil {
				println("sent bep2 claim error", err.Error())
				return
			}
			println("sent bep2 claim success", tx)
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func erc20(bnbNetwork int, configFile string) {
	types.Network = types.ChainNetwork(bnbNetwork)

	config := util.ParseConfigFromFile(configFile)

	var otherExecutor dc.Executor
	if config.EthConfig.SwapType == dc.EthSwapTypeEth {
		otherExecutor = eth.NewEthExecutor(config.EthConfig.Provider, config.EthConfig.SwapContractAddr, config)
	} else {
		otherExecutor = eth.NewErc20Executor(config.EthConfig.Provider, config.EthConfig.SwapContractAddr, config)
	}
	bnbExecutor := bnb.NewExecutor(config.BnbConfig.RpcAddr, types.ChainNetwork(bnbNetwork), config.BnbConfig)

	toOnOtherChain := common.HexToAddress("0x042ccc750E1099068622Bb521003F207297a40b0")
	randomNum := GetRandomBytes()
	ts := time.Now().Unix()
	randomHash := CalculateRandomHash(randomNum, ts)
	swapId, _ := bnbExecutor.CalcSwapId(randomHash, bnbExecutor.GetDeputyAddress(), "")
	senderAddr := getAddressFromMnemonic(config.BnbConfig.Mnemonic)

	timeSpan := 100000
	util.Logger.Infof("bnb swap randomHash: %s", hex.EncodeToString(randomHash[:]))
	util.Logger.Infof("bnb swap randomNum: %s", hex.EncodeToString(randomNum[:]))
	util.Logger.Infof("bnb swap swapId: %s", hex.EncodeToString(swapId[:]))

	txHash, err := bnbExecutor.HTLT(randomHash, ts, int64(timeSpan), config.BnbConfig.DeputyAddr.String(), "", toOnOtherChain.String(), big.NewInt(100000000))
	if err != nil {
		util.Logger.Infof("sent htlt tx error, err=%s", err.Error())
		return
	}

	otherChainSwapId, _ := otherExecutor.CalcSwapId(randomHash, config.EthConfig.DeputyAddr.String(), senderAddr.String())

	util.Logger.Infof("init bep2 swap, other_chain_swap_id=%s, randomHash=%s, randomKey=%s, txHash=%s", hex.EncodeToString(otherChainSwapId), hex.EncodeToString(randomHash[:]),
		hex.EncodeToString(randomNum), txHash)

	time.Sleep(1 * time.Minute)

	for {
		claimable, err := otherExecutor.Claimable(common.BytesToHash(otherChainSwapId))

		util.Logger.Infof("query eth claimable, swap_id=%s", common.BytesToHash(otherChainSwapId[:]).String())

		if err != nil {
			util.Logger.Infof("get claimable error", err.Error())
			continue
		}

		if claimable {
			txHash, _ := otherExecutor.Claim(common.BytesToHash(otherChainSwapId[:]), common.BytesToHash(randomNum))
			if err != nil {
				util.Logger.Infof("claim eth error", err.Error())
				continue
			}
			util.Logger.Infof("sent eth claim tx_hash=%s", txHash)
			return
		}
		time.Sleep(15 * time.Second)
	}
}
