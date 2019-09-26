package deputy

import (
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/require"

	"github.com/binance-chain/bep3-deputy/common"
	"github.com/binance-chain/bep3-deputy/executor/mock"
	"github.com/binance-chain/bep3-deputy/store"
	"github.com/binance-chain/bep3-deputy/util"
)

func TestDeputy_CompensateNewSwap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("other_chain")

	de := NewDeputy(db, config, nil, otherChainExecutor)
	swap := &store.Swap{
		Type:             store.SwapTypeOtherToBEP2,
		BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		Status:           store.SwapStatusOtherHTLTConfirmed,
	}
	de.DB.Create(swap)

	bep2HTLTTx := &store.TxLog{
		TxType:           store.TxTypeBEP2HTLT,
		SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		Status:           store.TxStatusConfirmed,
		Chain:            "BNB",
	}
	de.DB.Create(bep2HTLTTx)

	de.CompensateNewSwap(de.DB, "BNB", []*store.Swap{swap})
	updatedSwap := &store.Swap{}
	de.DB.Where("bnb_chain_swap_id = ?", swap.BnbChainSwapId).First(updatedSwap)
	require.Equal(t, updatedSwap.Status, store.SwapStatusBEP2HTLTConfirmed)
}

func TestDeputy_ConfirmTx(t *testing.T) {
	testCases := []struct {
		swap      *store.Swap
		txLog     *store.TxLog
		newStatus store.SwapStatus
	}{
		// TxTypeBEP2HTLT
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherHTLTConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2HTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2HTLTConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2HTLTSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2HTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2HTLTConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2HTLTSentFailed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2HTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2HTLTConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2ClaimConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2HTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2ClaimConfirmed,
		},
		// TxTypeBEP2Claim
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherHTLTConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Claim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2ClaimConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2HTLTSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Claim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2ClaimConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2HTLTConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Claim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2ClaimConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherClaimSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Claim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusOtherClaimSent,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2HTLTSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Claim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2ClaimConfirmed,
		},
		// TxTypeBEP2Refund
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherRefundConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Refund,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2RefundConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2RefundSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Refund,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusBEP2RefundConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				BnbChainSwapId:   "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherRefundConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeBEP2Refund,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "BNB",
			},
			newStatus: store.SwapStatusOtherRefundConfirmed,
		},
		//TxTypeOtherHTLT
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2HTLTConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherHTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherHTLTConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherHTLTSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherHTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherHTLTConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherHTLTSentFailed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherHTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherHTLTConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2ClaimConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherHTLT,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusBEP2ClaimConfirmed,
		},
		// TxTypeOtherClaim
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2HTLTConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherClaim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherClaimConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherHTLTSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherClaim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherClaimConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherHTLTConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherClaim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherClaimConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2ClaimConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherClaim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusBEP2ClaimConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherClaimSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherClaim,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherClaimConfirmed,
		},
		//TxTypeOtherRefund
		{
			swap: &store.Swap{
				Type:             store.SwapTypeOtherToBEP2,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2RefundConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherRefund,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherRefundConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusOtherRefundSent,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherRefund,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusOtherRefundConfirmed,
		},
		{
			swap: &store.Swap{
				Type:             store.SwapTypeBEP2ToOther,
				OtherChainSwapId: "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.SwapStatusBEP2RefundConfirmed,
			},
			txLog: &store.TxLog{
				TxType:           store.TxTypeOtherRefund,
				SwapId:           "4f637e501c0567ffe3f1895fa86bfb6a383799ece9bd3dfea4218ce369db7179",
				RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
				Status:           store.TxStatusConfirmed,
				Chain:            "OTHER_CHAIN",
			},
			newStatus: store.SwapStatusBEP2RefundConfirmed,
		},
	}

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)
	for _, testCase := range testCases {
		de.DB.Delete(store.Swap{})

		de.DB.Create(testCase.swap)
		de.ConfirmTx(de.DB, testCase.txLog)

		updatedSwap := &store.Swap{}
		if testCase.txLog.Chain == otherChainExecutor.GetChain() {
			de.DB.Where("other_chain_swap_id = ?", testCase.swap.OtherChainSwapId).First(updatedSwap)
		} else {
			de.DB.Where("bnb_chain_swap_id = ?", testCase.swap.BnbChainSwapId).First(updatedSwap)
		}
		require.Equal(t, updatedSwap.Status, testCase.newStatus)
	}
}

