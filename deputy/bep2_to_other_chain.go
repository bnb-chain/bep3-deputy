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

func (deputy *Deputy) BEP2SendHTLT() {
	for {
		swaps := deputy.GetSwapsByTypeAndStatuses(store.SwapTypeBEP2ToOther,
			[]store.SwapStatus{store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSent})

		for _, swap := range swaps {
			if swap.Status == store.SwapStatusBEP2HTLTConfirmed {
				_, err := deputy.sendOtherHTLT(swap)
				if err != nil {
					util.Logger.Error(err.Error())
				}
			} else {
				deputy.handleTxSent(swap, deputy.OtherExecutor.GetChain(), store.TxTypeOtherHTLT,
					store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSentFailed)
			}
		}

		time.Sleep(common.DeputySendTxInterval)
	}
}

func (deputy *Deputy) sendOtherHTLT(swap *store.Swap) (string, error) {
	if !deputy.ShouldSendHTLT() {
		return "", errors.New(fmt.Sprintf("current mode is %s, we should not send HTLT tx now", deputy.mode))
	}

	outAmount := big.NewInt(0)
	outAmount.SetString(swap.OutAmount, 10)

	if (swap.ExpireHeight-swap.Height < deputy.Config.ChainConfig.BnbMinAcceptExpireHeightSpan) ||
		(outAmount.Cmp(deputy.Config.ChainConfig.BnbMaxSwapAmount) > 0) ||
		(outAmount.Cmp(deputy.Config.ChainConfig.BnbMinSwapAmount) < 0) {

		// Reject swap request
		deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
		return "", fmt.Errorf("reject swap for wrong params, bnb_swap_id=%s, amount=%s remaining_height_span=%d",
			swap.BnbChainSwapId, outAmount.String(), swap.ExpireHeight-swap.Height)
	} else {
		bigIntDecimal := util.GetBigIntForDecimal(deputy.Config.ChainConfig.OtherChainDecimal)

		actualOutAmount := big.NewInt(1)
		actualOutAmount.Mul(outAmount, bigIntDecimal).Div(actualOutAmount, common.Fixed8Decimals)

		actualOutAmount = util.CalcActualOutAmount(actualOutAmount, deputy.Config.ChainConfig.OtherChainRatio,
			deputy.Config.ChainConfig.OtherChainFixedFee)

		// reject if params error
		if actualOutAmount.Cmp(big.NewInt(0)) <= 0 || actualOutAmount.Cmp(deputy.Config.ChainConfig.OtherChainMaxDeputyOutAmount) > 0 {
			deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
			return "", fmt.Errorf("reject swap for wrong out_amount, bnb_swap_id=%s, out_amount=%s",
				swap.BnbChainSwapId, actualOutAmount.String())
		}

		otherChainSwapId := ec.HexToHash(swap.OtherChainSwapId)
		bnbChainSwapId := ec.HexToHash(swap.BnbChainSwapId)

		isExist, err := deputy.OtherExecutor.HasSwap(otherChainSwapId)
		if err != nil {
			return "", fmt.Errorf("query chain %s swap error, other_chain_swap_id=%s, err=%s",
				deputy.OtherExecutor.GetChain(), swap.OtherChainSwapId, err.Error())
		} else if isExist {
			return "", fmt.Errorf("chain %s swap already exists, other_chain_swap_id=%s",
				deputy.OtherExecutor.GetChain(), swap.OtherChainSwapId)
		}

		bnbCurHeight, err := deputy.BnbExecutor.GetHeight()
		if err != nil {
			return "", fmt.Errorf("query binance chain current height error error, err=%s", err.Error())
		}

		// update status if height remaining in binance chain is not enough
		if swap.ExpireHeight-bnbCurHeight < deputy.Config.ChainConfig.BnbMinRemainHeight {
			deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
			return "", fmt.Errorf("reject swap for remaining binance chain height diff is not enough, current height=%d, expire height=%d, bnb_swap_id=%s",
				bnbCurHeight, swap.ExpireHeight, swap.BnbChainSwapId)
		}

		bnbSwapRequest, err := deputy.BnbExecutor.GetSwap(bnbChainSwapId)
		if err != nil {
			return "", fmt.Errorf("get bnb swap request error, err=%s", err.Error())
		}

		// check parameters against swap request on binance chain in case of corrupted database
		if bnbSwapRequest.OutAmount.String() != swap.OutAmount ||
			bnbSwapRequest.SenderAddress != swap.SenderAddr ||
			bnbSwapRequest.RecipientAddress != swap.ReceiverAddr ||
			bnbSwapRequest.RecipientOtherChain != swap.OtherChainAddr ||
			bnbSwapRequest.ExpireHeight != swap.ExpireHeight {

			deputy.UpdateSwapStatus(swap, store.SwapStatusRejected, "")
			return "", fmt.Errorf("reject swap for mismatch of parameters, sender_addr=%s, recipient_addr=%s, out_amount=%s",
				bnbSwapRequest.SenderAddress, bnbSwapRequest.RecipientAddress, bnbSwapRequest.OutAmount.String())
		}

		txSent := &store.TxSent{
			Chain:            deputy.OtherExecutor.GetChain(),
			Type:             store.TxTypeOtherHTLT,
			SwapId:           swap.OtherChainSwapId,
			RandomNumberHash: swap.RandomNumberHash,
		}

		randomNumberHash := ec.HexToHash(swap.RandomNumberHash)
		txHash, cmnErr := deputy.OtherExecutor.HTLT(randomNumberHash, swap.Timestamp, deputy.Config.ChainConfig.OtherChainExpireHeightSpan,
			swap.OtherChainAddr, swap.SenderAddr, deputy.BnbExecutor.GetDeputyAddress(), actualOutAmount)

		if cmnErr != nil {
			errMsg := fmt.Sprintf(
				"send chain %s HTLT tx error, bnb_swap_id=%s, is_retryable=%t, err=%s",
				deputy.OtherExecutor.GetChain(), swap.BnbChainSwapId, cmnErr.Retryable(), cmnErr.Error(),
			)
			deputy.sendTgMsg(errMsg)

			// is error retryable
			if !cmnErr.Retryable() {
				txSent.ErrMsg = cmnErr.Error()
				txSent.Status = store.TxSentStatusFailed
				// TODO should txHash be set? Should it always be returned from HTLT?
				// if not then CheckTxSent will try and lookup txs based on an empty hash, seems wrong
				deputy.UpdateSwapStatus(swap, store.SwapStatusOtherHTLTSentFailed, actualOutAmount.String())
				deputy.DB.Create(txSent)
			}
			return "", fmt.Errorf(errMsg)
		}
		util.Logger.Infof("send chain %s HTLT tx success, other_chain_swap_id=%s, tx_hash=%s", deputy.OtherExecutor.GetChain(),
			swap.OtherChainSwapId, txHash)

		txSent.TxHash = txHash

		deputy.UpdateSwapStatus(swap, store.SwapStatusOtherHTLTSent, actualOutAmount.String())
		deputy.DB.Create(txSent)
		return txHash, nil
	}
}

