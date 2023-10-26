package v10_1_0_testnet

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"

	"github.com/persistenceOne/persistenceCore/v10/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running stkosmo upgrade handler")

		const HostChainID = "osmo-test-5"
		const HostChainMintDenom = "stk/uosmo"

		// Burn all coins on other addresses.
		args.Keepers.BankKeeper.IterateAllBalances(
			ctx,
			func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {

				if coin.Denom == HostChainMintDenom {
					// Send the whole address balance to the module account.
					if err := args.Keepers.BankKeeper.SendCoinsFromAccountToModule(
						ctx, address,
						liquidstakeibctypes.ModuleName,
						sdk.NewCoins(coin),
					); err != nil {
						return false
					}
				}

				return false
			},
		)

		// Burn all the coins present in the module account.
		if err := args.Keepers.BankKeeper.BurnCoins(
			ctx,
			liquidstakeibctypes.ModuleName,
			sdk.NewCoins(args.Keepers.BankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(liquidstakeibctypes.ModuleName), HostChainMintDenom)),
		); err != nil {
			return nil, err
		}

		// Remove the associated deposits from the store.
		mStore := prefix.NewStore(ctx.KVStore(args.Keepers.GetKey(liquidstakeibctypes.StoreKey)), liquidstakeibctypes.DepositKey)
		iterator := sdk.KVStorePrefixIterator(mStore, nil)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			deposit := liquidstakeibctypes.Deposit{}
			args.Codec.MustUnmarshal(iterator.Value(), &deposit)

			if deposit.ChainId == HostChainID {
				mStore.Delete(liquidstakeibctypes.GetDepositStoreKey(deposit.ChainId, deposit.Epoch))
			}
		}

		// Remove the associated LSM deposits from the store.
		mStore = prefix.NewStore(ctx.KVStore(args.Keepers.GetKey(liquidstakeibctypes.StoreKey)), liquidstakeibctypes.LSMDepositKey)
		iterator = sdk.KVStorePrefixIterator(mStore, nil)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			deposit := liquidstakeibctypes.LSMDeposit{}
			args.Codec.MustUnmarshal(iterator.Value(), &deposit)

			if deposit.ChainId == HostChainID {
				mStore.Delete(liquidstakeibctypes.GetLSMDepositStoreKey(deposit.ChainId, deposit.DelegatorAddress, deposit.Denom))
			}
		}

		// Remove the associated unbondings from the store.
		mStore = prefix.NewStore(ctx.KVStore(args.Keepers.GetKey(liquidstakeibctypes.StoreKey)), liquidstakeibctypes.UnbondingKey)
		iterator = sdk.KVStorePrefixIterator(mStore, nil)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			ub := liquidstakeibctypes.Unbonding{}
			args.Codec.MustUnmarshal(iterator.Value(), &ub)

			if ub.ChainId == HostChainID {
				mStore.Delete(liquidstakeibctypes.GetUnbondingStoreKey(ub.ChainId, ub.EpochNumber))
			}
		}

		// Remove the associated user unbondings from the store.
		mStore = prefix.NewStore(ctx.KVStore(args.Keepers.GetKey(liquidstakeibctypes.StoreKey)), liquidstakeibctypes.UserUnbondingKey)
		iterator = sdk.KVStorePrefixIterator(mStore, nil)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			ub := liquidstakeibctypes.UserUnbonding{}
			args.Codec.MustUnmarshal(iterator.Value(), &ub)

			if ub.ChainId == HostChainID {
				mStore.Delete(liquidstakeibctypes.GetUserUnbondingStoreKey(ub.ChainId, ub.Address, ub.EpochNumber))
			}
		}

		// Remove the associated validator unbondings from the store.
		mStore = prefix.NewStore(ctx.KVStore(args.Keepers.GetKey(liquidstakeibctypes.StoreKey)), liquidstakeibctypes.ValidatorUnbondingKey)
		iterator = sdk.KVStorePrefixIterator(mStore, nil)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			ub := liquidstakeibctypes.ValidatorUnbonding{}
			args.Codec.MustUnmarshal(iterator.Value(), &ub)

			if ub.ChainId == HostChainID {
				mStore.Delete(liquidstakeibctypes.GetValidatorUnbondingStoreKey(ub.ChainId, ub.ValidatorAddress, ub.EpochNumber))
			}
		}

		// Remove the host chain from the store.
		lsStore := prefix.NewStore(ctx.KVStore(args.Keepers.GetKey(liquidstakeibctypes.StoreKey)), liquidstakeibctypes.HostChainKey)
		lsStore.Delete([]byte(HostChainID))

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
