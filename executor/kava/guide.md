# Set up
## Kava 
1. Clone kava [repo]( https://github.com/Kava-Labs/kava)
2. Checkout branch `develop`
3. Install the `kvd` and `kvcli` binaries with `make install`
4. Initialize the blockchain with chain-id `testing`
5. Create 2 accounts `deputy` and `user` that both have a bnb balance of 100000000000.
6. In the genesis file located in `/contrib/testnet-5000/genesis-examples/genesis-bep3.json`, enter `deputy` address in the `deputy_address` field under bep3 params.

## Bnbchain
### Option 1: use the existing cloud server
We've set up a cloud server containing `bnbchaind` with the expected accounts/mnemonics and IP address. It can be used for testing, contact system admin @karzak for ssh access.
```bash
ssh ubuntu@ec2-3-231-211-245.compute-1.amazonaws.com
```

### Option 2: manual set up
1. Set up _local_ bnbchain with the steps provided [here](https://docs.binance.org/fullnode.html) using flag `--chain-id Binance-Chain-Tigris`.
**Note**: do not join the bnbchain mainnet. You must set up a local chain with chain id "Binance-Chain-Tigris".
2. With `bnbchaind` and `bnbcli` installed, create 2 accounts `deputy` and `user` that both have a BNB balance of 10000000000000.

## bep3-deputy
1. Clone this repo
2. Checkout branch `kava-deputy`
3. A sample config file is located at `/config/config.json`. Sections `bnb_config` and `kava_config` must be modified to match your bnbchain and kava chain variables. Enter kava `deputy` address and mnemonic and bnbchain `deputy` address and mnemonic. The `rpc_addr` field may also need to be updated.
4. Build bep3-deputy with `make build`

# Start processes
## Kava
Start Kava
```bash
kvd unsafe-reset-all
kvd start
```

## Bnbchain
Start bnbchain
```bash
bnbchaind unsafe-reset-all
bnbchaind start
```

If you're using the cloud server, you'll need to fund `deputy` and `user1` with BNB
```bash
bnbcli send --amount 10000000000000:BNB --to bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7 --from bnb1j9j2yfzs3x4xkqt3fgmyjn7kug4np3r6786y7c --chain-id Binance-Chain-Tigris
# enter password: "password"

bnbcli send --amount 100000000000:BNB --to bnb1urfermcg92dwq36572cx4xg84wpk3lfpksr5g7 --from bnb1j9j2yfzs3x4xkqt3fgmyjn7kug4np3r6786y7c --chain-id Binance-Chain-Tigris
# enter password: "password"
```

## bep3-deputy
Remove database:
```bash
cd /build
rm -rf deputy.db # deputy.db database must be removed between each run
```

Start deputy from `/build` dir:
```bash
./deputy --bnb-network 1 --kava-network 0 --config-type local --config-path "../config/config.json"
```

# Transfer BNB from Bnbchaind to Kava

##bnbchaind

Create Cross Chain HTLT (from bnbchain _user_ -> to kava _user_)
```bash
bnbcli token HTLT \
  --recipient-addr ${BNB_DEPUTY_ADDRESS} \
  --amount 30000000:BNB \
  --expected-income 10000000:BNB \
  --height-span 10001 \
  --from ${BNB_USER_ADDRESS} \
  --cross-chain true \
  --recipient-other-chain ${KAVA_USER_ADDRESS} \
  --chain-id Binance-Chain-Tigris
```

The HTLT creation process generates a random number and timestamp and requires a password:
```bash
Random number: 90a62df3efb640ea361008f619adf878e88e72bb877a5aba62230c2c3bb2c94f
Timestamp: 1583429130
Random number hash: 36371b5a6793cf411c1e55bbb5b6104981d79cce402259569f18bd044fb07303
Password to sign with 'user':
# If you're using the cloud server enter password: "password"
```

Get the `swapID` from the successful HTLT creation tx result. Tx result should be similar to:
```bash
Committed at block 332 (tx hash: C3D31985F37BB892883C91582D9D1235CA802E87DABCBF064BAB91E8D057696F, response: {Code:0 Data:[134 192 154 45 198 112 54 2 217 65 59 208 153 156 212 105 166 72 239 14 19 136 176 188 147 245 107 229 212 214 157 109] Log:Msg 0: swapID: 86c09a2dc6703602d9413bd0999cd469a648ef0e1388b0bc93f56be5d4d69d6d Info: GasWanted:0 GasUsed:0 Events:[{Type: Attributes:[{Key:[115 101 110 100 101 114] Value:[98 110 98 49 117 114 102 101 114 109 99 103 57 50 100 119 113 51 54 53 55 50 99 120 52 120 103 56 52 119 112 107 51 108 102 112 107 115 114 53 103 55] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 105 112 105 101 110 116] Value:[98 110 98 49 119 120 101 112 108 121 119 55 120 56 97 97 104 121 57 51 119 57 54 121 104 119 109 55 120 99 113 51 107 101 52 102 56 103 101 57 51 117] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[97 99 116 105 111 110] Value:[72 84 76 84] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

Query the new swap on bnbchain by its ID
```bash
  bnbcli token query-swap \
  --swap-id ${BNBCHAIN_SWAP_ID} \
  --chain-id Binance-Chain-Tigris
```

The deputy process will create an atomic swap on kava with the same information after _n_ blocks. It will log the expected kava swap ID as `swap.OtherChainSwapId` (it will be different than the bnbchain swap ID):
```bash
swap.OtherChainSwapId: edaf7deaf96d0ea583fa8a1cf3b547089418ad1065b4bca6ed856fcf8aaa110e
```

You'll need the kava swap ID and random number from above in order to claim the swap on kava.

## Kava

Claim atomic swap
```bash
kvcli tx bep3 claim ${KAVA_SWAP_ID} ${RANDOM_NUMBER} --from user
# enter password
```

A successful claim should output a tx result log similar to:
```bash
logs:
- msgindex: 0
  success: true
  log: ""
  events:
  - type: claimAtomicSwap
    attributes:
    - key: claim_sender
      value: kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj
    - key: recipient
      value: kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj
    - key: atomic_swap_id
      value: 53ab009facd0a9b2e02971cd4228d969f2f4fc41712b6984975d42b655159b7e
    - key: random_number_hash
      value: 36371b5a6793cf411c1e55bbb5b6104981d79cce402259569f18bd044fb07303
    - key: random_number
      value: 90a62df3efb640ea361008f619adf878e88e72bb877a5aba62230c2c3bb2c94f
  - type: message
    attributes:
    - key: action
      value: claimAtomicSwap
    - key: sender
      value: kava1eyugkwc74zejgwdwl7mvm7pad4hzdnka4wmdmu
    - key: module
      value: bep3
    - key: sender
      value: kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj
  - type: transfer
    attributes:
    - key: recipient
      value: kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj
    - key: amount
      value: 29999000bnb
```

Funds have been transferred from the user's bnbchain address to the user's kava address. After _n_ blocks the deputy will relay the successful claim to bnbchain, closing the swap.

# Transfer bnb from Kava to Bnbchain

## Kava
Create a new Atomic Swap
```bash
kvcli tx bep3 create ${KAVA_DEPUTY_ADDRESS} ${BNBCHAIN_USER_ADDRESS} ${BNBCHAIN_DEPUTY_ADDRESS} now 1111111bnb 1111111bnb 360 true --from user
# Note: you may need to include `--chain-id testing`
```

The Atomic Swap creation process generates a random number and timestamp and requires a password
```bash
Random number: ead368db570c229960e1c3bed0707b484210cb6ec40e7ecdd87a9c476a74b8ee
Timestamp: 1583545652
Random number hash: e3c58f4611e1a64c69714c39d1d8fdcaa7814873e3500de46dcbc36cd3db43d7

confirm transaction before signing and broadcasting [y/N]: y
Password to sign with 'user':
# enter password
```

Successful Atomic Swap creation will log a `createAtomicSwap` event similar to:
```bash
events:
- type: createAtomicSwap
  attributes:
  - key: sender
    value: kava1xy7hrjy9r0algz9w3gzm8u6mrpq97kwta747gj
  - key: recipient
    value: kava15qdefkmwswysgg4qxgqpqr35k3m49pkx2jdfnw
  - key: atomic_swap_id
    value: 084fbf068810a99340536580db51e2a2777b3a1b25af69eebfc9da02971c8e7c
  - key: random_number_hash
    value: 02a3987b7c6a1ff6ee035b47f6592392469eb9a492b7a99056245c2a6e33cdb7
  - key: timestamp
    value: "1583866396"
  - key: sender_other_chain
    value: bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7
  - key: expire_height
    value: "367"
  - key: amount
    value: 1111111bnb
  - key: expected_income
    value: 1111111bnb
```

Query the new swap on kava by its ID
```bash
  kvcli q bep3 swap ${KAVA_SWAP_ID}
```

After n blocks, the transaction will be relayed by the bep3-deputy to bnbchain. The bep3-deputy will log an INFO message similar to:
```bash
INFO sendBEP2HTLT send bep2 HTLT tx success, bnb_swap_id=448a4cc0e1d2b4bce793919fcb2e557aae44d96bc715af8e1a110f774747667d, tx_hash=0DAD58181C6537B05394F46AA42FD9C73002A6E8601205FED6FBAEEDDAE7E1D1
```

The swap on bnbchain can now be claimed
```bash
bnbcli token claim \
  --swap-id ${BNB_SWAP_ID} \
  --random-number ${RANDOM_NUMBER} \
  --from ${BNB_USER_ADDRESS} \
  --chain-id Binance-Chain-Tigris
```

Funds have been transferred from the user's kava address to the user's bnbchain address. After _n_ blocks the deputy will relay the successful claim to kava, closing the swap.

