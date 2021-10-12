package bech32ics20

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/osmosis-labs/bech32-ibc/x/bech32ics20/keeper"
)

// NewHandler returns a handler for "bech32ics20" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	// msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized bech32ics20 message type: %T", msg)
		}
	}
}
