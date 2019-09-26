package deputy

import (
	"fmt"
	"math/big"
	"time"

	"github.com/binance-chain/bep3-deputy/util"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/store"
)

func (deputy *Deputy) Status() (*common.DeputyStatus, error) {
	deputyStatus := &common.DeputyStatus{
		Mode: deputy.GetMode(),
	}

	otherCurrentBlockLog := deputy.GetCurrentBlockLog(deputy.OtherExecutor.GetChain())
	deputyStatus.OtherChainSyncHeight = otherCurrentBlockLog.Height
	deputyStatus.OtherChainLastBlockFetchedAt = time.Unix(otherCurrentBlockLog.CreateTime, 0)

	otherHeight, err := deputy.OtherExecutor.GetHeight()
	if err != nil {
		return nil, err
	}
	deputyStatus.OtherChainHeight = otherHeight

	bnbCurrentBlockLog := deputy.GetCurrentBlockLog(common.ChainBinance)
	deputyStatus.BnbSyncHeight = bnbCurrentBlockLog.Height
	deputyStatus.BnbChainLastBlockFetchedAt = time.Unix(bnbCurrentBlockLog.CreateTime, 0)

	bnbHeight, err := deputy.BnbExecutor.GetHeight()
	if err != nil {
		return nil, err
	}
	deputyStatus.BnbChainHeight = bnbHeight

	bnbStatus, err := deputy.BnbExecutor.GetStatus()
	if err != nil {
		return nil, err
	}
	deputyStatus.BnbStatus = bnbStatus

	otherStatus, err := deputy.OtherExecutor.GetStatus()
	if err != nil {
		return nil, err
	}
	deputyStatus.OtherChainStatus = otherStatus

	return deputyStatus, nil
}

func (deputy *Deputy) FailedSwaps(offset int, pageCount int) ([]*common.SwapStatus, int, error) {
	var totalCount = 0
	deputy.DB.Model(store.Swap{}).Where("status in (?)", []store.SwapStatus{
		store.SwapStatusOtherHTLTSentFailed, store.SwapStatusBEP2HTLTSentFailed, store.SwapStatusBEP2ClaimSentFailed,
		store.SwapStatusBEP2RefundSentFailed, store.SwapStatusOtherClaimSentFailed, store.SwapStatusOtherRefundSentFailed,
	}).Count(&totalCount)

	if totalCount <= offset {
		return nil, totalCount, fmt.Errorf("the number of total failed swaps is %d, the offset of query is %d", totalCount, offset)
	}

	failedSwaps := make([]*store.Swap, 0)
	deputy.DB.Where("status in (?)", []store.SwapStatus{
		store.SwapStatusOtherHTLTSentFailed, store.SwapStatusBEP2HTLTSentFailed, store.SwapStatusBEP2ClaimSentFailed,
		store.SwapStatusBEP2RefundSentFailed, store.SwapStatusOtherClaimSentFailed, store.SwapStatusOtherRefundSentFailed,
	}).Order("id desc").Offset(offset).Limit(pageCount).Find(&failedSwaps)

	swapStatuses := make([]*common.SwapStatus, 0, len(failedSwaps))
	for _, failedSwap := range failedSwaps {
		swapStatus := &common.SwapStatus{
			Id:               failedSwap.Id,
			Type:             failedSwap.Type,
			SenderAddr:       failedSwap.SenderAddr,
			ReceiverAddr:     failedSwap.ReceiverAddr,
			OtherChainAddr:   failedSwap.OtherChainAddr,
			InAmount:         failedSwap.InAmount,
			OutAmount:        failedSwap.OutAmount,
			RandomNumberHash: failedSwap.RandomNumberHash,
			ExpireHeight:     failedSwap.ExpireHeight,
			Height:           failedSwap.Height,
			Timestamp:        failedSwap.Timestamp,
			RandomNumber:     failedSwap.RandomNumber,
			Status:           failedSwap.Status,
		}

		txsSent := make([]*store.TxSent, 0)
		deputy.DB.Where("random_number_hash = ?", failedSwap.RandomNumberHash).Order("id desc").Find(&txsSent)
		swapStatus.TxsSent = txsSent

		swapStatuses = append(swapStatuses, swapStatus)
	}
	return swapStatuses, totalCount, nil
}

