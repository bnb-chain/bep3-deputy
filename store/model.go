package store

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BlockLog struct {
	Id         int64
	Chain      string
	BlockHash  string
	ParentHash string
	Height     int64
	BlockTime  int64
	CreateTime int64
}

func (BlockLog) TableName() string {
	return "block_log"
}
func (l *BlockLog) BeforeCreate() (err error) {
	l.CreateTime = time.Now().Unix()
	return nil
}

type TxLogStatus string

const (
	TxStatusInit      TxLogStatus = "INIT"
	TxStatusConfirmed TxLogStatus = "CONFIRMED"
)

type TxLog struct {
	Id    int64
	Chain string
	// swap id should be hex encoded bytes without '0x' prefix
	SwapId       string
	TxType       TxType
	TxHash       string
	ContractAddr string
	// sender address should be encoded by each executor
	SenderAddr string
	// receiver address should be encoded by each executor
	ReceiverAddr     string
	SenderOtherChain string
	OtherChainAddr   string
	InAmount         string
	OutAmount        string
	OutCoin          string
	// random number hash should be hex encoded bytes without '0x' prefix
	RandomNumberHash string
	ExpireHeight     int64
	Timestamp        int64
	RandomNumber     string
	BlockHash        string
	Height           int64
	Status           TxLogStatus
	ConfirmedNum     int64
	CreateTime       int64
	UpdateTime       int64
}

func (TxLog) TableName() string {
	return "tx_log"
}

func (t *TxLog) BeforeCreate() (err error) {
	t.CreateTime = time.Now().Unix()
	t.UpdateTime = time.Now().Unix()
	if t.Status == "" {
		t.Status = TxStatusInit
	}
	return nil
}

type SwapType string

const (
	SwapTypeOtherToBEP2 SwapType = "OTHER_TO_BEP2"
	SwapTypeBEP2ToOther SwapType = "BEP2_TO_OTHER"
)

type SwapStatus string

const (
	SwapStatusOtherHTLTConfirmed    SwapStatus = "OTHER_HTLT_CONFIRMED"
	SwapStatusOtherHTLTSent         SwapStatus = "OTHER_HTLT_SENT"
	SwapStatusOtherHTLTExpired      SwapStatus = "OTHER_HTLT_EXPIRED"
	SwapStatusOtherHTLTSentFailed   SwapStatus = "OTHER_HTLT_SENT_FAILED"
	SwapStatusOtherClaimSent        SwapStatus = "OTHER_CLAIM_SENT"
	SwapStatusOtherClaimConfirmed   SwapStatus = "OTHER_CLAIM_CONFIRMED"
	SwapStatusOtherClaimSentFailed  SwapStatus = "OTHER_CLAIM_SENT_FAILED"
	SwapStatusOtherRefundSent       SwapStatus = "OTHER_REFUND_SENT"
	SwapStatusOtherRefundConfirmed  SwapStatus = "OTHER_REFUND_CONFIRMED"
	SwapStatusOtherRefundSentFailed SwapStatus = "OTHER_REFUND_SENT_FAILED"

	SwapStatusBEP2HTLTSent         SwapStatus = "BEP2_HTLT_SENT"
	SwapStatusBEP2HTLTConfirmed    SwapStatus = "BEP2_HTLT_CONFIRMED"
	SwapStatusBEP2HTLTExpired      SwapStatus = "BEP2_HTLT_EXPIRED"
	SwapStatusBEP2HTLTSentFailed   SwapStatus = "BEP2_HTLT_SENT_FAILED"
	SwapStatusBEP2ClaimConfirmed   SwapStatus = "BEP2_CLAIM_CONFIRMED"
	SwapStatusBEP2ClaimSent        SwapStatus = "BEP2_CLAIM_SENT"
	SwapStatusBEP2ClaimSentFailed  SwapStatus = "BEP2_CLAIM_SENT_FAILED"
	SwapStatusBEP2RefundSent       SwapStatus = "BEP2_REFUND_SENT"
	SwapStatusBEP2RefundConfirmed  SwapStatus = "BEP2_REFUND_CONFIRMED"
	SwapStatusBEP2RefundSentFailed SwapStatus = "BEP2_REFUND_SENT_FAILED"

	SwapStatusRejected SwapStatus = "REJECTED"
)