func TestDeputy_sendBEP2HTLT_Reject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	deputyAddressStr := "other_deputy_address"
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)
	otherChainExecutor.EXPECT().GetSentTxStatus(gomock.Any()).AnyTimes().Return(store.TxSentStatusSuccess)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "reject_expire_height",
		OtherChainSwapId: "other_swap_id_reject_expire_height",
		BnbChainSwapId:   "bnb_swap_id_reject_expire_height",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan - 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainMaxSwapAmount.Int64()), big.NewInt(-1)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	updateSwap := &store.Swap{}

	de.DB.Where("bnb_chain_swap_id = ?", swap.BnbChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)

	swap1 := &store.Swap{
		RandomNumberHash: "reject_swap_amount",
		OtherChainSwapId: "other_swap_id_reject_swap_amount",
		BnbChainSwapId:   "bnb_swap_id_reject_swap_amount",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainMaxSwapAmount.Int64()), big.NewInt(10)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap1)

	_, err = de.sendBEP2HTLT(swap1)
	require.NotNil(t, err, "error should not be nil")
	updateSwap = &store.Swap{}

	de.DB.Where("bnb_chain_swap_id = ?", swap1.BnbChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)
}

func TestDeputy_sendBEP2HTLT_WrongOutAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	deputyAddressStr := "other_deputy_address"
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, nil, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "reject_wrong_out_amount",
		OtherChainSwapId: "other_swap_id_reject_wrong_out_amount",
		BnbChainSwapId:   "bnb_swap_id_reject_wrong_out_amount",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(-1)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	updateSwap := &store.Swap{}

	de.DB.Where("bnb_chain_swap_id = ?", swap.BnbChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)

	swap1 := &store.Swap{
		RandomNumberHash: "reject_wrong_amount",
		OtherChainSwapId: "other_swap_id_reject_max_out_amount",
		BnbChainSwapId:   "bnb_swap_id_reject_max_out_amount",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainMaxSwapAmount.Int64()), big.NewInt(-1)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap1)

	config.ChainConfig.BnbMaxDeputyOutAmount = config.ChainConfig.BnbFixedFee

	_, err = de.sendBEP2HTLT(swap1)
	require.NotNil(t, err, "error should not be nil")
	updateSwap = &store.Swap{}

	de.DB.Where("other_chain_swap_id = ?", swap1.OtherChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)
}

func TestDeputy_sendBEP2HTLT_SwapExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	deputyAddressStr := "other_deputy_address"
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "swap_exist",
		OtherChainSwapId: "other_swap_id_swap_exist",
		BnbChainSwapId:   "bnb_swap_id_swap_exist",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	bnbChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, errors.New("any error"))

	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "swap error")

	bnbChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(true, nil)

	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "swap already exists")
}

func TestDeputy_sendBEP2HTLT_WrongExpireHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	deputyAddressStr := "other_deputy_address"
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "wrong_expire_height",
		OtherChainSwapId: "other_swap_id_wrong_expire_height",
		BnbChainSwapId:   "bnb_swap_id_wrong_expire_height",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	bnbChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	otherChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.OtherChainMinRemainHeight+1, nil)
	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")

	updateSwap := &store.Swap{}

	de.DB.Where("bnb_chain_swap_id = ?", swap.BnbChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)
}

