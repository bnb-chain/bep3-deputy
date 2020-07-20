package deputy

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	ec "github.com/ethereum/go-ethereum/common"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

func (deputy *Deputy) OtherSendHTLT() {
	for {
		swaps := deputy.GetSwapsByTypeAndStatuses(store.SwapTypeOtherToBEP2,
			[]store.SwapStatus{store.SwapStatusOtherHTLTConfirmed, store.SwapStatusBEP2HTLTSent})

		for _, swap := range swaps {
			if swap.Status == store.SwapStatusOtherHTLTConfirmed {
				util.Logger.Info(fmt.Sprintf("attempting to send bnb chain HTLT for swap other_id=%s", swap.OtherChainSwapId))
				_, err := deputy.sendBEP2HTLT(swap)
				if err != nil {
					util.Logger.Error("submit bnb chain HTLT failed: %s", err)
				}
			} else {
				deputy.handleTxSent(swap, deputy.BnbExecutor.GetChain(), store.TxTypeBEP2HTLT,
					store.SwapStatusOtherHTLTConfirmed, store.SwapStatusBEP2HTLTSentFailed)
			}
		}

		time.Sleep(common.DeputySendTxInterval)
	}
}

func (deputy *Deputy) sendBEP2HTLT(swap *store.Swap) (string, error) {
	if !deputy.ShouldSendHTLT() {
		return "", fmt.Errorf("current mode is %s, we should not send HTLT tx now", deputy.mode)
	}

	outAmount := big.NewInt(0)
	outAmount.SetString(swap.OutAmount, 10)

	if (swap.ExpireHeight-swap.Height < deputy.Config.ChainConfig.OtherChainMinAcceptExpireHeightSpan) ||
		(outAmount.Cmp(deputy.Config.ChainConfig.OtherChainMaxSwapAmount) > 0) ||
		(outAmount.Cmp(deputy.Config.ChainConfig.OtherChainMinSwapAmount) < 0) {

		// Reject swap request
		deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
		errMsg := fmt.Sprintf(
			"set swap status to %s other_id=%s bnb_id=%s: height span too small or swap amount outside allowed range",
			store.SwapStatusRejected, swap.OtherChainSwapId, swap.BnbChainSwapId,
		)
		deputy.sendTgMsg(errMsg)
		return "", errors.New(errMsg)
	} else {
		bigIntDecimal := util.GetBigIntForDecimal(deputy.Config.ChainConfig.OtherChainDecimal)

		actualOutAmount := big.NewInt(1)
		actualOutAmount.Mul(outAmount, common.Fixed8Decimals).Div(actualOutAmount, bigIntDecimal)

		actualOutAmount = util.CalcActualOutAmount(actualOutAmount, deputy.Config.ChainConfig.BnbRatio,
			deputy.Config.ChainConfig.BnbFixedFee)

		// reject if params error
		if actualOutAmount.Cmp(big.NewInt(0)) <= 0 || actualOutAmount.Cmp(deputy.Config.ChainConfig.BnbMaxDeputyOutAmount) > 0 {
			deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
			errMsg := fmt.Sprintf(
				"set swap status to %s other_id=%s bnb_id=%s: transfer amount after fees subtracted is too small or big",
				store.SwapStatusRejected, swap.OtherChainSwapId, swap.BnbChainSwapId,
			)
			deputy.sendTgMsg(errMsg)
			return "", fmt.Errorf(errMsg)
		}

		bnbSwapId := ec.HexToHash(swap.BnbChainSwapId)
		otherChainSwapId := ec.HexToHash(swap.OtherChainSwapId)

		// do not send htlt tx if swap already exist or query failed
		isExist, err := deputy.BnbExecutor.HasSwap(bnbSwapId)
		if err != nil {
			return "", fmt.Errorf("query bep2 swap error, err=%s", err.Error())
		} else if isExist {
			return "", fmt.Errorf("bep2 swap already exists, bnb_swap_id=%s", swap.BnbChainSwapId)
		}

		otherChainCurHeight, err := deputy.OtherExecutor.GetHeight()
		if err != nil {
			return "", fmt.Errorf("query chain %s current height error, err=%s", deputy.OtherExecutor.GetChain(), err.Error())
		}

		// update status if height remaining in other chain is not enough
		if swap.ExpireHeight-otherChainCurHeight < deputy.Config.ChainConfig.OtherChainMinRemainHeight {
			deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
			errMsg := fmt.Sprintf(
				"set swap status to %s other_id=%s bnb_id=%s: not enough time left before swap expires",
				store.SwapStatusRejected, swap.OtherChainSwapId, swap.BnbChainSwapId,
			)
			deputy.sendTgMsg(errMsg)
			return "", errors.New(errMsg)
		}

		otherSwapRequest, err := deputy.OtherExecutor.GetSwap(otherChainSwapId)
		if err != nil {
			return "", fmt.Errorf("get other chain swap request error, err=%s", err.Error())
		}

		// check parameters against swap request on other chain in case of corrupted database
		if otherSwapRequest.OutAmount.String() != swap.OutAmount ||
			otherSwapRequest.SenderAddress != swap.SenderAddr ||
			otherSwapRequest.RecipientAddress != swap.ReceiverAddr ||
			otherSwapRequest.RecipientOtherChain != swap.OtherChainAddr ||
			otherSwapRequest.ExpireHeight != swap.ExpireHeight {

			deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
			errMsg := fmt.Sprintf(
				"set swap status to %s other_id=%s bnb_id=%s: bnb swap on chain doesn't match version in database",
				store.SwapStatusRejected, swap.OtherChainSwapId, swap.BnbChainSwapId,
			)
			deputy.sendTgMsg(errMsg)
			return "", errors.New(errMsg)
		}

		txSent := &store.TxSent{
			Chain:            deputy.BnbExecutor.GetChain(),
			Type:             store.TxTypeBEP2HTLT,
			SwapId:           swap.BnbChainSwapId,
			RandomNumberHash: swap.RandomNumberHash,
		}

		randomNumberHash := ec.HexToHash(swap.RandomNumberHash)
		txHash, cmnErr := deputy.BnbExecutor.HTLT(randomNumberHash, swap.Timestamp, deputy.Config.ChainConfig.BnbExpireHeightSpan,
			swap.OtherChainAddr, swap.SenderAddr, deputy.OtherExecutor.GetDeputyAddress(), actualOutAmount)
		if cmnErr != nil {
			// is error retryable
			if !cmnErr.Retryable() {
				txSent.ErrMsg = cmnErr.Error()
				txSent.Status = store.TxSentStatusFailed
				deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2HTLTSentFailed, actualOutAmount.String())
				deputy.sendTgMsg(fmt.Sprintf(
					"set swap status to %s other_id=%s bnb_id=%s: got non retryable error from sending htlt: %s",
					store.SwapStatusBEP2HTLTSentFailed, swap.OtherChainSwapId, swap.BnbChainSwapId, cmnErr.Error(),
				))
				deputy.DB.Create(txSent)
			}
			return "", fmt.Errorf("could not send HTLT: %w", cmnErr)
		}
		util.Logger.Infof("send bep2 HTLT tx success, bnb_swap_id=%s, tx_hash=%s", swap.BnbChainSwapId, txHash)

		txSent.TxHash = txHash

		deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2HTLTSent, actualOutAmount.String())
		deputy.DB.Create(txSent)
		return txHash, nil
	}
}

