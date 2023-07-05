
For now, we only support atomic swaps between bnb chain and ethereum. Our goal is to support
bnb chain and any other chain based on blocks and supports atomic swap(evm based or other platform).

For overview of deputy, you can refer to [this doc](./deputy.md).

Actually, it would be easy to support other chain if atomic swap procedure of other chain is the
same as current swap procedure, you can refer to [BEP3](https://github.com/binance-chain/BEPs/blob/master/BEP3.md)
for more details.

There are three main components: blockchain executor, blockchain observer and deputy. blockchain executor is an 
interface which interacts with blockchains, it has all methods we need so far. blockchain observer is a generic 
component which store txs and block info of blockchain, but we only support block based blockchains like ethereum, 
bnb chain, eos and etc. deputy is the component responsible for managing life cycle of swaps.

But if you want to add support for other chain like eos, you need to do is to implement a block executor based on eos.

## Step 1: add config

For the chain you want to support, you may need dedicated configuration for this chain, like deputy address, the way 
you want to initiate your deputy account and etc. You can refer to `eth_config` part in config.json.

What you need to do is add your config in `util/config.go` and do not forget to validate your config options.

## Step 2: implement executor

Suppose you have added your specific config, you need to do now is implement an executor. You can refer to [executor doc](./executor.md) 
here.

You can refer to bnb chain executor or ethereum executor. The interface contains all the methods we need for 
observer and deputy components. 

## Step 3: init your executor when starting 

If you have implemented your executor, you can init your executor in `main.go`, and replace ethereum executor with it.

```go
	dp := deputy.NewDeputy(db, config, bnbExecutor, ethExecutor)
	dp.Start()

	ob := observer.NewObserver(db, config, bnbExecutor, ethExecutor)
	ob.Start()
```

