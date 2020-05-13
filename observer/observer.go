package observer

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jinzhu/gorm"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

type Observer struct {
	DB *gorm.DB

	Config *util.Config

	OtherExecutor common.Executor
	BnbExecutor   common.Executor
}

func NewObserver(db *gorm.DB, cfg *util.Config, bnbExecutor common.Executor, otherExecutor common.Executor) *Observer {
	return &Observer{
		DB:            db,
		Config:        cfg,
		OtherExecutor: otherExecutor,
		BnbExecutor:   bnbExecutor,
	}
}

func (ob *Observer) Start() {
	go ob.fetch(ob.OtherExecutor, ob.Config.ChainConfig.OtherChainStartHeight)
	go ob.fetch(ob.BnbExecutor, ob.Config.ChainConfig.BnbStartHeight)

	go ob.Prune()
	go ob.Alert()
}

func (ob *Observer) fetch(executor common.Executor, startHeight int64) {
	for {
		curBlockLog := ob.GetCurrentBlockLog(executor.GetChain())
		util.Logger.Infof("%s cur height: %d", executor.GetChain(), curBlockLog.Height)

		nextHeight := curBlockLog.Height + 1
		if curBlockLog.Height == 0 && startHeight != 0 {
			nextHeight = startHeight
		}

		err := ob.fetchBlock(executor, curBlockLog.Height, nextHeight, curBlockLog.BlockHash)
		if err != nil {
			normalizedErr := strings.ToLower(err.Error())
			if strings.Contains(normalizedErr, "height must be less than or equal to the current blockchain height") ||
				strings.Contains(normalizedErr, "not found") {
				util.Logger.Infof("try to get ahead block, chain=%s, height=%d", executor.GetChain(), nextHeight)
			} else {
				util.Logger.Error(normalizedErr)
			}

			time.Sleep(executor.GetFetchInterval())
		}
	}
}

