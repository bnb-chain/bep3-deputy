
## Overview

Executor interacts with blockchain and it provides methods that deputy needs. It can get blocks and txs from blockchain, 
get swap status, send swap related txs to blockchain and so on.

Executor is a interface, if you want to add support for other chain besides ethereum and binance chain, you just need
to implement an executor.

Below is the interface detail:
```go
	// GetChain returns unique name of the chain(like BNB, ETH and etc)
	GetChain() string
	// GetHeight returns current height of chain
	GetHeight() (int64, error)
	// GetBlockAndTxs returns block info and txs included in this block
	GetBlockAndTxs(height int64) (*BlockAndTxLogs, error)
	// GetFetchInterval returns fetch interval of the chain like average blocking time, it is used in observer
	GetFetchInterval() time.Duration
	// GetDeputyAddress returns deputy account address
	GetDeputyAddress() string
	// GetSentTxStatus returns status of tx sent
	GetSentTxStatus(hash string) store.TxStatus
	// GetBalance returns balance of swap token
	GetBalance() (*big.Int, error)
	// GetStatus returns status of deputy account, like balance of deputy account
	GetStatus() (interface{}, error)
	// GetBalanceAlertMsg returns balance alert message if necessary, like account balance is less than amount in config
	GetBalanceAlertMsg() (string, error)
	// IsSameAddress returns is addrA the same with addrB
	IsSameAddress(addrA string, addrB string) bool
	// CalcSwapId calculate swap id for each chain
	CalcSwapId(randomNumberHash common.Hash, sender string, senderOtherChain string) ([]byte, error)
	// Claimable returns is swap claimable
	Claimable(swapId common.Hash) (bool, error)
	// Refundable returns is swap refundable
	Refundable(swapId common.Hash) (bool, error)
	// GetSwap returns swap request detail
	GetSwap(swapId common.Hash) (*SwapRequest, error)
	// HasSwap returns does swap exist
	HasSwap(swapId common.Hash) (bool, error)
	// HTLT sends htlt tx
	HTLT(randomNumberHash common.Hash, timestamp int64, heightSpan int64, recipientAddr string, otherChainSenderAddr string,
		otherChainRecipientAddr string, outAmount *big.Int) (string, *Error)
	// Claim sends claim tx
	Claim(swapId common.Hash, randomNumber common.Hash) (string, *Error)
	// Refund sends refund tx
	Refund(swapId common.Hash) (string, *Error)
```

For implement detail, you can refer to executors implemented now.

