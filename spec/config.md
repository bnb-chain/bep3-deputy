## DB config

DB config is config of database. Deputy will prune blocks and txs for we just need recent blocks and txs of blocks and 
storage will not grow without limit.

+ dialect: it should be `sqlite3` or `mysql`, only sqlite and mysql are supported for now.
+ db_path: db file path or mysql db config, eg(`root:12345678@(127.0.0.1:3306)/deputy?charset=utf8&parseTime=True&loc=Local`).
+ max_bnb_kept_block_height: how many recent blocks of binance chain you want to keep.
+ max_other_kept_block_height: how many recent blocks of other chain you want to keep.


## Alert config

Deputy will send alert messages to telegram group if block is not fetched for a long time or tx sent is failed.

+ telegram_bot_id: `telegram_bot_id` is your telegram bot id.
+ telegram_chat_id: `telegram_chat_id` is chat id of group your bot joined.
+ bnb_block_update_time_out: `bnb_block_update_time_out` is how long(in seconds) that block is not fetched in binance chain you want 
deputy to send alert messages.
+ other_chain_block_update_time_out: `other_chain_block_update_time_out` is how long(in seconds) that block is not fetched in other chain
you want deputy to send alert messages.
+ reconciliation_diff_amount: deputy will sum up all tokens of both chain and compare the result with last result. If diff between 
two result is larger than `reconciliation_diff_amount`, deputy will send alert msg to tg.

References:
+ [create a bot](https://core.telegram.org/bots#6-botfather)
+ [get bot id and chat id](https://stackoverflow.com/questions/32423837/telegram-bot-how-to-get-a-group-chat-id)

## Chain config

Chain common config for deputy. Pls note that `swap_amount` and `fixed_fee` below are number with decimal. For example, decimal in binance chain 
is 8 which means 100000000 is 1 actually. You need to handle decimal and amount with decimal.

+ bnb_confirm_num: number of confirmations in binance chain. 2 is enough for binance does not have forks.
+ bnb_auto_retry_num: number of retry if tx sent by deputy in binance chain is lost(node returned tx hash but it's not included in blockchain).
+ bnb_auto_retry_timeout: how long(in seconds) you think tx is lost if it is not included in binance chain. 
+ bnb_expire_height_span: expire height span of HTLT tx sent by deputy.
+ bnb_min_accept_expire_height_span: min expire height span for received HTLT tx in binance chain.
+ bnb_min_remain_height: min expire remaining height for HTLT tx in binance chain when send HTLT tx in other chain.
+ bnb_min_swap_amount: min swap amount(with decimal) for each swap request.
+ bnb_max_swap_amount: max swap amount(with decimal) for each swap request.
+ bnb_max_deputy_out_amount: max deputy out amount(with decimal) for each swap request.
+ bnb_ratio: ratio of token swap out from deputy in binance chain. for example, if ratio is 0.8, when someone send deputy 100 tokens, deputy
will swap 100*0.8=80 tokens out.
+ bnb_fixed_fee: fee(with decimal) deputy wants to charge for every swap request for swap related txs cost `BNB`. for example, if fixed fee is 100 and deputy 
need to swap out 1000 tokens, it will deduct fixed fee first, and it will swap 1000-100=900 tokens out at last.
+ bnb_start_height: start height of binance chain you want to sync like block chain height when you start your deputy.

+ other_chain: chain name of other chain, `ETH` is supported only for now.
+ other_chain_confirm_num: number of confirmations in other chain. 
+ other_chain_decimal: decimal of token you want to swap in other chain. for example, decimal of ETH if 18
+ other_chain_auto_retry_num: number of retry if tx sent by deputy in other chain is lost(node returned tx hash but it's not included in blockchain).
+ other_chain_auto_retry_timeout: how long(in seconds) you think tx is lost if it is not included in other chain(pending or can not found in blockchain). 
+ other_chain_expire_height_span: expire height span of HTLT tx sent by deputy.
+ other_chain_min_accept_expire_height_span: min expire height span for received HTLT tx in other chain.
+ other_chain_min_remain_height: min expire remaining height for HTLT tx in other chain when send HTLT tx in binance chain.
+ other_chain_min_swap_amount: min swap amount(with decimal) for each swap request.
+ other_chain_max_swap_amount: max swap amount(with decimal) for each swap request.
+ other_chain_max_deputy_out_amount: max deputy out amount(with decimal) for each swap request.
+ other_chain_ratio: ratio of token swap out from deputy in other chain. for example, if ratio is 0.8, when someone send deputy 100 tokens, deputy
will swap 100*0.8=80 tokens out.
+ other_chain_fixed_fee": fee(with decimal) deputy wants to charge for every swap request for swap related txs cost gas in blockchain. for example, 
if fixed fee is 100 and deputy need to swap out 1000 tokens, it will deduct fixed fee first, and it will swap 1000-100=900 tokens out at last.
+ other_chain_start_height: start height of other chain you want to sync like block chain height when you start your deputy.


## Admin config

+ listen_addr: listen address of admin server. If you want to deploy deputy in container, you may need to listen on `0.0.0.0`.

## Instrumentation config

+ prometheus: export prometheus metrics or not
+ prometheus_listen_addr: listen address of prometheus. If you want to deploy deputy in container, you may need to listen on `0.0.0.0`.

## Binance chain config

+ key_type:  `mnemonic` and `aws_mnemonic` supported. `mnemonic` will use mnemonic provided below and `aws_mnemonic`
 will fetch mnemonic from aws secret manager.
+ aws_region: region of aws
+ aws_secret_name: secret name of private key in aws
+ mnemonic: mnemonic of deputy account
+ rpc_addr": rpc address of binance chain
+ symbol: symbol of token you want to swap
+ deputy_addr: address of deputy
+ fetch_interval: block fetch interval.
+ token_balance_alert_threshold: balance(with decimal) of token you want to swap that you want deputy to send alert message to telegram.
+ bnb_balance_alert_threshold: balance(with decimal) of `BNB` that you want deputy to send alert message to telegram.

## Ethereum config

+ swap_type": `erc20_swap` for erc20 token swap and `eth_swap` for eth swap.
+ key_type": `private_key` and `aws_private_key` supported. `private_key` will use private key provided below and `aws_private_key`
will fetch private key from aws secret manager.
+ aws_region: region of aws
+ aws_secret_name: secret name of private key in aws
+ private_key": private key of deputy account
+ provider: provider address of ethereum
+ swap_contract_addr: swap contract address
+ token_contract_addr: erc20 token contract address
+ deputy_addr: deputy address
+ gas_limit: gas limit of ethereum
+ fetch_interval: block fetch interval of ethereum
+ token_balance_alert_threshold: balance(with decimal) of token you want to swap that you want deputy to send alert message to telegram.
+ eth_balance_alert_threshold: balance(with decimal) of `ETH` that you want deputy to send alert message to telegram.
+ allowance_balance_alert_threshold: balance(with decimal) of token approved to swap contract that you want deputy to send alert message to telegram.

## Log config

+ level: level of log, `CRITICAL`,`ERROR`,`WARNING`,`NOTICE`,`INFO`,`DEBUG` are supported.
+ filename: log file path if `use_console_logger` is true
+ max_file_size_in_mb: max log file size
+ max_backups_of_log_files: max backups of log files
+ max_age_to_retain_log_files_in_days: max days to retain log files
+ use_console_logger: use console logger or not
+ use_file_logger: use file logger or not
+ compress: compress log file or not
