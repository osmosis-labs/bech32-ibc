package keeper

import (
	"github.com/osmosis-labs/bech32-ibc/x/bech32ibc/types"
)

var _ types.QueryServer = Keeper{}
