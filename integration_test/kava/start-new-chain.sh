#! /bin/bash
set -e # fail on first error

validatorMnemonic="law assault face proud fan slim genius boring portion delay team rude vapor timber noble absorb laugh dilemma patch actress brisk tissue drift flock"
# kava179k8at2krka7snzmp7tpdvl0p8zssu0yvycxc5
deputyMnemonic="equip town gesture square tomorrow volume nephew minute witness beef rich gadget actress egg sing secret pole winter alarm law today check violin uncover"
# kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c
coldWalletMnemonic="skin wolf lemon pond lizard then drip garlic elegant clutch word domain vote topple alter assist hope fork teach shuffle define bright chuckle elbow"
# kava1kt283f9gtkyq3ndkd67yj3jlvlg94lgp7gcz33
testUserMnemonic="very health column only surface project output absent outdoor siren reject era legend legal twelve setup roast lion rare tunnel devote style random food"
# kava1ypjp0m04pyp73hwgtc0dgkx0e9rrydecm054da

# Remove any existing data directory
rm -rf ~/.kvd
rm -rf ~/.kvcli

# Create new data directory
kvd init --chain-id=testing validator # doesn't need to be the same as the validator
# Copy in template genesis file
cp ~/kava/config/genesis.json ~/.kvd/config/genesis.json

kvcli config chain-id testing
# Set the cli to wait until a block is acceted before printing results
kvcli config broadcast-mode block
# avoid having to use password for keys
kvcli config keyring-backend test

# Create validator keys and add account to genesis
printf "$validatorMnemonic\n" | kvcli keys add validator --recover
kvd add-genesis-account $(kvcli keys show validator -a) 1000000000ukava
# Create deputy keys and add account to genesis
printf "$deputyMnemonic\n" | kvcli keys add deputy --recover
kvd add-genesis-account $(kvcli keys show deputy -a) 10000000ukava,100000000000000bnb
# TODO kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c
# kavapub1addwnpepqfjcpzn0xapflfpperm520hzmahdhhjaqs0ses4vv4r9a6cmaaq8kq7r6q0
# # Create deputy keys but don't add account to genesis
printf "$coldWalletMnemonic\n" | kvcli keys add cold-wallet --recover
# Create test user keys and add account to genesis
printf "$testUserMnemonic\n" | kvcli keys add test-user --recover
kvd add-genesis-account $(kvcli keys show test-user -a) 10000000ukava

# Create a delegation tx for the validator and add to genesis
kvd gentx --name validator --amount 100000000ukava --keyring-backend test
kvd collect-gentxs

# Sanity check to make sure genesis hasn't got messed up
kvd validate-genesis

# start the blockchain in the background, record the process id so that it can be stopped, wait until it starts making blocks
kvd start --pruning nothing --rpc.laddr "tcp://0.0.0.0:26657"