func (deputy *Deputy) OtherSendClaim() {
	for {
		swaps := deputy.GetSwapsByTypeAndStatuses(store.SwapTypeOtherToBEP2,
			[]store.SwapStatus{store.SwapStatusBEP2ClaimConfirmed, store.SwapStatusOtherClaimSent})

		for _, swap := range swaps {
			if swap.Status == store.SwapStatusBEP2ClaimConfirmed {
				util.Logger.Info(fmt.Sprintf("attempting to send other chain Claim for swap other_id=%s", swap.OtherChainSwapId))
				_, err := deputy.sendOtherClaim(swap)
				if err != nil {
					util.Logger.Error("submit other chain Claim failed: %s", err)
				}
			} else {
				deputy.handleTxSent(swap, deputy.OtherExecutor.GetChain(), store.TxTypeOtherClaim,
					store.SwapStatusBEP2ClaimConfirmed, store.SwapStatusOtherClaimSentFailed)
			}
		}

		time.Sleep(common.DeputySendTxInterval)
	}
}

func (deputy *Deputy) sendOtherClaim(swap *store.Swap) (string, error) {
	otherChainSwapId := ec.HexToHash(swap.OtherChainSwapId)
	randomNumber := ec.HexToHash(swap.RandomNumber)

	claimable, err := deputy.OtherExecutor.Claimable(otherChainSwapId)
	if err != nil {
		return "", fmt.Errorf("query chain %s swap error, other_chain_swap_id=%s, err=%s", deputy.OtherExecutor.GetChain(),
			otherChainSwapId.String(), err.Error())
	}

	// if swap is not claimable, swap may expired or claimed, it would safe to update swap status to SwapStatusOtherHTLTExpired,
	// for status will be updated when claim tx is confirmed.
	if !claimable {
		curBlock := deputy.GetCurrentBlockLog(deputy.OtherExecutor.GetChain())
		if curBlock.Height > swap.ExpireHeight {
			deputy.UpdateSwapStatus(swap, store.SwapStatusOtherHTLTExpired, "")
			deputy.sendTgMsg(fmt.Sprintf(
				"set swap status to %s other_id=%s bnb_id=%s: tried to send claim but other chain htlt expired",
				store.SwapStatusOtherHTLTExpired, swap.OtherChainSwapId, swap.BnbChainSwapId,
			))
		}
		return "", fmt.Errorf("chain %s swap is not claimable, other_chain_swap_id=%s", deputy.OtherExecutor.GetChain(),
			otherChainSwapId.String())
	}

	txSent := &store.TxSent{
		Chain:            deputy.OtherExecutor.GetChain(),
		Type:             store.TxTypeOtherClaim,
		SwapId:           swap.OtherChainSwapId,
		RandomNumberHash: swap.RandomNumberHash,
	}

	txHash, cmnErr := deputy.OtherExecutor.Claim(otherChainSwapId, randomNumber)
	if cmnErr != nil {
		// is error retryable
		if !cmnErr.Retryable() {
			txSent.ErrMsg = cmnErr.Error()
			txSent.Status = store.TxSentStatusFailed
			deputy.UpdateSwapStatus(swap, store.SwapStatusOtherClaimSentFailed, "")
			deputy.sendTgMsg(fmt.Sprintf(
				"set swap status to %s other_id=%s bnb_id=%s: got non retryable error from sending claim: %s",
				store.SwapStatusOtherClaimSentFailed, swap.OtherChainSwapId, swap.BnbChainSwapId, cmnErr.Error(),
			))
			deputy.DB.Create(txSent)
		}
		return "", fmt.Errorf("could not send Claim: %w", cmnErr)
	}
	util.Logger.Infof("send chain %s claim tx success, other_chain_swap_id=%s, tx_hash=%s", deputy.OtherExecutor.GetChain(),
		otherChainSwapId.String(), txHash)

	txSent.TxHash = txHash
	deputy.UpdateSwapStatus(swap, store.SwapStatusOtherClaimSent, "")
	deputy.DB.Create(txSent)
	return txHash, nil
}