func (deputy *Deputy) BEP2SendClaim() {
	for {
		swaps := deputy.GetSwapsByTypeAndStatuses(store.SwapTypeBEP2ToOther,
			[]store.SwapStatus{store.SwapStatusOtherClaimConfirmed, store.SwapStatusBEP2ClaimSent})

		for _, swap := range swaps {
			if swap.Status == store.SwapStatusOtherClaimConfirmed {
				_, err := deputy.sendBEP2Claim(swap)
				if err != nil {
					util.Logger.Error(err.Error())
				}
			} else {
				deputy.handleTxSent(swap, deputy.BnbExecutor.GetChain(), store.TxTypeBEP2Claim,
					store.SwapStatusOtherClaimConfirmed, store.SwapStatusBEP2ClaimSentFailed)
			}
		}

		time.Sleep(common.DeputySendTxInterval)
	}
}

func (deputy *Deputy) sendBEP2Claim(swap *store.Swap) (string, error) {
	bnbSwapId := ec.HexToHash(swap.BnbChainSwapId)
	randomNumber := ec.HexToHash(swap.RandomNumber)

	claimable, err := deputy.BnbExecutor.Claimable(bnbSwapId)
	if err != nil {
		return "", fmt.Errorf("query bep2 swap error, err=%s", err.Error())
	}

	// if swap is not claimable, swap may expired or claimed, it would safe to update swap status to SwapStatusBEP2HTLTExpired,
	// for status will be updated when claim tx is confirmed.
	if !claimable {
		curBlock := deputy.GetCurrentBlockLog(deputy.BnbExecutor.GetChain())
		if curBlock.Height > swap.ExpireHeight {
			deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2HTLTExpired, "")
		}
		return "", fmt.Errorf("bep2 swap can not be claimed, bnb_swap_id=%s", swap.BnbChainSwapId)
	}

	txSent := &store.TxSent{
		Chain:            deputy.BnbExecutor.GetChain(),
		Type:             store.TxTypeBEP2Claim,
		SwapId:           swap.BnbChainSwapId,
		RandomNumberHash: swap.RandomNumberHash,
	}

	txHash, cmnErr := deputy.BnbExecutor.Claim(bnbSwapId, randomNumber)
	if cmnErr != nil {
		errMsg := fmt.Sprintf("send bep2 claim tx error, bnb_swap_id=%s, err=%s", swap.BnbChainSwapId, cmnErr.Error())
		// send alert msg if it is not Invalid sequence
		if !strings.Contains(errMsg, "Invalid sequence") {
			deputy.sendTgMsg(errMsg)
		}

		// is error retryable
		if !cmnErr.Retryable() {
			txSent.ErrMsg = cmnErr.Error()
			txSent.Status = store.TxSentStatusFailed
			deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2ClaimSentFailed, "")
			deputy.DB.Create(txSent)
		}
		return "", fmt.Errorf(errMsg)
	}
	util.Logger.Infof("send bep2 claim tx success, bnb_swap_id=%s, random_number=%s, tx_hash=%s",
		swap.BnbChainSwapId, randomNumber.Hex(), txHash)

	txSent.TxHash = txHash

	deputy.UpdateSwapStatus(swap, store.SwapStatusBEP2ClaimSent, "")
	deputy.DB.Create(txSent)
	return txHash, nil
}

