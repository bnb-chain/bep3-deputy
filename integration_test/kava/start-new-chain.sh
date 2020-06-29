#! /bin/bash
set -e # fail on first error

validatorMnemonic="equip town gesture square tomorrow volume nephew minute witness beef rich gadget actress egg sing secret pole winter alarm law today check violin uncover"
# kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c
deputyMnemonic="slab twist stumble inmate predict parent repair crystal celery swarm memory loan rabbit blanket shell talk attend charge inside denial harbor music board steak"
# kava1sl8glhaa9f9tep0d9h8gdcfmwcatghtdrfcd2x
coldWalletMnemonic="skin wolf lemon pond lizard then drip garlic elegant clutch word domain vote topple alter assist hope fork teach shuffle define bright chuckle elbow"
# kava1kt283f9gtkyq3ndkd67yj3jlvlg94lgp7gcz33

# bnb cold wallet: merry gain mass calm judge border know fetch crouch depend deer over leg airport agree crisp case birth design patch truck butter patrol praise
# tbnb1y6kw4ztun2vp026ytp597nuqse7e3xxha2hyks

# Remove any existing data directory
rm -rf ~/.kvd
rm -rf ~/.kvcli

# Create new data directory
kvd init --chain-id=testing validator # doesn't need to be the same as the validator

kvcli config chain-id testing
# Set the cli to wait until a block is acceted before printing results
kvcli config broadcast-mode block
# avoid having to use password for keys
kvcli config keyring-backend test

# Create validator keys and add account to genesis
printf "$validatorMnemonic\n" | kvcli keys add validator --recover
kvd add-genesis-account $(kvcli keys show validator -a) 1000000000ukava
# Create faucet keys and add account to genesis
printf "$deputyMnemonic\n" | kvcli keys add deputy --recover
kvd add-genesis-account $(kvcli keys show deputy -a) 1000000000ukava,100000000000000bnb
# Create cold wallet fees
printf "$coldWalletMnemonic\n" | kvcli keys add cold-wallet --recover

# Create a delegation tx for the validator and add to genesis
kvd gentx --name validator --amount 100000000ukava --keyring-backend test
kvd collect-gentxs

# Replace stake with ukava
sed -in-place='' 's/stake/ukava/g' ~/.kvd/config/genesis.json

# start the blockchain in the background, record the process id so that it can be stopped, wait until it starts making blocks
kvd start --pruning everything --rpc.laddr "tcp://0.0.0.0:26657"
# & kvdPid="$!"
# sleep 6
# start the rest server. Ctrl-C  will stop both rest server and the blockchain (on crtl-c the kill thing runs and stops the blockchain process)
# kvcli rest-server ; kill $kvdPid