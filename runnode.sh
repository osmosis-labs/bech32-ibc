#!/bin/sh
set -o errexit -o nounset
CHAINID="testing"
rm -rf ~/.bech32ibc
coins="10000000000stake,100000000000samoleans"
bech32ibcd init --chain-id $CHAINID $CHAINID

bech32ibcd keys add validator --keyring-backend="test"
bech32ibcd add-genesis-account $(bech32ibcd keys show validator -a --keyring-backend="test")  $coins
bech32ibcd gentx validator 5000000000stake --keyring-backend="test" --chain-id $CHAINID
bech32ibcd collect-gentxs

sed -i '' 's#"nativeHRP": "osmo"#"nativeHRP": "stake"#g' ~/.bech32ibc/config/genesis.json
sed -i '' 's#"voting_period": "172800s"#"voting_period": "20s"#g' ~/.bech32ibc/config/genesis.json
sed -i '' 's#"quorum": "0.334000000000000000"#"quorum": "0.100000000000000000"#g' ~/.bech32ibc/config/genesis.json
sed -i '' 's#"threshold": "0.500000000000000000"#"threshold": "0.100000000000000000"#g' ~/.bech32ibc/config/genesis.json

sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.bech32ibc/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.bech32ibc/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.bech32ibc/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' ~/.bech32ibc/config/config.toml

bech32ibcd start --pruning=nothing 
