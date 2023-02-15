package wasmbindings

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	oraclekeeper "github.com/persistenceOne/persistence-sdk/v2/x/oracle/keeper"

	"github.com/persistenceOne/persistenceCore/v7/wasmbindings/bindings"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(checkersKeeper *oraclekeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:        old,
			checkersKeeper: checkersKeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped        wasmkeeper.Messenger
	checkersKeeper *oraclekeeper.Keeper
}

func (c CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.CheckersMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "checkers msg")
		}

		if contractMsg.UpdateExchangeRate != nil {
			// TODO: do this later since we are checking only query
			return nil, nil, nil
		}
	}

	return c.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}