func TestDeputy_sendBEP2HTLT_HTLTFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	deputyAddressStr := "other_deputy_address"
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "send_htlt_failed",
		OtherChainSwapId: "other_swap_id_send_htlt_failed",
		BnbChainSwapId:   "bnb_swap_id_send_htlt_failed",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)).String(),
		SenderAddr:       "sender_address",
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	bnbChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	otherChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.OtherChainMinRemainHeight-1, nil)
	otherChainExecutor.EXPECT().GetSwap(gomock.Any()).AnyTimes().Return(&common.SwapRequest{
		SenderAddress:    swap.SenderAddr,
		ExpireHeight:     swap.ExpireHeight,
		RecipientAddress: swap.ReceiverAddr,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)),
	}, nil)
	bnbChainExecutor.EXPECT().HTLT(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("",
		common.NewError(errors.New("Invalid sequence"), true))

	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "send bep2 HTLT tx error")

	txSent := &store.TxSent{}
	db.Where("swap_id = ?", swap.BnbChainSwapId).First(txSent)
	require.Equal(t, txSent.RandomNumberHash, "")

	bnbChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	otherChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.OtherChainMinRemainHeight-1, nil)
	bnbChainExecutor.EXPECT().HTLT(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("",
		common.NewError(errors.New("one error"), false))

	// other error, will create tx sent record
	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "send bep2 HTLT tx error")

	txSent = &store.TxSent{}
	db.Where("swap_id = ?", swap.BnbChainSwapId).First(txSent)
	require.EqualValues(t, txSent.Status, store.TxSentStatusFailed)
}

func TestDeputy_sendBEP2HTLT_HTLTFailed_MismatchedParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	deputyAddressStr := "other_deputy_address"
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "send_htlt_failed",
		OtherChainSwapId: "other_swap_id_send_htlt_failed",
		BnbChainSwapId:   "bnb_swap_id_send_htlt_failed",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)).String(),
		SenderAddr:       "sender_address",
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	bnbChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	otherChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.OtherChainMinRemainHeight-1, nil)
	otherChainExecutor.EXPECT().GetSwap(gomock.Any()).AnyTimes().Return(&common.SwapRequest{
		SenderAddress:    "wrong sender",
		ExpireHeight:     swap.ExpireHeight,
		RecipientAddress: swap.ReceiverAddr,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)),
	}, nil)

	_, err = de.sendBEP2HTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "reject swap for mismatch of parameters")
}

func TestDeputy_sendBEP2HTLT_HTLTSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	deputyAddressStr := "other_deputy_address"
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("OTHER_CHAIN")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "send_htlt_success",
		OtherChainSwapId: "other_swap_id_send_htlt_success",
		BnbChainSwapId:   "bnb_swap_id_send_htlt_success",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.OtherChainMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     deputyAddressStr,
		SenderAddr:       "sender_address",
	}
	de.DB.Create(swap)

	bnbChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	otherChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.OtherChainMinRemainHeight-1, nil)
	otherChainExecutor.EXPECT().GetSwap(gomock.Any()).AnyTimes().Return(&common.SwapRequest{
		SenderAddress:    swap.SenderAddr,
		RecipientAddress: swap.ReceiverAddr,
		ExpireHeight:     swap.ExpireHeight,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.OtherChainFixedFee.Int64()), big.NewInt(100)),
	}, nil)
	bnbChainExecutor.EXPECT().HTLT(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("tx_hash", nil)
	_, err = de.sendBEP2HTLT(swap)
	require.Nil(t, err, "error should be nil")

	updateSwap := &store.Swap{}

	de.DB.Where("bnb_chain_swap_id = ?", swap.BnbChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusBEP2HTLTSent)
}

func TestDeputy_sendOtherHTLT_Reject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	deputyAddressStr := "bnb_deputy_address"

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)

	de := NewDeputy(db, config, bnbChainExecutor, nil)

	swap := &store.Swap{
		RandomNumberHash: "reject_expire_height",
		OtherChainSwapId: "other_swap_id_reject_expire_height",
		BnbChainSwapId:   "bnb_swap_id_reject_expire_height",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.BnbMinAcceptExpireHeightSpan - 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbMaxSwapAmount.Int64()), big.NewInt(-1)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	_, err = de.sendOtherHTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	updateSwap := &store.Swap{}

	de.DB.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)

	swap1 := &store.Swap{
		RandomNumberHash: "reject_swap_amount",
		OtherChainSwapId: "other_swap_id_reject_swap_amount",
		BnbChainSwapId:   "bnb_swap_id_reject_swap_amount",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbMaxSwapAmount.Int64()), big.NewInt(10)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap1)

	_, err = de.sendOtherHTLT(swap1)
	require.NotNil(t, err, "error should not be nil")
	updateSwap = &store.Swap{}

	de.DB.Where("other_chain_swap_id = ?", swap1.OtherChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)
}

