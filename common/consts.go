package common

import "time"

const ObserverMaxBlockNumber = 10000
const ObserverPruneInterval = 10 * time.Minute
const ObserverAlertInterval = 5 * time.Minute

const DeputyConfirmTxInterval = 5 * time.Second
const DeputySendTxInterval = 5 * time.Second
const DeputyCheckTxSentInterval = 5 * time.Second
const DeputyExpireUserHTLTInterval = 10 * time.Second
const DeputyAlertInterval = 5 * time.Minute
const DeputyReconInterval = 1 * time.Hour
const DeputyMetricsInterval = 10 * time.Second
