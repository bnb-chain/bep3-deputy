#! /bin/bash

set -e

BNCHOME=${HOME}/.bnbchaind

# Setup keys
# deputy bnb1uky3me9ggqypmrsvxk7ur6hqkzq7zmv4ed4ng7
printf "password\nclinic soap symptom alter mango orient punch table seek among broken bundle best dune hurt predict liquid subject silver once kick metal okay moment\n" | bnbcli keys add deputy --recover
# cold wallet bnb13acaej0d676zya5q9ghz7hdpc4743aw5ekx2kf
printf "password\nfancy lazy report bird holiday original save early fun lunar secret enact also tennis sentence morning rebel program ocean income used ranch census next\n" | bnbcli keys add cold-wallet --recover
testUserMnemonics=()
# bnb18s9h2nsjlecwpntgyjd9g7h4wx5zez899azk04
testUserMnemonics[0]="then nuclear favorite advance plate glare shallow enhance replace embody list dose quick scale service sentence hover announce advance nephew phrase order useful this"
# bnb1zfa5vmsme2v3ttvqecfleeh2xtz5zghh49hfqe
testUserMnemonics[1]="almost design doctor exist destroy candy zebra insane client grocery govern idea library degree two rebuild coffee hat scene deal average fresh measure potato"
# bnb1nva9yljftdf6m2dwhufk5kzg204jg060sw0fv2
testUserMnemonics[2]="welcome bean crystal pave chapter process bless tribe inside bottom exhaust hollow display envelope rally moral admit round hidden junk silly afraid awesome muffin"
# bnb1m5k0g7q0n8nzsre8ysuhsv09j03jhgcrrnwur3
testUserMnemonics[3]="end bicycle walnut empty bus silly camera lift fancy symptom office pluck detail unable cry sense scrap tuition relax amateur hold win debate hat"
# bnb15udkjukldcejs3y3pep2jvwza0lfzwshv7ndfr
testUserMnemonics[4]="cloud deal hurdle sound scout merit carpet identify fossil brass ancient keep disorder save lobster whisper course intact winter bullet flame mother upgrade install"
# bnb1gzr37hrqlqk4wpdfqn6pp3dy9ek28tahy7dx2v
testUserMnemonics[5]="mutual duck begin remind release brave patrol squeeze abandon pact valid close fragile plastic disorder saddle bring inspire corn kitten reduce candy side honey"
for i in {0..5}; do
    printf "password\n${testUserMnemonics[$i]}\n" | bnbcli keys add test-user$i --recover
done

# Init a new chain
bnbchaind init --moniker validatorName --chain-id Binance-Chain-Tigris --overwrite --home $BNCHOME

# Add the accounts to the genesis file (when the chain starts these will be populated with coins)
jq ".app_state.accounts[1].name=\"deputy\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json
jq ".app_state.accounts[1].address=\"$(bnbcli keys show deputy --address)\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json

for i in {0..5}; do
    jq ".app_state.accounts[$((2+$i))].name=\"test_user_$i\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json
    jq ".app_state.accounts[$((2+$i))].address=\"$(bnbcli keys show test-user$i --address)\"" ${BNCHOME}/config/genesis.json | sponge ${BNCHOME}/config/genesis.json
done

# Turn on console logging
# sed -i 's/logToConsole = false/logToConsole = true/g' ${BNCHOME}/config/app.toml

# Start chain
bnbchaind --home $BNCHOME start