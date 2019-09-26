
## Overview

![](./assets/deputy.svg)

### blockchain executor

blockchain executor will interact with blockchains. It can get blocks and txs from blockchains and it 
can send swap related txs to blockchains.

### blockchain observer

blockchain observer will fetch blocks and txs from blockchains and will store recent block and swap related 
txs.

observer will also send alert messages to telegram bot if block is not fetched for a long time in blockchains.

### deputy

Deputy will create swap and change swap status when related txs confirmed. It will also send corresponding txs automatically.

Deputy will send alert messages to telegram bot if tx sent is failed or swap is failed.

### admin

Admin is responsible for query deputy status, like current synced blockchain height, current blockchain height, deputy balance 
and so on.

You can query failed swaps which need resend corresponding tx manually. And you can resend tx via admin.

You can also change deputy work mode like stopping send HTLT txs.