func (deputy *Deputy) ResendTx(id int64) (string, error) {
	swap := &store.Swap{}
	deputy.DB.Where("id = ?", id).First(&swap)

	if swap.Id == 0 {
		return "", fmt.Errorf("swap %d does no exist", id)
	}

	var err error = nil
	var txHash = ""
	if swap.Status == store.SwapStatusOtherRefundSentFailed {
		txHash, err = deputy.sendOtherRefund(swap)
	} else if swap.Status == store.SwapStatusOtherClaimSentFailed {
		txHash, err = deputy.sendOtherClaim(swap)
	} else if swap.Status == store.SwapStatusOtherHTLTSentFailed {
		txHash, err = deputy.sendOtherHTLT(swap)
	} else if swap.Status == store.SwapStatusBEP2HTLTSentFailed {
		txHash, err = deputy.sendBEP2HTLT(swap)
	} else if swap.Status == store.SwapStatusBEP2ClaimSentFailed {
		txHash, err = deputy.sendBEP2Claim(swap)
	} else if swap.Status == store.SwapStatusBEP2RefundSentFailed {
		txHash, err = deputy.sendBEP2Refund(swap)
	} else {
		return "", fmt.Errorf("status of swap %d is %s, do not need to resend tx", id, swap.Status)
	}
	return txHash, err
}

func (deputy *Deputy) GetReconciliationStatus() (*common.ReconciliationStatus, error) {
	reconStatus := &common.ReconciliationStatus{}

	bnbChainTokenBalance, err := deputy.BnbExecutor.GetBalance()
	if err != nil {
		return nil, err
	}
	reconStatus.Bep2TokenBalance = util.QuoBigInt(bnbChainTokenBalance, util.GetBigIntForDecimal(common.BEP2Decimal))

	otherChainTokenBalance, err := deputy.OtherExecutor.GetBalance()
	if err != nil {
		return nil, err
	}
	reconStatus.OtherChainTokenBalance = util.QuoBigInt(otherChainTokenBalance, util.GetBigIntForDecimal(deputy.Config.ChainConfig.OtherChainDecimal))

	pendingBep2ToOtherSwaps := make([]*store.Swap, 0)
	deputy.DB.Where("type = ? and status in (?)", store.SwapTypeBEP2ToOther, []store.SwapStatus{
		store.SwapStatusOtherHTLTSent, store.SwapStatusOtherHTLTConfirmed, store.SwapStatusOtherClaimConfirmed, store.SwapStatusBEP2ClaimSentFailed,
		store.SwapStatusOtherHTLTExpired, store.SwapStatusOtherRefundSentFailed,
	}).Find(&pendingBep2ToOtherSwaps)
	pendingOtherOutAmount := big.NewInt(0)
	for _, swap := range pendingBep2ToOtherSwaps {
		deputyOutAmount := big.NewInt(0)
		deputyOutAmount.SetString(swap.DeputyOutAmount, 10)
		pendingOtherOutAmount = pendingOtherOutAmount.Add(pendingOtherOutAmount, deputyOutAmount)
	}
	reconStatus.OtherChainTokenOutPending = util.QuoBigInt(pendingOtherOutAmount, util.GetBigIntForDecimal(deputy.Config.ChainConfig.OtherChainDecimal))

	pendingOtherToBep2Swaps := make([]*store.Swap, 0)
	deputy.DB.Where("type = ? and status in (?)", store.SwapTypeOtherToBEP2, []store.SwapStatus{
		store.SwapStatusBEP2HTLTSent, store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusBEP2RefundSentFailed,
		store.SwapStatusBEP2ClaimConfirmed, store.SwapStatusOtherClaimSentFailed,
	}).Find(&pendingOtherToBep2Swaps)
	pendingBep2OutAmount := big.NewInt(0)
	for _, swap := range pendingOtherToBep2Swaps {
		deputyOutAmount := big.NewInt(0)
		deputyOutAmount.SetString(swap.DeputyOutAmount, 10)
		pendingBep2OutAmount = pendingBep2OutAmount.Add(pendingBep2OutAmount, deputyOutAmount)
	}
	reconStatus.Bep2TokenOutPending = util.QuoBigInt(pendingBep2OutAmount, util.GetBigIntForDecimal(common.BEP2Decimal))

	return reconStatus, nil
}
