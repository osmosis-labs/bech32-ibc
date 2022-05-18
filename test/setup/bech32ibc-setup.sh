#!/bin/sh

set -o errexit -o nounset

CHAINID=$1
GENACCT=$2
PRIVPATH=$3

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi

if [ -z "$3" ]; then
  echo "Need to input path of priv_validator_key json file"
  exit 1
fi

# Build genesis file incl account for passed address
coins="10000000000stake,100000000000samoleans"
bech32ibcd init --chain-id $CHAINID $CHAINID
bech32ibcd keys add validator --keyring-backend="test"
bech32ibcd add-genesis-account $(bech32ibcd keys show validator -a --keyring-backend="test") $coins
bech32ibcd add-genesis-account $GENACCT $coins
cp $PRIVPATH ~/.bech32ibc/config/priv_validator_key.json
bech32ibcd gentx validator 5000000000stake --keyring-backend="test" --chain-id $CHAINID
bech32ibcd collect-gentxs

# modify genesis of bech32ibc module
# sed -i 's#"nativeHRP": "osmo"#"nativeHRP": "akash"#g' ~/.bech32ibc/config/genesis.json

# modify genesis of governance voting period to be faster (20s)
sed -i 's#"voting_period": "172800s"#"voting_period": "20s"#g' ~/.bech32ibc/config/genesis.json
sed -i 's#"quorum": "0.334000000000000000"#"quorum": "0.100000000000000000"#g' ~/.bech32ibc/config/genesis.json
sed -i 's#"threshold": "0.500000000000000000"#"threshold": "0.100000000000000000"#g' ~/.bech32ibc/config/genesis.json

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.bech32ibc/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.bech32ibc/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.bech32ibc/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.bech32ibc/config/config.toml

# Start the bech32ibcd
bech32ibcd start --pruning=nothing
