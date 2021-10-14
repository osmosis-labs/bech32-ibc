# bech32ibc

**bech32ibc** is a blockchain built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).

## Get started

```
starport serve
```

`serve` command installs dependencies, builds, initializes and starts your blockchain in development.

## Configure

Your blockchain in development can be configured with `config.yml`. To learn more see the [reference](https://github.com/tendermint/starport#documentation).

## Launch

To launch your blockchain live on mutliple nodes use `starport network` commands. Learn more about [Starport Network](https://github.com/tendermint/spn).

## What needs to be tested on Althea testnet

### Make IBC channels
### Make governance proposal to connect HRP and IBC channel

- Native HRP should be set in `bech32ibc` module's genesis
- Connect HRP to IBC Channel via governance proposal (`bech32ibc` module's `UpdateHrpIbcChannelProposal`), e.g. connect `osmo1` prefix to the IBC channel with Osmosis.

```sh
<daemon> tx bech32ibc update-hrp-ibc-record [human-readable-prefix] [channel-id] --title="set hrp for x network" --description="set hrp for x network description." --deposit="" --ics-to-height-offset=1000 ics-to-timeout-offset="0s" 
```

### Test out IBC sends to (1) a live chain, (2) a chain that is offline, and recover the funds that get stuck in Althea

Broadcast `banktypes.MsgSend` where target address is set to native chain address or altchain address - execution of these messages is handled by `bech32ics20` module.

```sh
bech32ibcd tx bank send validator <native_chain_or_altchain_address> 100uosmo --keyring-backend=test --chain-id=testing --yes
```

## Learn more

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Discord](https://discord.gg/W8trcGV)
