package common

import "time"

const (
	ObserverMaxBlockNumber = 10000
	ObserverPruneInterval  = 10 * time.Minute
	ObserverAlertInterval  = 5 * time.Minute

	DeputyConfirmTxInterval      = 5 * time.Second
	DeputySendTxInterval         = 5 * time.Second
	DeputyCheckTxSentInterval    = 5 * time.Second
	DeputyExpireUserHTLTInterval = 10 * time.Second
	DeputyAlertInterval          = 5 * time.Minute
	DeputyReconInterval          = 1 * time.Hour
	DeputyMetricsInterval        = 10 * time.Second
	// Overflow interval must be long compared to the time it takes for a tx to be confirmed.
	// Otherwise it funds could be sent from the deputy while a previous tx is still processing, resulting in too much being sent out of the account.
	DeputyRunOverflowInterval = 15 * time.Second // 10 * time.Minute
)