type Swap struct {
	Id   int64
	Type SwapType
	// bnb chain swap id should be hex encoded bytes without '0x' prefix
	BnbChainSwapId string
	// other chain swap id should be hex encoded bytes without '0x' prefix
	OtherChainSwapId string
	SenderAddr       string
	ReceiverAddr     string
	OtherChainAddr   string
	InAmount         string
	OutAmount        string
	DeputyOutAmount  string
	RandomNumberHash string `gorm:"not null"`
	ExpireHeight     int64
	Height           int64
	Timestamp        int64
	RandomNumber     string
	Status           SwapStatus
	CreateTime       int64
	UpdateTime       int64
}

func (Swap) TableName() string {
	return "swap"
}

func (t *Swap) BeforeCreate() (err error) {
	t.CreateTime = time.Now().Unix()
	t.UpdateTime = time.Now().Unix()
	return nil
}

type TxStatus string

const (
	TxSentStatusInit     TxStatus = "INIT"
	TxSentStatusNotFound TxStatus = "NOT_FOUND"
	TxSentStatusPending  TxStatus = "PENDING"
	TxSentStatusFailed   TxStatus = "FAILED"
	TxSentStatusSuccess  TxStatus = "SUCCESS"
	TxSentStatusLost     TxStatus = "LOST"
)

type TxType string

const (
	TxTypeOtherHTLT   TxType = "OTHER_HTLT"
	TxTypeOtherClaim  TxType = "OTHER_CLAIM"
	TxTypeOtherRefund TxType = "OTHER_REFUND"
	TxTypeBEP2HTLT    TxType = "BEP2_HTLT"
	TxTypeBEP2Claim   TxType = "BEP2_CLAIM"
	TxTypeBEP2Refund  TxType = "BEP2_REFUND"
)

type TxSent struct {
	Id               int64    `json:"id"`
	Chain            string   `json:"chain"`
	SwapId           string   `json:"swap_id"`
	Type             TxType   `json:"type"`
	TxHash           string   `json:"tx_hash"`
	RandomNumberHash string   `json:"random_number_hash"`
	ErrMsg           string   `json:"err_msg"`
	Status           TxStatus `json:"status"`
	CreateTime       int64    `json:"create_time"`
	UpdateTime       int64    `json:"update_time"`
}

func (TxSent) TableName() string {
	return "tx_sent"
}

func (t *TxSent) BeforeCreate() (err error) {
	t.UpdateTime = time.Now().Unix()
	if t.CreateTime == 0 {
		t.CreateTime = t.UpdateTime
	}
	if t.Status == "" {
		t.Status = TxSentStatusInit
	}
	return nil
}

func InitTables(db *gorm.DB) {
	if !db.HasTable(&BlockLog{}) {
		db.CreateTable(&BlockLog{})
		db.Model(&BlockLog{}).AddUniqueIndex("idx_block_log_height", "chain", "height")
		db.Model(&BlockLog{}).AddIndex("idx_block_log_create_time", "create_time")
	}

	if !db.HasTable(&TxLog{}) {
		db.CreateTable(&TxLog{})
		db.Model(&TxLog{}).AddIndex("idx_tx_log_height", "chain", "height")
		db.Model(&TxLog{}).AddUniqueIndex("idx_tx_log_tx_hash", "tx_hash")
		db.Model(&TxLog{}).AddIndex("idx_tx_log_chain_status", "chain", "status")
		db.Model(&TxLog{}).AddIndex("idx_tx_log_swap_id", "swap_id")
		db.Model(&TxLog{}).AddIndex("idx_tx_log_create_time", "create_time")

	}

	if !db.HasTable(&Swap{}) {
		db.CreateTable(&Swap{})
		db.Model(&Swap{}).AddIndex("idx_swap_type_status", "type", "status")
		db.Model(&Swap{}).AddUniqueIndex("idx_swap_bnb_swap_id", "bnb_chain_swap_id")
		db.Model(&Swap{}).AddUniqueIndex("idx_swap_other_swap_id", "other_chain_swap_id")
		db.Model(&Swap{}).AddIndex("idx_swap_create_time", "create_time")
		db.Model(&Swap{}).AddIndex("idx_swap_update_time", "update_time")
	}

	if !db.HasTable(&TxSent{}) {
		db.CreateTable(&TxSent{})
		db.Model(&TxSent{}).AddUniqueIndex("idx_tx_sent_tx_hash", "tx_hash")
		db.Model(&TxSent{}).AddIndex("idx_tx_sent_swap_id", "swap_id")
		db.Model(&TxSent{}).AddIndex("idx_tx_sent_status", "status")
		db.Model(&TxSent{}).AddIndex("idx_tx_sent_create_time", "create_time")
	}
}
