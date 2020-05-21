# Set up

## kava

Clone kava [repo](https://github.com/Kava-Labs/kava)

Remove existing kava binaries

```bash
rm -rf ~/.kvd
rm -rf ~/.kvcli
```

Checkout develop branch

```bash
git checkout develop
```

Install the `kvd` and `kvcli` binaries

```bash
make install
```

Initialize the blockchain with chain-id `testing`

```bash
# moniker is your preferred nickname
kvd init --chain-id=testing ${MONIKER}
```

```bash
# Copy sample bep3 genesis file to config
cp ./contrib/testnet-5000/genesis_examples/genesis_bep3.json ~/.kvd/config/genesis.json
```

Add genesis accounts `deputy` and `user` with bnb balance.

```bash
kvcli keys add deputy
# enter a new password
# kvcli will print the deputy's mnemonic phrase, we'll need this later

kvd add-genesis-account $(kvcli keys show deputy -a) 1000000000000bnb

kvcli keys add user
# enter a new password
kvd add-genesis-account $(kvcli keys show user -a) 1000000000000bnb
```

Populate genesis file with custom values

```bash
brew install moreutils

# replace KAVA_DEPUTY_ADDRESS with the deputy address from above
jq '.app_state.bep3.params.bnb_deputy_address="KAVA_DEPUTY_ADDRESS"' ~/.kvd/config/genesis.json|sponge ~/.kvd/config/genesis.json
```

Create validator and collect gentxs

```bash
kvcli keys add validator
# enter a new password
kvd add-genesis-account $(kvcli keys show validator -a) 5000000000000ukava
kvd gentx --name validator --amount 100000000ukava
# enter password

kvd collect-gentxs
kvcli config trust-node true
kvcli config chain-id testing

# check genesis is valid
kvd validate-genesis
```

## bnbchain testnet

1. Download tbnbcli by following these [steps](https://docs.binance.org/fullnode.html).
2. Create two new testnet accounts `deputy` and `user` and load them with testnet BNB from the [faucet](https://www.binance.vision/tutorials/binance-dex-funding-your-testnet-account). Save the deputy's mnemonic.

## bep3-deputy

Clone this repo and checkout the `kava-deputy` branch

```bash
git clone git@github.com:Kava-Labs/bep3-deputy.git
git checkout kava-deputy
```

A sample config file is located at `/config/config_kava.json`. Modify the following sections:

- `chain_config`: enter latest testnet block height in `bnb_start_height`
- `bnb_config`: enter bnbchain deputy details in `deputy_addr` and `mnemonic`, check `rpc_addr`
- `kava_config`: enter kava deputy details in `deputy_addr` and `mnemonic`, check `rpc_addr`

Build bep3-deputy

```bash
make build
```

# Start processes

## kava

Start kava:

```bash
kvd unsafe-reset-all
kvd start
```

## bep3-deputy

Start deputy:

```bash
cd /build

# deputy.db database must be removed between each run
rm -rf deputy.db

# Start deputy from `/build` dir:
# --kava-network (0 for local chain ID 'testing', 1 for testnet chain ID 'kava-testnet-6000', 2 for mainnet 'kava-3')
./deputy --bnb-network 0 --kava-network 1 --config-type local --config-path "../config/test_config_kava.json"
```

# Transfer BNB from bnbchain to kava

## bnbchain testnet

Create cross-chain HTLT (from bnbchain _user_ -> to kava _user_)

```bash
tbnbcli token HTLT \
  --recipient-addr ${BNB_DEPUTY_ADDRESS} \
  --amount 10000000:BNB \
  --expected-income 10000000:BNB \
  --height-span 10001 \
  --from ${BNB_USER_ADDRESS} \
  --cross-chain true \
  --recipient-other-chain ${KAVA_USER_ADDRESS} \
  --chain-id Binance-Chain-Nile \
  --node ${BNB_RPC_URL} \
  --trust-node
```

The HTLT creation process generates a random number and timestamp and requires a password:

```bash
Random number: 90a62df3efb640ea361008f619adf878e88e72bb877a5aba62230c2c3bb2c94f
Timestamp: 1583429130
Random number hash: 36371b5a6793cf411c1e55bbb5b6104981d79cce402259569f18bd044fb07303
Password to sign with 'user':
# enter password
```

Get the `swapID` from the successful HTLT creation tx result. Tx result should be similar to:

```bash
Committed at block 332 (tx hash: C3D31985F37BB892883C91582D9D1235CA802E87DABCBF064BAB91E8D057696F, response: {Code:0 Data:[134 192 154 45 198 112 54 2 217 65 59 208 153 156 212 105 166 72 239 14 19 136 176 188 147 245 107 229 212 214 157 109] Log:Msg 0: swapID: 86c09a2dc6703602d9413bd0999cd469a648ef0e1388b0bc93f56be5d4d69d6d Info: GasWanted:0 GasUsed:0 Events:[{Type: Attributes:[{Key:[115 101 110 100 101 114] Value:[98 110 98 49 117 114 102 101 114 109 99 103 57 50 100 119 113 51 54 53 55 50 99 120 52 120 103 56 52 119 112 107 51 108 102 112 107 115 114 53 103 55] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 105 112 105 101 110 116] Value:[98 110 98 49 119 120 101 112 108 121 119 55 120 56 97 97 104 121 57 51 119 57 54 121 104 119 109 55 120 99 113 51 107 101 52 102 56 103 101 57 51 117] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[97 99 116 105 111 110] Value:[72 84 76 84] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

Query the new swap on bnbchain by its ID

```bash
tbnbcli token query-swap \
  --swap-id ${BNBCHAIN_SWAP_ID} \
  --chain-id Binance-Chain-Nile \
  --node ${BNB_RPC_URL} \
  --trust-node
```

The deputy process will create an atomic swap on kava with the same information after _n_ blocks

```bash
# The deputy will log the expected other_chain_swap_id
2020-03-26 15:09:59 INFO sendOtherHTLT send chain KAVA HTLT tx success, other_chain_swap_id=da89ae0c4f341ffa38345c635725bcc0d4e221b807fc7f143fabdd1e13c3b4d5, tx_hash=A3DC82B10373B30D00D30BE253DF34DADF0D57CCEAD319F09DDD3553ED2B36FC
```

The `other_chain_swap_id` and the generated random number from above can be used to claim the swap on kava.

## kava

Claim atomic swap

```bash
kvcli tx bep3 claim ${KAVA_SWAP_ID} ${RANDOM_NUMBER} --from user
# enter password
```

Funds will be transferred to the intended recipient's kava address. The claim should output a tx result log containing:

```bash
logs:
- msgindex: 0
  success: true
  log: ""
  events:
  - type: message
    attributes:
    - key: action
      value: claimAtomicSwap
```

After _n_ blocks the deputy will relay the successful claim to bnbchain, closing the swap:

```bash
# the deputy will log an INFO message similar to:
2020-03-26 15:17:17 INFO sendBEP2Claim send bep2 claim tx success, bnb_swap_id=4c4abc3fcc7a7e9b4f7d586f439c16312a20f49ab8129a00a86dd54257d79b6f, random_number=0x90e2cbb1a04a24553736adcee3c7862b536e7afaa34634d65f397812702030f2, tx_hash=37A4AB493F607041CE1367BC3FBDE382629D60808A07203C44A74A21D9D2D19A
```

# Transfer bnb from kava to bnbchain

## kava

Create a new Atomic Swap

```bash
kvcli tx bep3 create ${KAVA_DEPUTY_ADDRESS} ${BNBCHAIN_USER_ADDRESS} ${BNBCHAIN_DEPUTY_ADDRESS} now 1111111bnb 1111111bnb 360 true --from user
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

```bash
# the deputy will log an INFO message similar to:
2020-03-26 15:23:28 INFO sendBEP2HTLT send bep2 HTLT tx success, bnb_swap_id=67af6df0af8817d6c6240f8f7c2139df9e15185f324b08b53caed14f83d511a9, tx_hash=AD283B9F0E5147A7AD02B846AACE1709166173A1C2BC72D1FF8BCC8F140DFDEA
```

```bash
# query the new swap on kava by its ID
kvcli q bep3 swap ${KAVA_SWAP_ID}
```

After n blocks, the swap will created on bnbchain by the deputy

```bash
# the deputy will log an INFO message similar to:
INFO sendBEP2HTLT send bep2 HTLT tx success, bnb_swap_id=448a4cc0e1d2b4bce793919fcb2e557aae44d96bc715af8e1a110f774747667d, tx_hash=0DAD58181C6537B05394F46AA42FD9C73002A6E8601205FED6FBAEEDDAE7E1D1
```

The swap on bnbchain can now be claimed

```bash
tbnbcli token claim \
  --swap-id ${BNB_SWAP_ID} \
  --random-number ${RANDOM_NUMBER} \
  --from ${BNB_USER_ADDRESS} \
  --chain-id Binance-Chain-Nile \
  --node tcp://data-seed-pre-0-s1.binance.org:80 \
  --trust-node
```

Funds are transferred to the intended recipient's bnbchain address. After _n_ blocks the deputy will relay the successful claim to kava, closing the swap.
