#! /bin/bash
set -e

config=~/Projects/Kava/bep3-deputy/config/test_config_kava.json
tmp=$(mktemp)

bnbchainheight=$(curl -s https://testnet-dex.binance.org/api/v1/node-info | jq '.sync_info.latest_block_height | tonumber')
echo bnbchain height: $bnbchainheight
setBnbchainHeight=$(( $bnbchainheight + 5))
jq '.chain_config.bnb_start_height = $newVal' --argjson newVal $setBnbchainHeight $config > "$tmp" && mv "$tmp" $config

kavaheight=$(curl -s http://127.0.0.1:26657/abci_info | jq '.result.response.last_block_height | tonumber')
echo kava height: $kavaheight
setKavaHeight=$(( $kavaheight ))
jq '.chain_config.other_chain_start_height = $newVal' --argjson newVal $setKavaHeight $config > "$tmp" && mv "$tmp" $config

cd ~/Projects/Kava/bep3-deputy
make build
cd ./build
rm -rf ./deputy.db
./deputy --bnb-network 0 --kava-network 0 --config-type local --config-path "../config/test_config_kava.json"