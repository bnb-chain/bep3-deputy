package eth

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	bnbTypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

type ContractType int

const HTLTEventName = "HTLT"
const ClaimEventName = "Claimed"

var HTLTEventHash = common.HexToHash("0xb3e26d98380491276a8dce9d38fd1049e89070230ff5f36ebb55ead64500ade1")
var ClaimEventHash = common.HexToHash("0x9f46b1606087bdf4183ec7dfdbe68e4ab9129a6a37901c16a7b320ae11a96018")
var RefundEventHash = common.HexToHash("0x04eb8ae268f23cfe2f9d72fa12367b104af16959f6a93530a4cc0f50688124f9")

type ContractEvent interface {
	ToTxLog() *store.TxLog
}

type ClaimEvent struct {
	MsgSender        common.Address
	RecipientAddr    common.Address
	SwapId           common.Hash
	RandomNumberHash common.Hash
	RandomNumber     common.Hash
}

func ParseClaimEvent(abi *abi.ABI, log *types.Log) (ContractEvent, error) {
	var ev ClaimEvent

	err := abi.Unpack(&ev, ClaimEventName, log.Data)
	if err != nil {
		return nil, err
	}

	ev.MsgSender = common.BytesToAddress(log.Topics[1].Bytes())
	ev.RecipientAddr = common.BytesToAddress(log.Topics[2].Bytes())
	ev.SwapId = common.BytesToHash(log.Topics[3].Bytes())

	util.Logger.Debugf("sender addr: %s", ev.MsgSender.String())
	util.Logger.Debugf("receiver addr: %s", ev.RecipientAddr.String())
	util.Logger.Debugf("swap id: %s", hex.EncodeToString(ev.SwapId[:]))
	util.Logger.Debugf("random number hash: %s", hex.EncodeToString(ev.RandomNumberHash[:]))
	util.Logger.Debugf("random number: %s", hex.EncodeToString(ev.RandomNumber[:]))
	return ev, nil
}

func (ev ClaimEvent) ToTxLog() *store.TxLog {
	return &store.TxLog{
		TxType:           store.TxTypeOtherClaim,
		SwapId:           hex.EncodeToString(ev.SwapId[:]),
		SenderAddr:       ev.MsgSender.Hex(),
		ReceiverAddr:     ev.RecipientAddr.Hex(),
		RandomNumber:     ev.RandomNumber.Hex(),
		RandomNumberHash: hex.EncodeToString(ev.RandomNumberHash[:]),
	}
}

type HTLTEvent struct {
	MsgSender        common.Address
	RecipientAddr    common.Address
	SwapId           common.Hash
	RandomNumberHash common.Hash
	Timestamp        uint64
	Bep2Addr         common.Address
	ExpireHeight     *big.Int
	OutAmount        *big.Int
	Bep2Amount       *big.Int
}

func ParseHTLTEvent(abi *abi.ABI, log *types.Log) (ContractEvent, error) {
	var ev HTLTEvent

	err := abi.Unpack(&ev, HTLTEventName, log.Data)
	if err != nil {
		return nil, err
	}

	ev.MsgSender = common.BytesToAddress(log.Topics[1].Bytes())
	ev.RecipientAddr = common.BytesToAddress(log.Topics[2].Bytes())
	ev.SwapId = common.BytesToHash(log.Topics[3].Bytes())

	util.Logger.Debugf("sender addr: %s", ev.MsgSender.String())
	util.Logger.Debugf("receiver addr: %s", ev.RecipientAddr.String())
	util.Logger.Debugf("swap id: %s", hex.EncodeToString(ev.SwapId[:]))
	util.Logger.Debugf("random number hash: %s", hex.EncodeToString(ev.RandomNumberHash[:]))
	util.Logger.Debugf("bep2 addr: %s", hex.EncodeToString(ev.Bep2Addr[:]))
	util.Logger.Debugf("timestamp: %d", ev.Timestamp)
	util.Logger.Debugf("expire height: %d", ev.ExpireHeight)
	util.Logger.Debugf("erc20 amount: %d", ev.OutAmount)
	util.Logger.Debugf("bep2 amount: %d", ev.Bep2Amount)
	return ev, nil
}

func (ev HTLTEvent) ToTxLog() *store.TxLog {
	return &store.TxLog{
		TxType:           store.TxTypeOtherHTLT,
		SwapId:           hex.EncodeToString(ev.SwapId[:]),
		SenderAddr:       ev.MsgSender.Hex(),
		ReceiverAddr:     ev.RecipientAddr.Hex(),
		OtherChainAddr:   bnbTypes.AccAddress(ev.Bep2Addr[:]).String(),
		RandomNumberHash: hex.EncodeToString(ev.RandomNumberHash[:]),
		Timestamp:        int64(ev.Timestamp),
		ExpireHeight:     ev.ExpireHeight.Int64(),
		InAmount:         ev.Bep2Amount.String(),
		OutAmount:        ev.OutAmount.String(),
	}
}

type RefundEvent struct {
	MsgSender        common.Address
	RecipientAddr    common.Address
	SwapId           common.Hash
	RandomNumberHash common.Hash
}

func ParseRefundEvent(log *types.Log) (ContractEvent, error) {
	var ev RefundEvent

	ev.MsgSender = common.BytesToAddress(log.Topics[1].Bytes())
	ev.RecipientAddr = common.BytesToAddress(log.Topics[2].Bytes())
	ev.SwapId = common.BytesToHash(log.Topics[3].Bytes())

	util.Logger.Debugf("sender addr: %s", ev.MsgSender.String())
	util.Logger.Debugf("swap id: %s", hex.EncodeToString(ev.SwapId[:]))
	util.Logger.Debugf("receiver addr: %s", ev.RecipientAddr.String())
	util.Logger.Debugf("random number hash: %s", hex.EncodeToString(ev.RandomNumberHash[:]))
	return ev, nil
}

func (ev RefundEvent) ToTxLog() *store.TxLog {
	return &store.TxLog{
		TxType:           store.TxTypeOtherRefund,
		SwapId:           hex.EncodeToString(ev.SwapId[:]),
		SenderAddr:       ev.MsgSender.Hex(),
		ReceiverAddr:     ev.RecipientAddr.Hex(),
		RandomNumberHash: hex.EncodeToString(ev.RandomNumberHash[:]),
	}
}

func ParseEvent(abi *abi.ABI, log *types.Log) (ContractEvent, error) {
	if bytes.Equal(log.Topics[0][:], ClaimEventHash[:]) {
		return ParseClaimEvent(abi, log)
	} else if bytes.Equal(log.Topics[0][:], HTLTEventHash[:]) {
		return ParseHTLTEvent(abi, log)
	} else if bytes.Equal(log.Topics[0][:], RefundEventHash[:]) {
		return ParseRefundEvent(log)
	}
	return nil, nil
}

func CalculateSwapID(randomNumberHash []byte, sender []byte, senderOtherChain []byte) []byte {
	data := randomNumberHash
	data = append(data, []byte(sender)...)
	data = append(data, []byte(senderOtherChain)...)
	return tmhash.Sum(data)
}

func getPrivateKey(config *util.EthConfig) (*ecdsa.PrivateKey, error) {
	var ethPrivateKey string
	if config.KeyType == util.KeyTypeAWSPrivateKey {
		awsPrivateKey, err := util.GetSecret(config.AWSSecretName, config.AWSRegion)
		if err != nil {
			return nil, err
		}
		ethPrivateKey = awsPrivateKey
	} else {
		ethPrivateKey = config.PrivateKey
	}

	privKey, err := crypto.HexToECDSA(ethPrivateKey)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}
