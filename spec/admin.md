Admin is a server for management purpose, it provides some endpoints like getting deputy status, resend tx for
failed swaps and etc.

## Endpoints

For example, if you start you admin server at `127.0.0.1:8000`, you can get endpoints provided
by admin server:

```bash
$ curl 127.0.0.1:8000
{
    "endpoints": [
        "/status",
        "/failed_swaps/{page}",
        "/resend_tx/{id}",
        "/set_mode/{mode}"
    ]
}
```

### Query deputy status

You can query deputy status via `127.0.0.1:8000/status`.

```bash
$ curl 127.0.0.1:8000/status
{
    "mode": "NormalMode",
    "bnb_chain_height": 74271,
    "bnb_sync_height": 74270,
    "other_chain_height": 6247835,
    "other_chain_sync_height": 6247835,
    "bnb_chain_last_block_fetched_at": "2019-08-23T16:21:15+08:00",
    "other_chain_last_block_fetched_at": "2019-08-23T16:21:12+08:00",
    "bnb_status": {
        "balance": [
            {
                "symbol": "BNB",
                "free": "197993354.10375000",
                "locked": "10.00000000",
                "frozen": "0.00000000"
            }
        ]
    },
    "other_chain_status": {
        "allowance": "768.0023",
        "erc20_balance": "48557.41037",
        "eth_balance": "9.342522577"
    }
}
```


### Query failed swaps

You can query failed swaps(mostly tx sent by deputy failed for some reason and need resend corresponding tx manually)
through `127.0.0.1:8000/failed_swaps/{page}`. You need to specify page and admin server will return 100 swaps per page.

```bash
$ curl 127.0.0.1:8000/failed_swaps/1
{
    "total_count": 1,
    "cur_page": 1,
    "page_num": 100,
    "swaps": [
        {   
            "id": 1,
            "type": "BEP2_TO_OTHER",
            "sender_addr": "bnb1t6haxvhczufp0g9jfafm44aw97amcsf29a22hh",
            "receiver_addr": "bnb1gqumpqlqmz8juyysrxac273j6fv56sztxmfr3e",
            "other_chain_addr": "042ccc750e1099068622bb521003f207297a40b0",
            "in_amount": "100000000",
            "out_amount": "100000000",
            "random_number_hash": "85b8c6bac7c7500b1bebd54c93937fe8fdcf0d6b91f5597d5b85427554462abc",
            "expire_height": 121719,
            "height": 21719,
            "timestamp": 1566475580,
            "random_number": "53c5891150dea1edb77ebfb83c353f85fc6f3b3e3ef57f00b43f8678fbbe0d22",
            "status": "BEP2_CLAIM_SENT_FAILED",
            "txs_sent": [
                {
                    "id": 329,
                    "chain": "BNB",
                    "type": "BEP2_CLAIM",
                    "tx_hash": "FEAB43A1998BD4ACCF839AF0751236ED437B01835269DAE2B66CACE56D3224D0",
                    "random_number_hash": "85b8c6bac7c7500b1bebd54c93937fe8fdcf0d6b91f5597d5b85427554462abc",
                    "err_msg": "",
                    "status": "FAILED",
                    "create_time": 1566478074,
                    "update_time": 1566478079
                },
                {
                    "id": 325,
                    "chain": "ETH",
                    "type": "OTHER_HTLT",
                    "tx_hash": "0x8f722899cd7411c39914d37bbfddab6b8515e30e604d1ac1d4b0daa036c403a4",
                    "random_number_hash": "85b8c6bac7c7500b1bebd54c93937fe8fdcf0d6b91f5597d5b85427554462abc",
                    "err_msg": "",
                    "status": "SUCCESS",
                    "create_time": 1566475597,
                    "update_time": 1566475613
                }
            ]
        }
    ]
}
```

Admin server will return detail of failed swaps and txs sent.

### Resend corresponding tx for failed swap

If you have figured out why the tx sent is failed, you can resend corresponding tx if needed.

```bash
$ curl 127.0.0.1:8000/resend_tx/1
```

### Change work mode of deputy

Sometimes, if we want to stop swap, we can change work mode of deputy to `1`(which means `StopSendHTLTMode`), deputy 
will stop send out HTLT txs. You can change the mode back to normal(0) if you are cool.

```bash
$ curl 127.0.0.1:8000/set_mode/1
```

### NOTICE

For endpoints of admin server does not have auth, so please **DO NOT** expose admin server to others.