func TestDeputy_sendOtherHTLT_WrongOutAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	deputyAddressStr := "bnb_deputy_address"

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)

	de := NewDeputy(db, config, bnbChainExecutor, nil)

	swap := &store.Swap{
		RandomNumberHash: "reject_wrong_amount",
		OtherChainSwapId: "other_swap_id_reject_wrong_amount",
		BnbChainSwapId:   "bnb_swap_id_reject_wrong_amount",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(-1)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	_, err = de.sendOtherHTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	updateSwap := &store.Swap{}

	de.DB.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)

	swap1 := &store.Swap{
		RandomNumberHash: "reject_wrong_amount",
		OtherChainSwapId: "other_swap_id_reject_max_out_amount",
		BnbChainSwapId:   "bnb_swap_id_reject_max_out_amount",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbMaxSwapAmount.Int64()), big.NewInt(-1)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap1)

	config.ChainConfig.OtherChainMaxDeputyOutAmount = config.ChainConfig.OtherChainFixedFee

	_, err = de.sendOtherHTLT(swap1)
	require.NotNil(t, err, "error should not be nil")
	updateSwap = &store.Swap{}

	de.DB.Where("other_chain_swap_id = ?", swap1.OtherChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)
}

func TestDeputy_sendOtherHTLT_SwapExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	deputyAddressStr := "bnb_deputy_address"

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return(deputyAddressStr)

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("other_chain")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "swap_exist",
		OtherChainSwapId: "other_swap_id_swap_exist",
		BnbChainSwapId:   "bnb_swap_id_swap_exist",
		Height:           10,
		ExpireHeight:     10 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     deputyAddressStr,
	}
	de.DB.Create(swap)

	otherChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, errors.New("any error"))

	_, err = de.sendOtherHTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "swap error")

	otherChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(true, nil)

	_, err = de.sendOtherHTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "swap already exists")
}

func TestDeputy_sendOtherHTLT_WrongExpireHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return("bnb_deputy_address")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("mock_chain")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "wrong_expire_height",
		OtherChainSwapId: "other_swap_id_wrong_expire_height",
		BnbChainSwapId:   "bnb_swap_id_wrong_expire_height",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     config.BnbConfig.DeputyAddr.String(),
	}
	de.DB.Create(swap)

	otherChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	bnbChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.BnbMinRemainHeight+1, nil)

	_, err = de.sendOtherHTLT(swap)
	require.NotNil(t, err, "error should not be nil")

	updateSwap := &store.Swap{}

	de.DB.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusRejected)
}

func TestDeputy_sendOtherHTLT_HTLTFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return("bnb_deputy_address")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("mock_chain")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "send_htl_failed",
		OtherChainSwapId: "other_swap_id_send_htl_failed",
		BnbChainSwapId:   "bnb_swap_id_send_htl_failed",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     config.BnbConfig.DeputyAddr.String(),
		SenderAddr:       "sender_address",
	}
	de.DB.Create(swap)

	otherChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	bnbChainExecutor.EXPECT().GetSwap(gomock.Any()).AnyTimes().Return(&common.SwapRequest{
		SenderAddress:    swap.SenderAddr,
		RecipientAddress: swap.ReceiverAddr,
		ExpireHeight:     swap.ExpireHeight,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)),
	}, nil)
	bnbChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.BnbMinRemainHeight-1, nil)
	otherChainExecutor.EXPECT().HTLT(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("",
		common.NewError(errors.New("one error"), false))

	// Invalid sequence error
	_, err = de.sendOtherHTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "HTLT tx error")
}

func TestDeputy_sendOtherHTLT_HTLTFailed_MismatchedParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return("bnb_deputy_address")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("mock_chain")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "send_htl_failed",
		OtherChainSwapId: "other_swap_id_send_htl_failed",
		BnbChainSwapId:   "bnb_swap_id_send_htl_failed",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     config.BnbConfig.DeputyAddr.String(),
		SenderAddr:       "sender_address",
	}
	de.DB.Create(swap)

	otherChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	bnbChainExecutor.EXPECT().GetSwap(gomock.Any()).AnyTimes().Return(&common.SwapRequest{
		SenderAddress:    "wrong_sender",
		RecipientAddress: swap.ReceiverAddr,
		ExpireHeight:     swap.ExpireHeight,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)),
	}, nil)
	bnbChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.BnbMinRemainHeight-1, nil)

	// Invalid sequence error
	_, err = de.sendOtherHTLT(swap)
	require.NotNil(t, err, "error should not be nil")
	require.Contains(t, err.Error(), "reject swap for mismatch of parameters")
}

