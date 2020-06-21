package deputy

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	ec "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

type Deputy struct {
	mtx  sync.RWMutex
	mode common.DeputyMode

	DB            *gorm.DB
	Config        *util.Config
	BnbExecutor   common.Executor
	OtherExecutor common.Executor

	lastReconStatus *common.ReconciliationStatus
}

func NewDeputy(db *gorm.DB, cfg *util.Config, bnbExecutor common.Executor, otherExecutor common.Executor) *Deputy {
	return &Deputy{
		DB:            db,
		Config:        cfg,
		BnbExecutor:   bnbExecutor,
		OtherExecutor: otherExecutor,
	}
}

func (deputy *Deputy) Start() {
	go deputy.ConfirmBEP2Tx()
	go deputy.ConfirmOtherTx()

	// swap from bep2 to other chain
	go deputy.BEP2SendHTLT()
	go deputy.BEP2SendClaim()
	go deputy.BEP2SendRefund()
	go deputy.BEP2ExpireUserHTLT()

	// swap from other chain to bep2
	go deputy.OtherSendHTLT()
	go deputy.OtherSendClaim()
	go deputy.OtherSendRefund()
	go deputy.OtherExpireUserHTLT()

	// common
	go deputy.CheckTxSentRoutine()
	go deputy.Alert()
	go deputy.ReconRoutine()

	// watch for funds accumulating in the hot wallet and move to cold wallet
	go deputy.RunBEP2HotWalletOverflow()
	go deputy.RunOtherHotWalletOverflow()

	// metrics
	if deputy.Config.InstrumentationConfig.Prometheus && deputy.Config.InstrumentationConfig.PrometheusListenAddr != "" {
		go deputy.Metrics()
	}
}

