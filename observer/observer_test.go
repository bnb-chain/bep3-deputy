package observer

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/require"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/executor/mock"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

func TestObserver_fetchBlock_FetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetBlockAndTxs(gomock.Any()).AnyTimes().Return(nil, errors.New("error"))

	ob := NewObserver(db, config, bnbChainExecutor, nil)

	err = ob.fetchBlock(bnbChainExecutor, 2, 3, "2")
	require.NotNil(t, err, "error should not be nil")

	require.Contains(t, err.Error(), "block info error")
}

func TestObserver_fetchBlock_WrongParentHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetBlockAndTxs(gomock.Any()).AnyTimes().Return(
		&common.BlockAndTxLogs{
			Height:          3,
			BlockHash:       "3",
			ParentBlockHash: "2_1",
			TxLogs:          nil,
		}, nil)

	ob := NewObserver(db, config, bnbChainExecutor, nil)

	blockLog1 := &store.BlockLog{
		Chain:      "BNB",
		BlockHash:  "1",
		Height:     1,
		ParentHash: "",
	}
	db.Create(blockLog1)

	blockLog2 := &store.BlockLog{
		Chain:      "BNB",
		BlockHash:  "2",
		Height:     2,
		ParentHash: "1",
	}
	db.Create(blockLog2)

	txLog2 := &store.TxLog{
		TxType: store.TxTypeBEP2Refund,
		Chain:  "BNB",
		Height: 2,
		TxHash: "tx_hash",
	}
	db.Create(txLog2)

	ob.fetchBlock(bnbChainExecutor, 2, 3, "2")

	deletedBlockLog := &store.BlockLog{}
	db.Where("height = ? and chain = ?", blockLog2.Height, "BNB").First(&deletedBlockLog)
	require.Equal(t, deletedBlockLog.Height, int64(0))

	deletedTxLog := &store.TxLog{}
	db.Where("height = ? and chain = ?", blockLog2.Height, "BNB").First(&deletedTxLog)
	require.Equal(t, deletedTxLog.Height, int64(0))
}

func TestObserver_fetchBlock_TrueParentHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetBlockAndTxs(gomock.Any()).AnyTimes().Return(
		&common.BlockAndTxLogs{
			Height:          3,
			BlockHash:       "3",
			ParentBlockHash: "2",
			TxLogs: []*store.TxLog{
				{
					TxType: store.TxTypeBEP2Refund,
					Chain:  "BNB",
					Height: 3,
					TxHash: "tx_hash_1",
				},
				{
					TxType: store.TxTypeBEP2Claim,
					Chain:  "BNB",
					Height: 3,
					TxHash: "tx_hash_2",
				},
			},
		}, nil)

	ob := NewObserver(db, config, bnbChainExecutor, nil)

	blockLog1 := &store.BlockLog{
		Chain:      "BNB",
		BlockHash:  "1",
		Height:     1,
		ParentHash: "",
	}
	db.Create(blockLog1)

	blockLog2 := &store.BlockLog{
		Chain:      "BNB",
		BlockHash:  "2",
		Height:     2,
		ParentHash: "1",
	}
	db.Create(blockLog2)

	txLog2 := &store.TxLog{
		TxType: store.TxTypeBEP2Refund,
		Chain:  "BNB",
		Height: 2,
		TxHash: "tx_hash",
	}
	db.Create(txLog2)

	err = ob.fetchBlock(bnbChainExecutor, 2, 3, "2")
	require.Nil(t, err, "error should be nil")

	newBlockLog := &store.BlockLog{}
	db.Where("height = ? and chain = ?", 3, "BNB").First(&newBlockLog)
	require.Equal(t, newBlockLog.Height, int64(3))

	newTxLogs := make([]*store.TxLog, 0)
	db.Where("height = ? and chain = ?", 3, "BNB").Find(&newTxLogs)
	require.Equal(t, len(newTxLogs), 2)
}
