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

## Process of sending assets via `bech32ibc` and `bech32ics20`

- Native HRP should be set in `bech32ibc` module's genesis
- Connect HRP to IBC Channel via governance proposal (`bech32ibc` module's `UpdateHrpIbcChannelProposal`), e.g. connect `osmo1` prefix to the IBC channel with Osmosis.
- Broadcast `MsgSend` or `MsgMultiSend` target address set to native chain address or altchain address - execution of these messages is handled by `bech32ics20` module.

## Learn more

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Discord](https://discord.gg/W8trcGV)
