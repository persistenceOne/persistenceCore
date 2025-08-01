package v12_0_0_rc0

import (
	"encoding/json"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

var denomsMap = map[string]struct{}{
	"stk/ubld":    {},
	"stk/uhuahua": {},
}

var contractAddr = "persistence1yljdn6nvt2k7dtz8erd9ytdaauhef5gtwwtevvtdqna0ms5afa0qunhryh"

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running upgrade handler")
		err := FixLSMData(ctx, args)
		if err != nil {
			return vm, err
		}
		ctx.Logger().Info("fixed lsm data")
		err = RedeemStkBalances(ctx, args)
		if err != nil {
			return vm, err
		}
		ctx.Logger().Info("redeemed stk balances")
		ctx.Logger().Info("running module migrations...")
		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

func FixLSMData(ctx sdk.Context, args upgrades.UpgradeHandlerArgs) error {
	k := args.Keepers.StakingKeeper
	err := k.RefreshTotalLiquidStaked(ctx)
	if err != nil {
		return err
	}
	for _, validator := range k.GetAllValidators(ctx) {
		validator.ValidatorBondShares = sdk.ZeroDec()
		k.SetValidator(ctx, validator)
	}

	// Sum up the total liquid tokens and increment each validator's liquid shares
	for _, delegation := range k.GetAllDelegations(ctx) {
		delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			return err
		}

		// disable valbond on delegator level.
		if k.DelegatorIsLiquidStaker(delegatorAddress) {
			delegation.ValidatorBond = false
			k.SetDelegation(ctx, delegation)
		}
	}

	return nil
}

func RedeemStkBalances(ctx sdk.Context, args upgrades.UpgradeHandlerArgs) error {
	userdata, err := GetUserAddressesWithBalance(ctx, args)
	if err != nil {
		return err
	}
	type Redeem struct{}
	type RedeemMsg struct {
		Redeem Redeem `json:"redeem"`
	}
	bz, err := json.Marshal(RedeemMsg{Redeem: Redeem{}})
	if err != nil {
		return err
	}
	for _, user := range userdata {
		msg := &wasmtypes.MsgExecuteContract{
			Sender:   user.Address,
			Contract: contractAddr,
			Msg:      bz,
			Funds:    sdk.NewCoins(user.Balance),
		}
		handler := args.Keepers.GovKeeper.Router().Handler(msg)
		res, err := handler(ctx, msg)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvents(res.GetEvents())
	}
	return nil
}

type AddressBalance struct {
	Address string
	Balance sdk.Coin
}

func GetUserAddressesWithBalance(ctx sdk.Context, args upgrades.UpgradeHandlerArgs) ([]AddressBalance, error) {

	contractAddrs := GetContractAddresses(ctx, args)
	ibcAddrs := GetIBCAddrs(ctx, args)
	// Define a struct to hold address and balance information
	userAddresses := make([]AddressBalance, 0)

	// Get all balances with non-zero amounts for the specific denom
	args.Keepers.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, balance sdk.Coin) bool {
		if _, ok := denomsMap[balance.Denom]; !ok {
			return false // continue iteration
		}

		// Get the account to check its type
		acc := args.Keepers.AccountKeeper.GetAccount(ctx, addr)
		if acc == nil {
			return false
		}

		if acc.GetPubKey() == nil {
			ctx.Logger().Info("found address without pubkey",
				"address", addr, "acc", acc)
		}

		// Skip module accounts
		if _, isModuleAccount := acc.(authtypes.ModuleAccountI); isModuleAccount {
			ctx.Logger().Info("found module account", "address", addr)
			return false
		}

		// Skip contract addresses (they start with specific prefix for contracts)
		if _, ok := contractAddrs[addr.String()]; ok {
			ctx.Logger().Info("found contract account", "address", addr)
			return false
		}

		// Skip IBC escrow accounts
		if _, ok := ibcAddrs[addr.String()]; ok {
			ctx.Logger().Info("found ibc escrow account", "address", addr)
			return false
		}

		userAddresses = append(userAddresses, AddressBalance{
			Address: addr.String(),
			Balance: balance,
		})

		return false // continue iteration
	})

	return userAddresses, nil
}

func GetContractAddresses(ctx sdk.Context, args upgrades.UpgradeHandlerArgs) map[string]struct{} {
	contracts := map[string]struct{}{}
	args.Keepers.WasmKeeper.IterateContractInfo(ctx, func(addr sdk.AccAddress, info wasmtypes.ContractInfo) bool {
		contracts[addr.String()] = struct{}{}
		return false
	})
	return contracts
}

// we are optimistic here
func GetIBCAddrs(ctx sdk.Context, args upgrades.UpgradeHandlerArgs) map[string]struct{} {
	addrs := map[string]struct{}{}
	channelID := args.Keepers.IBCKeeper.ChannelKeeper.GetNextChannelSequence(ctx)
	for i := uint64(0); i <= channelID; i++ {
		addr := transfertypes.GetEscrowAddress(transfertypes.PortID, channeltypes.FormatChannelIdentifier(i))
		addrs[addr.String()] = struct{}{}
	}
	return addrs
}