func TestDeputy_sendOtherHTLT_HTLTSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetDeputyAddress().AnyTimes().Return("bnb_deputy_address")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("mock_chain")

	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	swap := &store.Swap{
		RandomNumberHash: "send_htlt_success",
		OtherChainSwapId: "other_swap_id_send_htlt_success",
		BnbChainSwapId:   "bnb_swap_id_send_htlt_success",
		Height:           1000,
		ExpireHeight:     1000 + config.ChainConfig.BnbMinAcceptExpireHeightSpan + 1,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)).String(),
		ReceiverAddr:     config.BnbConfig.DeputyAddr.String(),
		SenderAddr:       "sender_address",
	}
	de.DB.Create(swap)

	otherChainExecutor.EXPECT().HasSwap(gomock.Any()).Return(false, nil)
	bnbChainExecutor.EXPECT().GetHeight().Return(swap.ExpireHeight-config.ChainConfig.BnbMinRemainHeight-1, nil)
	bnbChainExecutor.EXPECT().GetSwap(gomock.Any()).AnyTimes().Return(&common.SwapRequest{
		SenderAddress:    swap.SenderAddr,
		RecipientAddress: swap.ReceiverAddr,
		ExpireHeight:     swap.ExpireHeight,
		OutAmount:        new(big.Int).Add(big.NewInt(config.ChainConfig.BnbFixedFee.Int64()), big.NewInt(100)),
	}, nil)
	otherChainExecutor.EXPECT().HTLT(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("tx_hash", nil)
	_, err = de.sendOtherHTLT(swap)
	require.Nil(t, err, "error should be nil")

	updateSwap := &store.Swap{}

	de.DB.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updateSwap)
	require.Equal(t, updateSwap.Status, store.SwapStatusOtherHTLTSent)
}

func TestDeputy_sendOtherRefund(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	config := util.GetTestConfig()
	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	de := NewDeputy(nil, config, nil, otherChainExecutor)

	swap := &store.Swap{
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
	}

	otherChainExecutor.EXPECT().Refundable(gomock.Any()).Return(false, nil)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("mock_chain")

	_, err := de.sendOtherRefund(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), fmt.Sprintf("chain %s swap is not refundable", otherChainExecutor.GetChain()))

	otherChainExecutor.EXPECT().Refundable(gomock.Any()).Return(false, errors.New("any error"))
	_, err = de.sendOtherRefund(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), fmt.Sprintf("query chain %s swap error", otherChainExecutor.GetChain()))

	otherChainExecutor.EXPECT().Refundable(gomock.Any()).Return(true, nil)
	otherChainExecutor.EXPECT().Refund(gomock.Any()).Return("", common.NewError(errors.New("refund error"), true))

	_, err = de.sendOtherRefund(swap)
	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), fmt.Sprintf("send chain %s refund tx error", otherChainExecutor.GetChain()))

	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	db.Create(swap)
	de.DB = db

	otherChainExecutor.EXPECT().Refundable(gomock.Any()).Return(true, nil)
	otherChainExecutor.EXPECT().Refund(gomock.Any()).Return("txHash", nil)

	_, err = de.sendOtherRefund(swap)
	require.Nil(t, err, "refund can not fail")
	newSwap := &store.Swap{}
	db.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(newSwap)
	require.Equal(t, newSwap.Status, store.SwapStatusOtherRefundSent)
}