func (deputy *Deputy) OtherSendRefund() {
	for {
		swaps := deputy.GetSwapsByTypeAndStatuses(store.SwapTypeOtherToBEP2,
			[]store.SwapStatus{store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusBEP2HTLTExpired, store.SwapStatusBEP2RefundSent})

		for _, swap := range swaps {
			if swap.Status == store.SwapStatusBEP2HTLTConfirmed {
				// htlt tx sent by deputy expired
				htltTx := deputy.GetTxLogByTxType(deputy.BnbExecutor.GetChain(), store.TxTypeBEP2HTLT, swap)

				curBlock := deputy.GetCurrentBlockLog(deputy.BnbExecutor.GetChain())
				if curBlock.Height > htltTx.ExpireHeight {
					deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2HTLTExpired, "")
				}
			} else if swap.Status == store.SwapStatusBEP2HTLTExpired {
				_, err := deputy.sendBEP2Refund(swap)
				if err != nil {
					util.Logger.Error(err.Error())
				}
			} else if swap.Status == store.SwapStatusBEP2RefundSent {
				deputy.handleTxSent(swap, deputy.BnbExecutor.GetChain(), store.TxTypeBEP2Refund,
					store.SwapStatusBEP2HTLTExpired, store.SwapStatusBEP2RefundSentFailed)
			}
		}

		time.Sleep(common.DeputySendTxInterval)
	}
}

func (deputy *Deputy) sendBEP2Refund(swap *store.Swap) (string, error) {
	bnbSwapId := ec.HexToHash(swap.BnbChainSwapId)

	refundable, err := deputy.BnbExecutor.Refundable(bnbSwapId)
	if err != nil {
		return "", fmt.Errorf("query bep2 swap error, err=%s", err.Error())
	} else if !refundable {
		return "", fmt.Errorf("bep2 swap can not be refund, random_number_hash=%s", bnbSwapId.String())
	}

	txSent := &store.TxSent{
		Chain:            deputy.BnbExecutor.GetChain(),
		Type:             store.TxTypeBEP2Refund,
		SwapId:           swap.BnbChainSwapId,
		RandomNumberHash: swap.RandomNumberHash,
	}

	txHash, cmnErr := deputy.BnbExecutor.Refund(bnbSwapId)
	if cmnErr != nil {
		errMsg := fmt.Sprintf("send bep2 refund tx error, bnb_swap_id=%s, err=%s", swap.BnbChainSwapId, cmnErr.Error())
		// send alert msg if it is not Invalid sequence
		if !strings.Contains(errMsg, "Invalid sequence") {
			deputy.sendTgMsg(errMsg)
		}
		// is error retryable
		if !cmnErr.Retryable() {
			txSent.ErrMsg = cmnErr.Error()
			txSent.Status = store.TxSentStatusFailed
			deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2RefundSentFailed, "")
			deputy.DB.Create(txSent)
		}
		return "", fmt.Errorf(errMsg)
	}
	util.Logger.Infof("send bep2 refund tx success, bnb_swap_id=%s, tx_hash=%s", swap.BnbChainSwapId, txHash)

	txSent.TxHash = txHash

	deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2RefundSent, "")
	deputy.DB.Create(txSent)
	return txHash, nil
}

func (deputy *Deputy) OtherExpireUserHTLT() {
	for {
		curBlock := deputy.GetCurrentBlockLog(deputy.OtherExecutor.GetChain())

		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?) and expire_height < ?",
			store.SwapTypeOtherToBEP2, []store.SwapStatus{store.SwapStatusBEP2HTLTSentFailed,
				store.SwapStatusOtherClaimSentFailed, store.SwapStatusRejected}, curBlock.Height).Updates(
			map[string]interface{}{
				"status":      store.SwapStatusOtherHTLTExpired,
				"update_time": time.Now().Unix(),
			})

		time.Sleep(common.DeputyCheckTxSentInterval)
	}
}
