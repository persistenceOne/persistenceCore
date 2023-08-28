package wasmbindings

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	liquidstakeibckeeper "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(bank *bankkeeper.BaseKeeper, liquidStakeIBC *liquidstakeibckeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:        old,
			bank:           bank,
			liquidStakeIBC: liquidStakeIBC,
		}
	}
}

type CustomMessenger struct {
	wrapped        wasmkeeper.Messenger
	bank           *bankkeeper.BaseKeeper
	liquidStakeIBC *liquidstakeibckeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

// DispatchMsg executes on the contractMsg.
func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {

		var contractMsg PstakeMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "pstake msg")
		}
		return contractMsg.Process(m.liquidStakeIBC, ctx, contractAddr)

	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

type PstakeMsg struct {
	LiquidStake    *liquidstakeibctypes.MsgLiquidStake    `json:"msg_liquid_stake,omitempty"`
	RedeemStake    *liquidstakeibctypes.MsgRedeem         `json:"msg_redeem_stake,omitempty"`
	LiquidStakeLSM *liquidstakeibctypes.MsgLiquidStakeLSM `json:"msg_liquid_stake_lsm,omitempty"`
}

func (msg PstakeMsg) Process(liquidStakeIBCKeeper *liquidstakeibckeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress) ([]sdk.Event, [][]byte, error) {
	msgServer := liquidstakeibckeeper.NewMsgServerImpl(*liquidStakeIBCKeeper)
	if msg.LiquidStake != nil {
		err := msg.LiquidStake.ValidateBasic()
		if err != nil {
			return nil, nil, err
		}
		msg.LiquidStake.DelegatorAddress = contractAddr.String()

		_, err = msgServer.LiquidStake(
			sdk.WrapSDKContext(ctx),
			msg.LiquidStake,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "liquid stake")
		}
		return nil, nil, nil
	}
	if msg.RedeemStake != nil {
		err := msg.RedeemStake.ValidateBasic()
		if err != nil {
			return nil, nil, err
		}
		msg.RedeemStake.DelegatorAddress = contractAddr.String()

		_, err = msgServer.Redeem(
			sdk.WrapSDKContext(ctx),
			msg.RedeemStake,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "redeem stake")
		}
		return nil, nil, nil
	}
	if msg.LiquidStakeLSM != nil {
		err := msg.LiquidStakeLSM.ValidateBasic()
		if err != nil {
			return nil, nil, err
		}
		msg.LiquidStakeLSM.DelegatorAddress = contractAddr.String()

		_, err = msgServer.LiquidStakeLSM(
			sdk.WrapSDKContext(ctx),
			msg.LiquidStakeLSM,
		)
		if err != nil {
			return nil, nil, errorsmod.Wrap(err, "liquid stake")
		}
		return nil, nil, nil
	}
	return nil, nil, errorsmod.Wrap(sdkerrors.ErrNotSupported, "only")
}