func TestDeputy_sendBEP2Refund(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	config := util.GetTestConfig()
	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	de := NewDeputy(nil, config, bnbChainExecutor, nil)

	swap := &store.Swap{
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
	}

	bnbChainExecutor.EXPECT().Refundable(gomock.Any()).Return(false, nil)
	bnbChainExecutor.EXPECT().GetChain().Return("mock_chain").AnyTimes()

	_, err := de.sendBEP2Refund(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "bep2 swap can not be refund")

	bnbChainExecutor.EXPECT().Refundable(gomock.Any()).Return(false, errors.New("any error"))
	_, err = de.sendBEP2Refund(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "query bep2 swap error")

	bnbChainExecutor.EXPECT().Refundable(gomock.Any()).Return(true, nil)
	bnbChainExecutor.EXPECT().Refund(gomock.Any()).Return("", common.NewError(errors.New("Invalid sequence"), true))

	_, err = de.sendBEP2Refund(swap)
	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "send bep2 refund tx error")

	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	db.Create(swap)
	de.DB = db

	bnbChainExecutor.EXPECT().Refundable(gomock.Any()).Return(true, nil)
	bnbChainExecutor.EXPECT().Refund(gomock.Any()).Return("",
		common.NewError(errors.New("other error"), false))

	_, err = de.sendBEP2Refund(swap)
	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "send bep2 refund tx error")

	txSent := &store.TxSent{}
	db.Where("swap_id = ?", swap.BnbChainSwapId).First(txSent)
	require.EqualValues(t, txSent.Status, store.TxSentStatusFailed)

	bnbChainExecutor.EXPECT().Refundable(gomock.Any()).Return(true, nil)
	bnbChainExecutor.EXPECT().Refund(gomock.Any()).Return("txHash", nil)

	_, err = de.sendBEP2Refund(swap)
	require.Nil(t, err, "refund can not fail")
	newSwap := &store.Swap{}
	db.Where("bnb_chain_swap_id = ?", swap.BnbChainSwapId).First(newSwap)
	require.Equal(t, newSwap.Status, store.SwapStatusBEP2RefundSent)
}

func TestDeputy_sendBEP2Claim(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	config := util.GetTestConfig()
	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	de := NewDeputy(db, config, bnbChainExecutor, nil)

	swap := &store.Swap{
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		RandomNumber:     "3a900a02d7cd96fe3bc3d23c562c59161a03613a04ab8c31bc26e7cbce4d0549",
	}

	bnbChainExecutor.EXPECT().Claimable(gomock.Any()).Return(false, nil)
	bnbChainExecutor.EXPECT().GetChain().Return("mock_chain").AnyTimes()

	_, err = de.sendBEP2Claim(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "bep2 swap can not be claimed")

	bnbChainExecutor.EXPECT().Claimable(gomock.Any()).Return(false, errors.New("any error"))
	_, err = de.sendBEP2Claim(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "query bep2 swap error")

	bnbChainExecutor.EXPECT().Claimable(gomock.Any()).Return(true, nil)
	bnbChainExecutor.EXPECT().Claim(gomock.Any(), gomock.Any()).Return("",
		common.NewError(errors.New("Invalid sequence"), true))

	_, err = de.sendBEP2Claim(swap)
	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "send bep2 claim tx error")

	de.DB.Create(swap)

	bnbChainExecutor.EXPECT().Claimable(gomock.Any()).Return(true, nil)
	bnbChainExecutor.EXPECT().Claim(gomock.Any(), gomock.Any()).Return("",
		common.NewError(errors.New("other error"), false))

	_, err = de.sendBEP2Claim(swap)
	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "send bep2 claim tx error")

	txSent := &store.TxSent{}
	db.Where("swap_id = ?", swap.BnbChainSwapId).First(txSent)
	require.EqualValues(t, txSent.Status, store.TxSentStatusFailed)

	bnbChainExecutor.EXPECT().Claimable(gomock.Any()).Return(true, nil)
	bnbChainExecutor.EXPECT().Claim(gomock.Any(), gomock.Any()).Return("txHash", nil)

	_, err = de.sendBEP2Claim(swap)
	require.Nil(t, err, "refund can not fail")
	newSwap := &store.Swap{}
	db.Where("bnb_chain_swap_id = ?", swap.BnbChainSwapId).First(newSwap)
	require.Equal(t, newSwap.Status, store.SwapStatusBEP2ClaimSent)
}