func (deputy *Deputy) Metrics() {
	for {
		var (
			totalBEP2ToOther    int64
			pendingBEP2ToOther  int64
			failedBE2P2ToOther  int64
			rejectedBEP2ToOther int64

			totalOtherToBEP2    int64
			pendingOtherToBEP2  int64
			failedOtherToBEP2   int64
			rejectedOtherToBEP2 int64
		)

		deputy.DB.Model(store.Swap{}).Where("type = ?", store.SwapTypeBEP2ToOther).Count(&totalBEP2ToOther)
		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?)", store.SwapTypeBEP2ToOther, []store.SwapStatus{
			store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSent, store.SwapStatusOtherHTLTConfirmed,
			store.SwapStatusOtherHTLTExpired, store.SwapStatusOtherRefundSent, store.SwapStatusOtherClaimConfirmed,
			store.SwapStatusBEP2ClaimSent,
		}).Count(&pendingBEP2ToOther)
		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?)", store.SwapTypeBEP2ToOther, []store.SwapStatus{
			store.SwapStatusOtherHTLTSentFailed, store.SwapStatusOtherRefundSentFailed, store.SwapStatusBEP2ClaimSentFailed,
		}).Count(&failedBE2P2ToOther)
		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?)", store.SwapTypeBEP2ToOther, []store.SwapStatus{
			store.SwapStatusRejected,
		}).Count(&rejectedBEP2ToOther)

		deputy.DB.Model(store.Swap{}).Where("type = ?", store.SwapTypeOtherToBEP2).Count(&totalOtherToBEP2)
		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?)", store.SwapTypeOtherToBEP2, []store.SwapStatus{
			store.SwapStatusOtherHTLTConfirmed, store.SwapStatusBEP2HTLTSent, store.SwapStatusBEP2HTLTConfirmed,
			store.SwapStatusBEP2HTLTExpired, store.SwapStatusBEP2RefundSent, store.SwapStatusBEP2ClaimConfirmed,
			store.SwapStatusOtherClaimSent,
		}).Count(&pendingOtherToBEP2)
		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?)", store.SwapTypeOtherToBEP2, []store.SwapStatus{
			store.SwapStatusBEP2HTLTSentFailed, store.SwapStatusBEP2RefundSentFailed, store.SwapStatusOtherClaimSentFailed,
		}).Count(&failedOtherToBEP2)
		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?)", store.SwapTypeOtherToBEP2, []store.SwapStatus{
			store.SwapStatusRejected,
		}).Count(&rejectedOtherToBEP2)

		if util.PrometheusMetrics != nil {
			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeBEP2ToOther), "status": "all"}).Set(float64(totalBEP2ToOther))
			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeBEP2ToOther), "status": "pending"}).Set(float64(pendingBEP2ToOther))
			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeBEP2ToOther), "status": "failed"}).Set(float64(failedBE2P2ToOther))
			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeBEP2ToOther), "status": "rejected"}).Set(float64(rejectedBEP2ToOther))

			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeOtherToBEP2), "status": "all"}).Set(float64(totalOtherToBEP2))
			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeOtherToBEP2), "status": "pending"}).Set(float64(pendingOtherToBEP2))
			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeOtherToBEP2), "status": "failed"}).Set(float64(failedOtherToBEP2))
			util.PrometheusMetrics.NumSwaps.With(prometheus.Labels{"type": string(store.SwapTypeOtherToBEP2), "status": "rejected"}).Set(float64(rejectedOtherToBEP2))
		}

		bnbBalance, err := deputy.BnbExecutor.GetBalance(deputy.BnbExecutor.GetDeputyAddress())
		if err == nil && util.PrometheusMetrics != nil {
			bnbBalanceBigFloat := util.QuoBigInt(bnbBalance, util.GetBigIntForDecimal(common.BEP2Decimal))
			bnbBalanceFloat, _ := bnbBalanceBigFloat.Float64()
			util.PrometheusMetrics.Balance.With(prometheus.Labels{"chain": string(deputy.BnbExecutor.GetChain())}).Set(bnbBalanceFloat)
		}

		otherBalance, err := deputy.OtherExecutor.GetBalance(deputy.OtherExecutor.GetDeputyAddress())
		if err == nil && util.PrometheusMetrics != nil {
			otherBalanceBigFloat := util.QuoBigInt(otherBalance, util.GetBigIntForDecimal(deputy.Config.ChainConfig.OtherChainDecimal))
			otherBalanceFloat, _ := otherBalanceBigFloat.Float64()
			util.PrometheusMetrics.Balance.With(prometheus.Labels{"chain": string(deputy.OtherExecutor.GetChain())}).Set(otherBalanceFloat)
		}

		time.Sleep(common.DeputyMetricsInterval)
	}
}

func (deputy *Deputy) ReconRoutine() {
	for {
		deputy.Recon()
		time.Sleep(common.DeputyReconInterval)
	}
}

