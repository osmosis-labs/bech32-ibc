package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) banktypes.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ banktypes.MsgServer = msgServer{}

func (k msgServer) Send(goCtx context.Context, msg *banktypes.MsgSend) (*banktypes.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.SendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	prefix, _, err := bech32.DecodeAndConvert(msg.ToAddress)
	if err != nil {
		return nil, err
	}

	nativePrefix, err := k.hrpToChannelMapper.GetNativeHrp(ctx)
	if err != nil {
		return nil, err
	}

	if prefix == nativePrefix {
		bankMsgServer := bankkeeper.NewMsgServerImpl(k.Keeper.Keeper)
		return bankMsgServer.Send(goCtx, msg)
	}

	ibcRecord, err := k.hrpToChannelMapper.GetHrpIbcRecord(ctx, prefix)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Len() == 0 {
		return nil, sdkerrors.Wrap(banktypes.ErrNoInputs, "invalid send amount")
	}
	if msg.Amount.Len() > 1 {
		return nil, sdkerrors.Wrap(ibctransfertypes.ErrInvalidAmount, "cannot send multiple denoms via IBC")
	}

	portId := k.tk.GetPort(ctx)
	_, clientState, err := k.channelKeeper.GetChannelClientState(ctx, portId, ibcRecord.SourceChannel)
	if err != nil {
		return nil, err
	}

	latestHeight := clientState.GetLatestHeight()
	timeoutHeight := clienttypes.NewHeight(latestHeight.GetRevisionNumber(), latestHeight.GetRevisionHeight()+ibcRecord.IcsToHeightOffset)

	ibcTransferMsg := ibctransfertypes.NewMsgTransfer(
		portId,
		ibcRecord.SourceChannel,
		msg.Amount[0],
		from,
		msg.ToAddress,
		timeoutHeight, 0, // Use no timeouts for now.  Can add this in future.
	)

	_, err = k.ics20TransferMsgServer.Transfer(sdk.WrapSDKContext(ctx), ibcTransferMsg)

	return &banktypes.MsgSendResponse{}, err
}

func (k msgServer) MultiSend(goCtx context.Context, msg *banktypes.MsgMultiSend) (*banktypes.MsgMultiSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: totalIn == totalOut should already have been checked
	for _, in := range msg.Inputs {
		if err := k.SendEnabledCoins(ctx, in.Coins...); err != nil {
			return nil, err
		}
	}

	for _, out := range msg.Outputs {
		accAddr, err := sdk.AccAddressFromBech32(out.Address)
		if err != nil {
			panic(err)
		}
		if k.BlockedAddr(accAddr) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", out.Address)
		}
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, banktypes.AttributeValueCategory),
		),
	)

	return &banktypes.MsgMultiSendResponse{}, nil
}