func TestDeputy_sendOtherClaim(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	config := util.GetTestConfig()
	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""

	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	de := NewDeputy(db, config, nil, otherChainExecutor)

	swap := &store.Swap{
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		RandomNumber:     "3a900a02d7cd96fe3bc3d23c562c59161a03613a04ab8c31bc26e7cbce4d0549",
	}

	otherChainExecutor.EXPECT().Claimable(gomock.Any()).Return(false, nil)
	otherChainExecutor.EXPECT().GetChain().Return("mock_chain").AnyTimes()
	otherChainExecutor.EXPECT().GetChain().Return("mock_chain").AnyTimes()

	_, err = de.sendOtherClaim(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), "swap is not claimable")

	otherChainExecutor.EXPECT().Claimable(gomock.Any()).Return(false, errors.New("any error"))
	_, err = de.sendOtherClaim(swap)

	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), fmt.Sprintf("query chain %s swap error", otherChainExecutor.GetChain()))

	otherChainExecutor.EXPECT().Claimable(gomock.Any()).Return(true, nil)
	otherChainExecutor.EXPECT().Claim(gomock.Any(), gomock.Any()).Return("",
		common.NewError(errors.New("other error"), true))

	_, err = de.sendOtherClaim(swap)
	require.NotNil(t, err, "err should not be nil")
	require.Contains(t, err.Error(), fmt.Sprintf("send chain %s claim tx error", otherChainExecutor.GetChain()))

	de.DB.Create(swap)

	otherChainExecutor.EXPECT().Claimable(gomock.Any()).Return(true, nil)
	otherChainExecutor.EXPECT().Claim(gomock.Any(), gomock.Any()).Return("txHash", nil)

	_, err = de.sendOtherClaim(swap)
	require.Nil(t, err, "refund can not fail")
	newSwap := &store.Swap{}
	db.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(newSwap)
	require.Equal(t, newSwap.Status, store.SwapStatusOtherClaimSent)
}

func TestDeputy_handleTxSent_NoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("other_chain")

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	swap := &store.Swap{
		Type:             store.SwapTypeBEP2ToOther,
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		Status:           store.SwapStatusOtherHTLTSent,
	}

	db.Create(swap)

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""
	de := NewDeputy(db, config, nil, otherChainExecutor)

	de.handleTxSent(swap, "mock_chain", store.TxTypeOtherHTLT, store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSentFailed)

	updatedSwap := &store.Swap{}
	db.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updatedSwap)
	require.Equal(t, updatedSwap.Status, store.SwapStatusBEP2HTLTConfirmed)
}

func TestDeputy_handleTxSent_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("other_chain")

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	swap := &store.Swap{
		Type:             store.SwapTypeBEP2ToOther,
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		Status:           store.SwapStatusOtherHTLTSent,
	}

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""
	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	db.Create(swap)

	txSent := &store.TxSent{
		Chain:            "other_chain",
		SwapId:           swap.OtherChainSwapId,
		Type:             store.TxTypeOtherHTLT,
		TxHash:           "somehash",
		RandomNumberHash: swap.RandomNumberHash,
		Status:           store.TxSentStatusPending,
		CreateTime:       time.Now().Add(-time.Duration(de.Config.ChainConfig.OtherChainAutoRetryTimeout+10) * time.Second).Unix(),
	}
	db.Create(txSent)

	de.handleTxSent(swap, "other_chain", store.TxTypeOtherHTLT, store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSentFailed)

	updatedSwap := &store.Swap{}
	db.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updatedSwap)
	require.Equal(t, updatedSwap.Status, store.SwapStatusBEP2HTLTConfirmed)
}

func TestDeputy_handleTxSent_TimeoutAndExceedRetryNum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("other_chain")

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	swap := &store.Swap{
		Type:             store.SwapTypeBEP2ToOther,
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		Status:           store.SwapStatusOtherHTLTSent,
	}

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""
	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	db.Create(swap)

	txSent := &store.TxSent{
		Chain:            "other_chain",
		Type:             store.TxTypeOtherHTLT,
		SwapId:           swap.OtherChainSwapId,
		TxHash:           "somehash",
		RandomNumberHash: swap.RandomNumberHash,
		Status:           store.TxSentStatusPending,
		CreateTime:       time.Now().Add(-time.Duration(de.Config.ChainConfig.OtherChainAutoRetryTimeout+10) * time.Second).Unix(),
	}
	txSent1 := &store.TxSent{
		Chain:            "other_chain",
		Type:             store.TxTypeOtherHTLT,
		SwapId:           swap.OtherChainSwapId,
		TxHash:           "somehash1",
		RandomNumberHash: swap.RandomNumberHash,
		Status:           store.TxSentStatusPending,
		CreateTime:       time.Now().Add(-time.Duration(de.Config.ChainConfig.OtherChainAutoRetryTimeout+10) * time.Second).Unix(),
	}
	txSent2 := &store.TxSent{
		Chain:            "other_chain",
		Type:             store.TxTypeOtherHTLT,
		SwapId:           swap.OtherChainSwapId,
		TxHash:           "somehash2",
		RandomNumberHash: swap.RandomNumberHash,
		Status:           store.TxSentStatusPending,
		CreateTime:       time.Now().Add(-time.Duration(de.Config.ChainConfig.OtherChainAutoRetryTimeout+10) * time.Second).Unix(),
	}
	db.Create(txSent)
	db.Create(txSent1)
	db.Create(txSent2)

	de.handleTxSent(swap, "other_chain", store.TxTypeOtherHTLT, store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSentFailed)

	updatedSwap := &store.Swap{}
	db.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updatedSwap)
	require.Equal(t, updatedSwap.Status, store.SwapStatusOtherHTLTSentFailed)
}