func (deputy *Deputy) ConfirmOtherTx() {
	for {
		txLogs := make([]*store.TxLog, 0)
		deputy.DB.Where("chain = ? and status = ? and confirmed_num >= ?", deputy.OtherExecutor.GetChain(),
			store.TxStatusInit, deputy.Config.ChainConfig.OtherChainConfirmNum).Find(&txLogs)

		txHashes := make([]string, 0, len(txLogs))
		newSwaps := make([]*store.Swap, 0)

		for _, txLog := range txLogs {
			if txLog.TxType == store.TxTypeOtherHTLT {
				// reject swap request if receiver addr and other chain addr both are deputy addr
				if deputy.OtherExecutor.IsSameAddress(txLog.ReceiverAddr, deputy.OtherExecutor.GetDeputyAddress()) &&
					!deputy.BnbExecutor.IsSameAddress(txLog.OtherChainAddr, deputy.BnbExecutor.GetDeputyAddress()) {

					randomNumberHash := ec.HexToHash(txLog.RandomNumberHash)
					bnbChainSwapId, err := deputy.BnbExecutor.CalcSwapId(randomNumberHash, deputy.BnbExecutor.GetDeputyAddress(), txLog.SenderAddr)
					if err != nil {
						util.Logger.Errorf("calculate swap id parse error, random_number_hash=%s, sender_addr=%s, other_swap_id=%s", randomNumberHash.String(), txLog.SenderAddr, txLog.SwapId)
						continue
					}

					newSwap := &store.Swap{
						Type:             store.SwapTypeOtherToBEP2,
						OtherChainSwapId: txLog.SwapId,
						BnbChainSwapId:   hex.EncodeToString(bnbChainSwapId),
						SenderAddr:       txLog.SenderAddr,
						ReceiverAddr:     txLog.ReceiverAddr,
						OtherChainAddr:   txLog.OtherChainAddr,
						InAmount:         txLog.InAmount,
						OutAmount:        txLog.OutAmount,
						RandomNumberHash: txLog.RandomNumberHash,
						Timestamp:        txLog.Timestamp,
						ExpireHeight:     txLog.ExpireHeight,
						Height:           txLog.Height,
						Status:           store.SwapStatusOtherHTLTConfirmed,
					}
					newSwaps = append(newSwaps, newSwap)
				}
			}
			txHashes = append(txHashes, txLog.TxHash)
		}

		tx := deputy.DB.Begin()
		if tx.Error != nil {
			util.Logger.Errorf("begin tx error, err=%s", tx.Error)
			continue
		}

		err := tx.Model(store.TxLog{}).Where("tx_hash in (?)", txHashes).Updates(
			map[string]interface{}{
				"status":      store.TxStatusConfirmed,
				"update_time": time.Now().Unix(),
			}).Error
		if err != nil {
			util.Logger.Errorf("update tx status error, err=%s", err)
			tx.Rollback()
			continue
		}

		// create swap
		for _, swap := range newSwaps {
			err := tx.Create(swap).Error
			if err != nil {
				util.Logger.Errorf("create swap error, err=%s", err)
				tx.Rollback()
				continue
			}
		}

		for _, txLog := range txLogs {
			err := deputy.ConfirmTx(tx, txLog)
			if err != nil {
				util.Logger.Errorf("confirm tx error, err=%s", err)
				tx.Rollback()
				continue
			}
		}

		err = deputy.CompensateNewSwap(tx, deputy.BnbExecutor.GetChain(), newSwaps)
		if err != nil {
			util.Logger.Errorf("compensate new swap tx error, err=%s", err)
			tx.Rollback()
			continue
		}

		tx.Commit()

		time.Sleep(common.DeputyConfirmTxInterval)
	}
}

