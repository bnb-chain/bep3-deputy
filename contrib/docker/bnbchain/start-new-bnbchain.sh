#! /usr/bin/bash

# Remove old keys and state
rm -rf ~/.bnbcli ~/.bnbchaind

# Setup keys
# deputy bnb1egvrzk6ujsh9t7cf79r0nwdzu8aze7j8893k5m
printf "password\nresemble volume attend machine expose behave amazing alone ten coconut sponsor endless employ grocery write physical diagram crisp bubble accuse six cry brown envelope\n" | bnbcli keys add deputy --recover
# user bnb13acaej0d676zya5q9ghz7hdpc4743aw5ekx2kf
printf "password\nfancy lazy report bird holiday original save early fun lunar secret enact also tennis sentence morning rebel program ocean income used ranch census next\n" | bnbcli keys add user --recover

# Init a new chain
bnbchaind init --moniker validatorName --chain-id Binance-Chain-Tigris --overwrite

# Add the two accounts to the genesis file (when the chain starts these will be populated with coins)
jq ".app_state.accounts[1].name=\"deputy\"" ~/.bnbchaind/config/genesis.json | sponge ~/.bnbchaind/config/genesis.json
jq ".app_state.accounts[1].address=\"$(bnbcli keys show deputy --address)\"" ~/.bnbchaind/config/genesis.json | sponge ~/.bnbchaind/config/genesis.json
jq ".app_state.accounts[2].name=\"user\"" ~/.bnbchaind/config/genesis.json | sponge ~/.bnbchaind/config/genesis.json
jq ".app_state.accounts[2].address=\"$(bnbcli keys show user --address)\"" ~/.bnbchaind/config/genesis.json | sponge ~/.bnbchaind/config/genesis.json