func TestDeputy_handleTxSent_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("other_chain")

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	swap := &store.Swap{
		Type:             store.SwapTypeBEP2ToOther,
		OtherChainSwapId: "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		BnbChainSwapId:   "5fb486af1f27d12d87d12757fb6710be10b56b287963bfe9ae9090d034952e8f",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		Status:           store.SwapStatusOtherHTLTSent,
	}

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""
	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	db.Create(swap)

	txSent := &store.TxSent{
		Chain:            "other_chain",
		Type:             store.TxTypeOtherHTLT,
		SwapId:           swap.OtherChainSwapId,
		TxHash:           "somehash",
		RandomNumberHash: swap.RandomNumberHash,
		Status:           store.TxSentStatusFailed,
		CreateTime:       time.Now().Add(-time.Duration(de.Config.ChainConfig.OtherChainAutoRetryTimeout+10) * time.Second).Unix(),
	}
	db.Create(txSent)

	de.handleTxSent(swap, "other_chain", store.TxTypeOtherHTLT, store.SwapStatusBEP2HTLTConfirmed, store.SwapStatusOtherHTLTSentFailed)

	updatedSwap := &store.Swap{}
	db.Where("other_chain_swap_id = ?", swap.OtherChainSwapId).First(updatedSwap)
	require.Equal(t, updatedSwap.Status, store.SwapStatusOtherHTLTSentFailed)
}

func TestDeputy_CheckTxSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bnbChainExecutor := mock.NewMockExecutor(ctrl)
	bnbChainExecutor.EXPECT().GetChain().AnyTimes().Return("BNB")
	bnbChainExecutor.EXPECT().GetSentTxStatus(gomock.Any()).AnyTimes().Return(store.TxSentStatusSuccess)

	otherChainExecutor := mock.NewMockExecutor(ctrl)
	otherChainExecutor.EXPECT().GetChain().AnyTimes().Return("other_chain")
	otherChainExecutor.EXPECT().GetSentTxStatus(gomock.Any()).AnyTimes().Return(store.TxSentStatusSuccess)

	config := util.GetTestConfig()
	db, err := util.PrepareDB(config)
	require.Nil(t, err, "create db error")

	config.AlertConfig.TelegramChatId = ""
	config.AlertConfig.TelegramBotId = ""
	de := NewDeputy(db, config, bnbChainExecutor, otherChainExecutor)

	txSent := &store.TxSent{
		Chain:            "other_chain",
		SwapId:           "695487c5277b81ebb573000a16f3f2b8da052109d725abe1ffe1dd55fe4f965f",
		Type:             store.TxTypeOtherHTLT,
		TxHash:           "somehash",
		RandomNumberHash: "b7fa5a2ac10c3eed718a850bfa6a1b71723df1a36edda00554b97ad4aa340f77",
		Status:           store.TxSentStatusInit,
		CreateTime:       time.Now().Add(-time.Duration(de.Config.ChainConfig.OtherChainAutoRetryTimeout+10) * time.Second).Unix(),
	}
	db.Create(txSent)

	de.CheckTxSent()

	updatedTxSent := &store.TxSent{}
	db.Where("swap_id = ?", txSent.SwapId).First(updatedTxSent)
	require.EqualValues(t, updatedTxSent.Status, store.TxSentStatusSuccess)
}