func (deputy *Deputy) ConfirmBEP2Tx() {
	for {
		txLogs := make([]*store.TxLog, 0)
		deputy.DB.Where("chain = ? and status = ? and confirmed_num >= ?", deputy.BnbExecutor.GetChain(),
			store.TxStatusInit, deputy.Config.ChainConfig.BnbConfirmNum).Find(&txLogs)

		txHashes := make([]string, 0, len(txLogs))
		newSwaps := make([]*store.Swap, 0)

		for _, txLog := range txLogs {
			if txLog.TxType == store.TxTypeBEP2HTLT {
				// reject swap request if receiver addr and other chain addr both are deputy addr
				if deputy.BnbExecutor.IsSameAddress(txLog.ReceiverAddr, deputy.BnbExecutor.GetDeputyAddress()) &&
					!deputy.OtherExecutor.IsSameAddress(txLog.OtherChainAddr, deputy.OtherExecutor.GetDeputyAddress()) &&
					strings.ToUpper(txLog.OutCoin) == strings.ToUpper(deputy.Config.BnbConfig.Symbol) {

					randomNumberHash := ec.HexToHash(txLog.RandomNumberHash)
					otherChainSwapId, err := deputy.OtherExecutor.CalcSwapId(randomNumberHash, deputy.OtherExecutor.GetDeputyAddress(), txLog.SenderAddr)
					if err != nil {
						util.Logger.Errorf("calculate swap id parse error, random_number_hash=%s, sender_addr=%s, bnb_swap_id=%s", randomNumberHash.String(), txLog.SenderAddr, txLog.SwapId)
						continue
					}

					newSwap := &store.Swap{
						Type:             store.SwapTypeBEP2ToOther,
						BnbChainSwapId:   txLog.SwapId,
						OtherChainSwapId: hex.EncodeToString(otherChainSwapId),
						SenderAddr:       txLog.SenderAddr,
						ReceiverAddr:     txLog.ReceiverAddr,
						OtherChainAddr:   txLog.OtherChainAddr,
						InAmount:         txLog.InAmount,
						OutAmount:        txLog.OutAmount,
						RandomNumberHash: txLog.RandomNumberHash,
						Timestamp:        txLog.Timestamp,
						ExpireHeight:     txLog.ExpireHeight,
						Height:           txLog.Height,
						Status:           store.SwapStatusBEP2HTLTConfirmed,
					}
					newSwaps = append(newSwaps, newSwap)
				}
			}

			txHashes = append(txHashes, txLog.TxHash)
		}

		tx := deputy.DB.Begin()
		if tx.Error != nil {
			util.Logger.Errorf("begin tx error, err=%s", tx.Error)
			continue
		}

		err := tx.Model(store.TxLog{}).Where("tx_hash in (?)", txHashes).Updates(
			map[string]interface{}{
				"status":      store.TxStatusConfirmed,
				"update_time": time.Now().Unix(),
			}).Error
		if err != nil {
			util.Logger.Errorf("update tx status error, err=%s", err)
			tx.Rollback()
			continue
		}

		// create swap
		for _, swap := range newSwaps {
			err := tx.Create(swap).Error
			if err != nil {
				util.Logger.Errorf("create swap error, err=%s", err)
				tx.Rollback()
				continue
			}
		}

		for _, txLog := range txLogs {
			err := deputy.ConfirmTx(tx, txLog)
			if err != nil {
				util.Logger.Errorf("confirm tx error, err=%s", err)
				tx.Rollback()
				continue
			}
		}

		err = deputy.CompensateNewSwap(tx, deputy.OtherExecutor.GetChain(), newSwaps)
		if err != nil {
			util.Logger.Errorf("compensate new swap tx error, err=%s", err)
			tx.Rollback()
			continue
		}

		tx.Commit()
		time.Sleep(common.DeputyConfirmTxInterval)
	}
}

func (deputy *Deputy) CompensateNewSwap(tx *gorm.DB, chain string, newSwaps []*store.Swap) error {
	for _, swap := range newSwaps {
		txLogs, err := deputy.GetConfirmedTxsLog(tx, chain, swap)
		if err != nil {
			return err
		}

		if len(txLogs) == 0 {
			continue
		}

		err = deputy.ConfirmTx(tx, txLogs[0])
		if err != nil {
			return err
		}
	}
	return nil
}