func (deputy *Deputy) BEP2SendRefund() {
	for {
		swaps := deputy.GetSwapsByTypeAndStatuses(store.SwapTypeBEP2ToOther,
			[]store.SwapStatus{store.SwapStatusOtherHTLTConfirmed, store.SwapStatusOtherHTLTExpired, store.SwapStatusOtherRefundSent})

		for _, swap := range swaps {
			if swap.Status == store.SwapStatusOtherHTLTConfirmed {
				// htlt tx sent by deputy expired
				htltTx := deputy.GetTxLogByTxType(deputy.OtherExecutor.GetChain(), store.TxTypeOtherHTLT, swap)

				curBlock := deputy.GetCurrentBlockLog(deputy.OtherExecutor.GetChain())
				if curBlock.Height > htltTx.ExpireHeight {
					deputy.UpdateSwapStatus(swap, store.SwapStatusOtherHTLTExpired, "")
				}
			} else if swap.Status == store.SwapStatusOtherHTLTExpired {
				_, err := deputy.sendOtherRefund(swap)
				if err != nil {
					util.Logger.Error(err.Error())
				}
			} else if swap.Status == store.SwapStatusOtherRefundSent {
				deputy.handleTxSent(swap, deputy.OtherExecutor.GetChain(), store.TxTypeOtherRefund,
					store.SwapStatusOtherHTLTExpired, store.SwapStatusOtherRefundSentFailed)
			}
		}

		time.Sleep(common.DeputySendTxInterval)
	}
}

func (deputy *Deputy) sendOtherRefund(swap *store.Swap) (string, error) {
	otherChainSwapId := ec.HexToHash(swap.OtherChainSwapId)

	refundable, err := deputy.OtherExecutor.Refundable(otherChainSwapId)
	if err != nil {
		return "", fmt.Errorf("query chain %s swap error, other_chain_swap_id=%s, err=%s", deputy.OtherExecutor.GetChain(),
			swap.OtherChainSwapId, err.Error())
	} else if !refundable {
		return "", fmt.Errorf("chain %s swap is not refundable, other_chain_swap_id=%s",
			deputy.OtherExecutor.GetChain(), swap.OtherChainSwapId)
	}

	txSent := &store.TxSent{
		Chain:            deputy.OtherExecutor.GetChain(),
		Type:             store.TxTypeOtherRefund,
		SwapId:           swap.OtherChainSwapId,
		RandomNumberHash: swap.RandomNumberHash,
	}

	txHash, cmnErr := deputy.OtherExecutor.Refund(otherChainSwapId)
	if cmnErr != nil {
		errMsg := fmt.Sprintf("send chain %s refund tx error, other_chain_swap_id=%s, err=%s", deputy.OtherExecutor.GetChain(),
			swap.OtherChainSwapId, cmnErr.Error())
		deputy.sendTgMsg(errMsg)

		// is error retryable
		if !cmnErr.Retryable() {
			txSent.ErrMsg = cmnErr.Error()
			txSent.Status = store.TxSentStatusFailed
			deputy.UpdateSwapStatus(swap, store.SwapStatusOtherRefundSentFailed, "")
			deputy.DB.Create(txSent)
		}
		return "", fmt.Errorf(errMsg)
	}
	util.Logger.Infof("send chain %s refund tx success, other_chain_swap_id=%s, tx_hash=%s", deputy.OtherExecutor.GetChain(),
		swap.OtherChainSwapId, txHash)

	txSent.TxHash = txHash
	deputy.UpdateSwapStatus(swap, store.SwapStatusOtherRefundSent, "")
	deputy.DB.Create(txSent)
	return txHash, nil
}

func (deputy *Deputy) BEP2ExpireUserHTLT() {
	for {
		curBlock := deputy.GetCurrentBlockLog(deputy.BnbExecutor.GetChain())

		deputy.DB.Model(store.Swap{}).Where("type = ? and status in (?) and expire_height < ?",
			store.SwapTypeBEP2ToOther, []store.SwapStatus{store.SwapStatusOtherHTLTSentFailed,
				store.SwapStatusBEP2ClaimSentFailed, store.SwapStatusRejected}, curBlock.Height).Updates(
			map[string]interface{}{
				"status":      store.SwapStatusBEP2HTLTExpired,
				"update_time": time.Now().Unix(),
			})

		time.Sleep(common.DeputyExpireUserHTLTInterval)
	}
}