func (ob *Observer) fetchBlock(executor common.Executor, curHeight, nextHeight int64, curBlockHash string) error {
	blockAndTxLogs, err := executor.GetBlockAndTxs(nextHeight)
	if err != nil {
		return fmt.Errorf("get %s block info error, height=%d, err=%s", executor.GetChain(), nextHeight, err.Error())
	}

	parentHash := blockAndTxLogs.ParentBlockHash
	if curHeight != 0 && parentHash != curBlockHash {
		util.Logger.Infof("delete %s block at height %d, hash=%s", executor.GetChain(), curHeight, curBlockHash)
		return ob.DeleteBlockAndTxs(executor.GetChain(), curHeight)
	} else {
		nextBlockLog := store.BlockLog{
			Chain:      executor.GetChain(),
			BlockHash:  blockAndTxLogs.BlockHash,
			ParentHash: parentHash,
			Height:     blockAndTxLogs.Height,
			BlockTime:  blockAndTxLogs.BlockTime,
		}

		err := ob.SaveBlockAndTxs(executor.GetChain(), &nextBlockLog, blockAndTxLogs.TxLogs)
		if err != nil {
			return err
		}

		if util.PrometheusMetrics != nil {
			util.PrometheusMetrics.FetchedBlockHeight.With(prometheus.Labels{"chain": executor.GetChain()}).Set(float64(nextHeight))
		}

		err = ob.UpdateConfirmedNum(executor.GetChain(), nextBlockLog.Height)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ob *Observer) UpdateConfirmedNum(chain string, height int64) error {
	return ob.DB.Model(store.TxLog{}).Where("chain = ? and status = ?", chain, store.TxStatusInit).Updates(
		map[string]interface{}{
			"confirmed_num": gorm.Expr("? - height", height),
			"update_time":   time.Now().Unix(),
		}).Error
}

func (ob *Observer) SaveBlockAndTxs(chain string, blockLog *store.BlockLog, txLogs []*store.TxLog) error {
	tx := ob.DB.Begin()
	if err := tx.Error; err != nil {
		util.Logger.Errorf("begin tx error, err=%s", err)
		return err
	}

	if err := tx.Create(blockLog).Error; err != nil {
		util.Logger.Errorf("create block log error, err=%s", err)
		tx.Rollback()
		return err
	}

	for _, txLog := range txLogs {
		if err := tx.Create(txLog).Error; err != nil {
			util.Logger.Errorf("create tx log error, err=%s", err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (ob *Observer) DeleteBlockAndTxs(chain string, height int64) error {
	tx := ob.DB.Begin()
	if err := tx.Error; err != nil {
		util.Logger.Errorf("begin tx error, err=%s", err)
		return err
	}

	if err := tx.Where("height = ? and chain = ?", height, chain).Delete(store.BlockLog{}).Error; err != nil {
		util.Logger.Errorf("delete block log error, err=%s", err)
		tx.Rollback()
		return err
	}

	if err := tx.Where("height = ? and chain = ? and status = ?", height, chain, store.TxStatusInit).Delete(store.TxLog{}).Error; err != nil {
		util.Logger.Errorf("delete tx log error, err=%s", err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (ob *Observer) GetCurrentBlockLog(chain string) *store.BlockLog {
	blockLog := store.BlockLog{}
	ob.DB.Where("chain = ?", chain).Order("height desc").First(&blockLog)
	return &blockLog
}

func (ob *Observer) Prune() {
	for {
		// block log is for keeping track of block hash of blocks to prevent forks, so theoretically it would be good
		// if max block log number is larger than max possible fork height.
		curOtherChainBlockLog := ob.GetCurrentBlockLog(ob.OtherExecutor.GetChain())
		ob.DB.Where("chain = ? and height < ?", ob.OtherExecutor.GetChain(), curOtherChainBlockLog.Height-common.ObserverMaxBlockNumber).Delete(store.BlockLog{})
		if curOtherChainBlockLog.Height > 0 && ob.Config.DBConfig.MaxOtherKeptBlockHeight > 0 {
			ob.DB.Where("chain = ? and height < ?", ob.OtherExecutor.GetChain(), curOtherChainBlockLog.Height-ob.Config.DBConfig.MaxOtherKeptBlockHeight).Delete(store.TxLog{})
		}

		curBnbBlockLog := ob.GetCurrentBlockLog(ob.BnbExecutor.GetChain())
		ob.DB.Where("chain = ? and height < ?", ob.BnbExecutor.GetChain(), curBnbBlockLog.Height-common.ObserverMaxBlockNumber).Delete(store.BlockLog{})
		if curBnbBlockLog.Height > 0 && ob.Config.DBConfig.MaxBnbKeptBlockHeight > 0 {
			ob.DB.Where("chain = ? and height < ?", ob.BnbExecutor.GetChain(), curBnbBlockLog.Height-ob.Config.DBConfig.MaxBnbKeptBlockHeight).Delete(store.TxLog{})
		}

		time.Sleep(common.ObserverPruneInterval)
	}
}

func (ob *Observer) Alert() {
	for {
		curOtherChainBlockLog := ob.GetCurrentBlockLog(ob.OtherExecutor.GetChain())
		if curOtherChainBlockLog.Height > 0 {
			if time.Now().Unix()-curOtherChainBlockLog.CreateTime > ob.Config.AlertConfig.OtherChainBlockUpdateTimeOut {
				msg := fmt.Sprintf("chain %s last block fetched at %s, height=%d",
					ob.OtherExecutor.GetChain(), time.Unix(curOtherChainBlockLog.CreateTime, 0).String(), curOtherChainBlockLog.Height)
				util.SendTelegramMessage(ob.Config.AlertConfig.TelegramBotId, ob.Config.AlertConfig.TelegramChatId, msg)
			}
		}

		curBnbBlockLog := ob.GetCurrentBlockLog(ob.BnbExecutor.GetChain())
		if curBnbBlockLog.Height > 0 {
			if time.Now().Unix()-curBnbBlockLog.CreateTime > ob.Config.AlertConfig.BnbBlockUpdateTimeOut {
				msg := fmt.Sprintf("chain %s last block fetched at %s, height=%d",
					ob.BnbExecutor.GetChain(), time.Unix(curBnbBlockLog.CreateTime, 0).String(), curBnbBlockLog.Height)
				util.SendTelegramMessage(ob.Config.AlertConfig.TelegramBotId, ob.Config.AlertConfig.TelegramChatId, msg)
			}
		}

		time.Sleep(common.ObserverAlertInterval)
	}
}