func (deputy *Deputy) ConfirmTx(tx *gorm.DB, txLog *store.TxLog) error {
	switch txLog.TxType {
	case store.TxTypeBEP2HTLT:
		return deputy.UpdateSwapStatusWhenConfirmTx(tx, "", txLog, []store.SwapStatus{
			store.SwapStatusOtherHTLTConfirmed, store.SwapStatusBEP2HTLTSent, store.SwapStatusBEP2HTLTSentFailed},
			nil, store.SwapStatusBEP2HTLTConfirmed, "")
	case store.TxTypeBEP2Claim:
		err := deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeOtherToBEP2, txLog, []store.SwapStatus{
			store.SwapStatusOtherHTLTConfirmed, store.SwapStatusBEP2HTLTSent, store.SwapStatusBEP2HTLTConfirmed},
			nil, store.SwapStatusBEP2ClaimConfirmed, txLog.RandomNumber)
		if err != nil {
			return err
		}

		err = deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeBEP2ToOther, txLog, nil,
			nil, store.SwapStatusBEP2ClaimConfirmed, txLog.RandomNumber)
		if err != nil {
			return err
		}
	case store.TxTypeBEP2Refund:
		err := deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeBEP2ToOther, txLog, nil,
			nil, store.SwapStatusBEP2RefundConfirmed, "")
		if err != nil {
			return err
		}

		err = deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeOtherToBEP2, txLog, nil,
			[]store.SwapStatus{store.SwapStatusOtherHTLTExpired, store.SwapStatusOtherRefundConfirmed},
			store.SwapStatusBEP2RefundConfirmed, "")
		if err != nil {
			return err
		}
	case store.TxTypeOtherHTLT:
		return deputy.UpdateSwapStatusWhenConfirmTx(tx, "", txLog, []store.SwapStatus{
			store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSent, store.SwapStatusOtherHTLTSentFailed},
			nil, store.SwapStatusOtherHTLTConfirmed, "")
	case store.TxTypeOtherClaim:
		err := deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeOtherToBEP2, txLog, nil,
			nil, store.SwapStatusOtherClaimConfirmed, txLog.RandomNumber)
		if err != nil {
			return err
		}

		err = deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeBEP2ToOther, txLog, []store.SwapStatus{
			store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSent, store.SwapStatusOtherHTLTConfirmed},
			nil, store.SwapStatusOtherClaimConfirmed, txLog.RandomNumber)
		if err != nil {
			return err
		}
	case store.TxTypeOtherRefund:
		err := deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeOtherToBEP2, txLog, nil,
			nil, store.SwapStatusOtherRefundConfirmed, txLog.RandomNumber)
		if err != nil {
			return err
		}

		err = deputy.UpdateSwapStatusWhenConfirmTx(tx, store.SwapTypeBEP2ToOther, txLog, nil,
			[]store.SwapStatus{store.SwapStatusBEP2HTLTExpired, store.SwapStatusBEP2RefundConfirmed},
			store.SwapStatusOtherRefundConfirmed, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func (deputy *Deputy) CheckTxSentRoutine() {
	for {
		deputy.CheckTxSent()

		time.Sleep(common.DeputyCheckTxSentInterval)
	}
}

func (deputy *Deputy) CheckTxSent() {
	txsSent := deputy.GetTxsSentByStatus([]store.TxStatus{store.TxSentStatusInit, store.TxSentStatusNotFound, store.TxSentStatusPending})

	for _, txSent := range txsSent {
		var status = store.TxSentStatusInit
		if txSent.Chain == deputy.OtherExecutor.GetChain() {
			status = deputy.OtherExecutor.GetSentTxStatus(txSent.TxHash)
		} else if txSent.Chain == deputy.BnbExecutor.GetChain() {
			status = deputy.BnbExecutor.GetSentTxStatus(txSent.TxHash)
		} else {
			util.Logger.Errorf("unexpected sent tx %s", txSent.TxHash)
			continue
		}

		if status == store.TxSentStatusFailed {
			errMsg := fmt.Sprintf("tx sent on chain %s failed, swap_id=%s, tx_type=%s, err_msg=%s",
				txSent.Chain, txSent.SwapId, txSent.Type, txSent.ErrMsg)
			deputy.sendTgMsg(errMsg)
		}

		deputy.UpdateTxSentStatus(txSent, status)
	}
}

func (deputy *Deputy) handleTxSent(swap *store.Swap, chain string, txType store.TxType,
	backwardStatus store.SwapStatus, failedStatus store.SwapStatus) {
	txsSent := deputy.GetTxsSentByType(chain, txType, swap)

	if len(txsSent) == 0 {
		deputy.UpdateSwapStatus(swap, backwardStatus, "")
		return
	}

	latestTx := txsSent[0]
	timeElapsed := time.Now().Unix() - latestTx.CreateTime
	autoRetryTimeout, autoRetryNum := deputy.getAutoRetryConfig(chain)
	txStatus := latestTx.Status

	if timeElapsed > autoRetryTimeout &&
		(txStatus == store.TxSentStatusNotFound ||
			txStatus == store.TxSentStatusInit ||
			txStatus == store.TxSentStatusPending) {

		if len(txsSent) >= autoRetryNum {
			deputy.UpdateSwapStatus(swap, failedStatus, "")
		} else {
			deputy.UpdateSwapStatus(swap, backwardStatus, "")
		}
		deputy.UpdateTxSentStatus(latestTx, store.TxSentStatusLost)
	} else if txStatus == store.TxSentStatusFailed {
		errMsg := fmt.Sprintf("tx on chain %s failed, marked swap status to %s, tx_hash=%s, bnb_swap_id=%s, other_chain_swap_id=%s, random_number_hash=%s",
			chain, failedStatus, latestTx.TxHash, swap.BnbChainSwapId, swap.OtherChainSwapId, swap.RandomNumberHash)
		deputy.sendTgMsg(errMsg)

		deputy.UpdateSwapStatus(swap, failedStatus, "")
	}
}

func (deputy *Deputy) sendTgMsg(msg string) {
	util.SendTelegramMessage(deputy.Config.AlertConfig.TelegramBotId, deputy.Config.AlertConfig.TelegramChatId, msg)
}

func (deputy *Deputy) getAutoRetryConfig(chain string) (int64, int) {
	var autoRetryTimeout int64
	var autoRetryNum int
	if chain == deputy.BnbExecutor.GetChain() {
		autoRetryTimeout = deputy.Config.ChainConfig.BnbAutoRetryTimeout
		autoRetryNum = deputy.Config.ChainConfig.BnbAutoRetryNum
	} else {
		autoRetryTimeout = deputy.Config.ChainConfig.OtherChainAutoRetryTimeout
		autoRetryNum = deputy.Config.ChainConfig.OtherChainAutoRetryNum
	}
	return autoRetryTimeout, autoRetryNum
}

func (deputy *Deputy) GetCurrentBlockLog(chain string) *store.BlockLog {
	blockLog := store.BlockLog{}
	deputy.DB.Where("chain = ?", chain).Order("height desc").First(&blockLog)
	return &blockLog
}

func (deputy *Deputy) Alert() {
	for {
		bnbAlertMsg, err := deputy.BnbExecutor.GetBalanceAlertMsg()
		if err != nil {
			util.Logger.Errorf("get bnb balance alert msg error, err=%s", err.Error())
		} else if bnbAlertMsg != "" {
			deputy.sendTgMsg(bnbAlertMsg)
		}

		otherChainAlertMsg, err := deputy.OtherExecutor.GetBalanceAlertMsg()
		if err != nil {
			util.Logger.Errorf("get chain %s alert msg error, err=%s", deputy.OtherExecutor.GetChain(), err.Error())
		} else if otherChainAlertMsg != "" {
			deputy.sendTgMsg(otherChainAlertMsg)
		}

		time.Sleep(common.DeputyAlertInterval)
	}
}

func (deputy *Deputy) SetMode(mode common.DeputyMode) {
	if mode != common.DeputyModeNormal && mode != common.DeputyModeStopSendHTLT {
		return
	}

	deputy.mtx.Lock()
	defer deputy.mtx.Unlock()

	deputy.mode = mode
}

func (deputy *Deputy) ShouldSendHTLT() bool {
	return deputy.GetMode() != common.DeputyModeStopSendHTLT
}

func (deputy *Deputy) GetMode() common.DeputyMode {
	deputy.mtx.RLock()
	defer deputy.mtx.RUnlock()

	return deputy.mode
}

func (deputy *Deputy) GetSwapsByTypeAndStatuses(swapType store.SwapType, statuses []store.SwapStatus) []*store.Swap {
	swaps := make([]*store.Swap, 0)
	deputy.DB.Where("type = ? and status in (?)", swapType, statuses).Find(&swaps)
	return swaps
}

func (deputy *Deputy) UpdateSwapStatus(swap *store.Swap, status store.SwapStatus, deputyOutAmount string) {
	toUpdate := map[string]interface{}{
		"status":      status,
		"update_time": time.Now().Unix(),
	}
	if deputyOutAmount != "" {
		toUpdate["deputy_out_amount"] = deputyOutAmount
	}
	deputy.DB.Model(swap).Updates(toUpdate)
}

func (deputy *Deputy) GetTxsSentByStatus(status []store.TxStatus) []*store.TxSent {
	txsSent := make([]*store.TxSent, 0)
	deputy.DB.Where("status in (?)", status).Find(&txsSent)
	return txsSent
}

func (deputy *Deputy) GetTxsSentByType(chain string, txType store.TxType, swap *store.Swap) []*store.TxSent {
	txsSent := make([]*store.TxSent, 0)

	query := deputy.DB.Where("chain = ? and type = ?", chain, txType)

	if chain == deputy.OtherExecutor.GetChain() {
		query = query.Where("swap_id = ?", swap.OtherChainSwapId)
	} else {
		query = query.Where("swap_id = ?", swap.BnbChainSwapId)
	}
	query.Order("id desc").Find(&txsSent)
	return txsSent
}

func (deputy *Deputy) UpdateTxSentStatus(txSent *store.TxSent, status store.TxStatus) {
	deputy.DB.Model(txSent).Updates(
		map[string]interface{}{
			"status":      status,
			"update_time": time.Now().Unix(),
		})
}

func (deputy *Deputy) GetTxLogByTxType(chain string, txType store.TxType, swap *store.Swap) *store.TxLog {
	txLog := &store.TxLog{}

	query := deputy.DB.Where("chain = ? and tx_type = ?", chain, txType)
	if chain == deputy.OtherExecutor.GetChain() {
		query = query.Where("swap_id = ?", swap.OtherChainSwapId)
	} else {
		query = query.Where("swap_id = ?", swap.BnbChainSwapId)
	}
	query.First(txLog)

	return txLog
}

func (deputy *Deputy) GetConfirmedTxsLog(tx *gorm.DB, chain string, swap *store.Swap) ([]*store.TxLog, error) {
	txLogs := make([]*store.TxLog, 0)
	query := tx.Where("chain = ? and status = ?", chain, store.TxStatusConfirmed)

	if chain == deputy.OtherExecutor.GetChain() {
		query = query.Where("swap_id = ?", swap.OtherChainSwapId)
	} else {
		query = query.Where("swap_id = ?", swap.BnbChainSwapId)
	}
	if err := query.Order("id desc").Find(&txLogs).Error; err != nil {
		return txLogs, err
	}
	return txLogs, nil
}

func (deputy *Deputy) UpdateSwapStatusWhenConfirmTx(tx *gorm.DB, swapType store.SwapType, txLog *store.TxLog,
	inStatuses []store.SwapStatus, notInStatuses []store.SwapStatus, updateStatus store.SwapStatus, randomNumber string) error {
	query := tx.Model(store.Swap{})

	if txLog.Chain == deputy.OtherExecutor.GetChain() {
		query = query.Where("other_chain_swap_id = ?", txLog.SwapId)
	} else {
		query = query.Where("bnb_chain_swap_id = ?", txLog.SwapId)
	}

	if swapType != "" {
		query = query.Where("type = ?", swapType)
	}

	if len(inStatuses) != 0 {
		query = query.Where("status in (?)", inStatuses)
	}

	if len(notInStatuses) != 0 {
		query = query.Where("status not in (?)", notInStatuses)
	}

	toUpdate := map[string]interface{}{
		"status":      updateStatus,
		"update_time": time.Now().Unix(),
	}

	if randomNumber != "" {
		toUpdate["random_number"] = randomNumber
	}
	return query.Updates(toUpdate).Error
}

func (deputy *Deputy) Recon() {
	reconStatus, err := deputy.GetReconciliationStatus()
	if err != nil {
		util.Logger.Errorf("get recon status error, err=%s", err.Error())
		return
	}

	if deputy.lastReconStatus == nil {
		deputy.lastReconStatus = reconStatus
		return
	}

	lastTotalAmount := big.NewFloat(0)
	lastTotalAmount.Add(lastTotalAmount, deputy.lastReconStatus.Bep2TokenBalance)
	lastTotalAmount.Add(lastTotalAmount, deputy.lastReconStatus.Bep2TokenOutPending)
	lastTotalAmount.Add(lastTotalAmount, deputy.lastReconStatus.OtherChainTokenBalance)
	lastTotalAmount.Add(lastTotalAmount, deputy.lastReconStatus.OtherChainTokenOutPending)

	latestTotalAmount := big.NewFloat(0)
	latestTotalAmount.Add(latestTotalAmount, reconStatus.Bep2TokenBalance)
	latestTotalAmount.Add(latestTotalAmount, reconStatus.Bep2TokenOutPending)
	latestTotalAmount.Add(latestTotalAmount, reconStatus.OtherChainTokenBalance)
	latestTotalAmount.Add(latestTotalAmount, reconStatus.OtherChainTokenOutPending)

	deputy.lastReconStatus = reconStatus

	diffAmount := new(big.Float).Sub(lastTotalAmount, latestTotalAmount)
	if diffAmount.Abs(diffAmount).Cmp(deputy.Config.AlertConfig.ReconciliationDiffAmount) >= 0 {
		reconMsg := fmt.Sprintf("bep2 token amount: %s\n", reconStatus.Bep2TokenBalance.String())
		reconMsg += fmt.Sprintf("bep2 token out pending amount: %s\n", reconStatus.Bep2TokenOutPending.String())
		reconMsg += fmt.Sprintf("other chain token amount: %s\n", reconStatus.OtherChainTokenBalance.String())
		reconMsg += fmt.Sprintf("other chain token out pending amount: %s\n", reconStatus.OtherChainTokenOutPending.String())
		reconMsg += fmt.Sprintf("total amount: %s\n", latestTotalAmount.String())
		reconMsg += fmt.Sprintf("diff from last amount: %s", diffAmount.String())
		deputy.sendTgMsg(reconMsg)
	}
}

func (deputy *Deputy) RunBEP2HotWalletOverflow() {
	executor := deputy.BnbExecutor
	for {
		time.Sleep(common.DeputyRunOverflowInterval)

		deputyBalance, err := executor.GetBalance(executor.GetDeputyAddress())
		if err != nil {
			util.Logger.Errorf("could not get BNB deputy balance: %w", err)
			continue
		}
		var overflow big.Int
		overflow.Sub(deputyBalance, deputy.Config.ChainConfig.BnbHotWalletOverflow)
		if overflow.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		_, err = executor.SendAmount(executor.GetColdWalletAddress(), &overflow)
		if err != nil {
			util.Logger.Errorf("BNB overflow tx failed: %w", err)
		}
	}
}
func (deputy *Deputy) RunOtherHotWalletOverflow() {
	executor := deputy.OtherExecutor
	for {
		time.Sleep(common.DeputyRunOverflowInterval)

		deputyBalance, err := executor.GetBalance(executor.GetDeputyAddress())
		if err != nil {
			util.Logger.Errorf("could not get OTHER deputy balance: %w", err)
			continue
		}
		var overflow big.Int
		overflow.Sub(deputyBalance, deputy.Config.ChainConfig.OtherChainHotWalletOverflow)
		if overflow.Cmp(big.NewInt(0)) <= 0 {
			continue
		}

		_, err = executor.SendAmount(executor.GetColdWalletAddress(), &overflow)
		if err != nil {
			util.Logger.Errorf("OTHER overflow tx failed: %w", err)
		}
	}
}
