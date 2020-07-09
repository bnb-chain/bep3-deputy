#! /bin/bash

set -e

BNCHOME=${HOME}/.bnbchaind

# Setup keys
# deputy bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7
printf "password\nclinic soap symptom alter mango orient punch table seek among broken bundle best dune hurt predict liquid subject silver once kick metal okay moment\n" | bnbcli keys add deputy --recover
# cold wallet bnb13acaej0d676zya5q9ghz7hdpc4743aw5ekx2kf
printf "password\nfancy lazy report bird holiday original save early fun lunar secret enact also tennis sentence morning rebel program ocean income used ranch census next\n" | bnbcli keys add cold-wallet --recover
# test user bnb18s9h2nsjlecwpntgyjd9g7h4wx5zez899azk04
printf "password\nthen nuclear favorite advance plate glare shallow enhance replace embody list dose quick scale service sentence hover announce advance nephew phrase order useful this\n" | bnbcli keys add test-user --recover

# Init a new chain
bnbchaind init --moniker validatorName --chain-id Binance-Chain-Tigris --overwrite --home $BNCHOME

# Add the two accounts to the genesis file (when the chain starts these will be populated with coins)
jq ".app_state.accounts[1].name=\"deputy\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json
jq ".app_state.accounts[1].address=\"$(bnbcli keys show deputy --address)\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json

jq ".app_state.accounts[2].name=\"test_user\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json
jq ".app_state.accounts[2].address=\"$(bnbcli keys show test-user --address)\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json

# Turn on console logging
# sed -i 's/logToConsole = false/logToConsole = true/g' ${BNCHOME}/config/app.toml

# Start chain
bnbchaind --home $BNCHOME start