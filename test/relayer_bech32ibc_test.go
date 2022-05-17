package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/v2/relayer"
	bech32ibctypes "github.com/osmosis-labs/bech32-ibc/x/bech32ibc/types"
	"github.com/stretchr/testify/require"
)

var (
	bech32ibcChains = []testChain{
		{"ibc-0", 0, gaiaTestConfig},
		{"ibc-1", 1, bech32ibcTestConfig},
	}
)

// QueryHrpIbcRecords queries hrp ibc records
func QueryHrpIbcRecords(c *relayer.Chain) ([]bech32ibctypes.HrpIbcRecord, error) {
	done := c.UseSDKContext()
	done()

	params := &bech32ibctypes.QueryHrpIbcRecordsRequest{}
	queryClient := bech32ibctypes.NewQueryClient(c.CLIContext(0))

	res, err := queryClient.HrpIbcRecords(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return res.HrpIbcRecords, nil
}

// QueryProposals query proposals
func QueryProposals(c *relayer.Chain) ([]govtypes.Proposal, error) {
	done := c.UseSDKContext()
	done()

	params := &govtypes.QueryProposalsRequest{}
	queryClient := govtypes.NewQueryClient(c.CLIContext(0))

	res, err := queryClient.Proposals(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return res.Proposals, nil
}

func QueryValidators(c *relayer.Chain) ([]stakingtypes.Validator, error) {
	done := c.UseSDKContext()
	done()

	params := &stakingtypes.QueryValidatorsRequest{}
	queryClient := stakingtypes.NewQueryClient(c.CLIContext(0))

	res, err := queryClient.Validators(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return res.Validators, nil
}

func TestBech32IBCStreamingRelayer(t *testing.T) {
	chains := spinUpTestChains(t, bech32ibcChains...)

	var (
		src            = chains.MustGet("ibc-0")
		dst            = chains.MustGet("ibc-1")
		testDenom      = "samoleans"
		testCoin       = sdk.NewCoin(testDenom, sdk.NewInt(1000))
		twoTestCoin    = sdk.NewCoin(testDenom, sdk.NewInt(2000))
		initialDeposit = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(20000000))
		delegationAmt  = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3000000000))
	)

	path, err := genTestPathAndSet(src, dst, "transfer", "transfer")
	require.NoError(t, err)

	// query initial balances to compare against at the end
	srcExpected, err := src.QueryBalance(src.Key)
	require.NoError(t, err)
	dstExpected, err := dst.QueryBalance(dst.Key)
	require.NoError(t, err)

	// create path
	_, err = src.CreateClients(dst, true, true, false)
	require.NoError(t, err)
	testClientPair(t, src, dst)

	_, err = src.CreateOpenConnections(dst, 3, src.GetTimeout())
	require.NoError(t, err)
	testConnectionPair(t, src, dst)

	_, err = src.CreateOpenChannels(dst, 3, src.GetTimeout())
	require.NoError(t, err)
	testChannelPair(t, src, dst)

	// send a couple of transfers to the queue on src
	require.NoError(t, src.SendTransferMsg(dst, testCoin, dst.MustGetAddress(), 0, 0))
	require.NoError(t, src.SendTransferMsg(dst, testCoin, dst.MustGetAddress(), 0, 0))

	// send a couple of transfers to the queue on dst
	require.NoError(t, dst.SendTransferMsg(src, testCoin, src.MustGetAddress(), 0, 0))
	require.NoError(t, dst.SendTransferMsg(src, testCoin, src.MustGetAddress(), 0, 0))

	validators, err := QueryValidators(dst)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(validators), 1)

	resp, _, err := dst.SendMsg(&stakingtypes.MsgDelegate{
		DelegatorAddress: dst.MustGetAddress(),
		ValidatorAddress: validators[0].OperatorAddress,
		Amount:           delegationAmt,
	})
	require.NoError(t, err)
	dst.Log(fmt.Sprintln("MsgDelegate.Response", resp.Logs))

	// Native HRP is set to "stake" as part of genesis in `bech32ibc-setup.sh`
	// Send a proposal to connect hrp with channel
	msg := &govtypes.MsgSubmitProposal{
		InitialDeposit: sdk.Coins{initialDeposit},
		Proposer:       dst.MustGetAddress(),
	}
	err = msg.SetContent(bech32ibctypes.NewUpdateHrpIBCRecordProposal(
		"set hrp for gaia network",
		"set hrp for gaia network",
		gaiaTestConfig.accountPrefix,
		dst.PathEnd.ChannelID,
		1000, 0,
	))
	require.NoError(t, err)
	resp, _, err = dst.SendMsg(msg)
	require.NoError(t, err)

	dst.Log(fmt.Sprintln("MsgSubmitProposal.Response", resp.Logs))

	proposalIDStr := resp.Logs[0].Events[2].Attributes[0].Value
	dst.Log(fmt.Sprintln("ProposalIDStr", proposalIDStr))
	proposalID, err := strconv.Atoi(proposalIDStr)
	require.NoError(t, err)

	// approve the proposal
	resp, _, err = dst.SendMsg(&govtypes.MsgVote{
		Voter:      dst.MustGetAddress(),
		ProposalId: uint64(proposalID),
		Option:     govtypes.OptionYes,
	})
	require.NoError(t, err)

	dst.Log(fmt.Sprintln("MsgVote.Response", resp.Logs))

	// wait for voting period
	dst.WaitForNBlocks(20)

	dst.Log(fmt.Sprintln("Log after 20 blocks"))

	proposals, err := QueryProposals(dst)
	require.NoError(t, err)
	dst.Log(fmt.Sprintln("proposals.Response", proposals))

	// check hrp is updated correctly
	hrpRecords, err := QueryHrpIbcRecords(dst)
	require.NoError(t, err)

	dst.Log(fmt.Sprintln("hrpRecords.Response", hrpRecords))

	// check balance changes
	_, _, err = dst.SendMsg(&banktypes.MsgSend{
		FromAddress: dst.MustGetAddress(),
		ToAddress:   dst.MustGetAddress(),
		Amount:      sdk.Coins{testCoin},
	})
	require.NoError(t, err)

	// check balance changes
	_, _, err = dst.SendMsg(&banktypes.MsgSend{
		FromAddress: dst.MustGetAddress(),
		ToAddress:   src.MustGetAddress(),
		Amount:      sdk.Coins{testCoin},
	})
	require.NoError(t, err)

	// Wait for message inclusion in both chains
	require.NoError(t, dst.WaitForNBlocks(1))

	// start the relayer process in it's own goroutine
	rlyDone, err := relayer.RunStrategy(src, dst, path.MustGetStrategy())
	require.NoError(t, err)

	// Wait for relay message inclusion in both chains
	require.NoError(t, src.WaitForNBlocks(1))
	require.NoError(t, dst.WaitForNBlocks(1))

	// send those tokens from dst back to dst and src back to src
	require.NoError(t, src.SendTransferMsg(dst, twoTestCoin, dst.MustGetAddress(), 0, 0))
	require.NoError(t, dst.SendTransferMsg(src, twoTestCoin, src.MustGetAddress(), 0, 0))

	// wait for packet processing
	require.NoError(t, dst.WaitForNBlocks(6))

	// kill relayer routine
	rlyDone()

	// check balance on src against expected
	srcGot, err := src.QueryBalance(src.Key)
	require.NoError(t, err)
	require.Equal(t, srcExpected.AmountOf(testDenom).Int64()-4000, srcGot.AmountOf(testDenom).Int64())

	// check balance on dst against expected
	dstGot, err := dst.QueryBalance(dst.Key)
	require.NoError(t, err)
	require.Equal(t, dstExpected.AmountOf(testDenom).Int64()-5000, dstGot.AmountOf(testDenom).Int64())

	// check balance on src against expected
	srcGot, err = src.QueryBalance(src.Key)
	require.NoError(t, err)
	require.Equal(t, srcExpected.AmountOf(testDenom).Int64()-4000, srcGot.AmountOf(testDenom).Int64())

	// check balance on dst against expected
	dstGot, err = dst.QueryBalance(dst.Key)
	require.NoError(t, err)
	require.Equal(t, dstExpected.AmountOf(testDenom).Int64()-5000, dstGot.AmountOf(testDenom).Int64())
}
