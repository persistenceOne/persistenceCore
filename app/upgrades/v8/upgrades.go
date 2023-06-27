package v8

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icaMigrations "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/migrations/v6"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"

	"github.com/persistenceOne/persistenceCore/v8/app/keepers"
	"github.com/persistenceOne/persistenceCore/v8/app/upgrades"
)

func setInitialMinCommissionRate(ctx sdk.Context, keepers *keepers.AppKeepers) error {
	stakingParams := keepers.StakingKeeper.GetParams(ctx)
	stakingParams.MinCommissionRate = sdk.NewDecWithPrec(5, 2)
	if err := keepers.StakingKeeper.SetParams(ctx, stakingParams); err != nil {
		return fmt.Errorf("failed to set MinCommissionRate to 5%%: %w", err)
	}

	// Force update validator commission rate if it is lower than the minimum rate
	validators := keepers.StakingKeeper.GetAllValidators(ctx)
	for _, v := range validators {
		if v.Commission.Rate.LT(stakingParams.MinCommissionRate) {
			if v.Commission.MaxRate.LT(stakingParams.MinCommissionRate) {
				v.Commission.MaxRate = stakingParams.MinCommissionRate
			}

			v.Commission.Rate = stakingParams.MinCommissionRate
			v.Commission.UpdateTime = ctx.BlockHeader().Time

			// call the before-modification hook since we're about to update the commission
			if err := keepers.StakingKeeper.Hooks().BeforeValidatorModified(ctx, v.GetOperator()); err != nil {
				return fmt.Errorf("BeforeValidatorModified failed with: %w", err)
			}
			keepers.StakingKeeper.SetValidator(ctx, v)
		}
	}

	return nil
}

func setOraclePairListEmpty(ctx sdk.Context, keepers *keepers.AppKeepers) {
	oracleParams := keepers.OracleKeeper.GetParams(ctx)
	oracleParams.AcceptList = oracletypes.DenomList{}
	keepers.OracleKeeper.SetParams(ctx, oracleParams)
}

func disableMEVAuction(ctx sdk.Context, keepers *keepers.AppKeepers) error {
	builderParams, err := keepers.BuilderKeeper.GetParams(ctx)
	if err != nil {
		return err
	}
	// Setting MaxBundleSize to 0 means no auction txs will be accepted
	builderParams.MaxBundleSize = 0
	return keepers.BuilderKeeper.SetParams(ctx, builderParams)
}

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	baseAppLegacySS := getLegacySubspaces(args.Keepers.ParamsKeeper)

	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		// ibc v4-to-v5
		// -- nothing --

		// ibc v5-to-v6
		ctx.Logger().Info("migrating ics27 channel capability")
		err := icaMigrations.MigrateICS27ChannelCapability(
			ctx,
			args.Codec,
			args.Keepers.GetKey(capabilitytypes.StoreKey),
			args.Keepers.CapabilityKeeper,
			lscosmostypes.ModuleName,
		)
		if err != nil {
			return nil, err
		}

		// ibc v6-to-v7
		// (optional) prune expired tendermint consensus states to save storage space
		ctx.Logger().Info("pruning expired tendermint consensus states for ibc clients")
		_, err = ibctmmigrations.PruneExpiredConsensusStates(ctx, args.Codec, args.Keepers.IBCKeeper.ClientKeeper)
		if err != nil {
			return nil, err
		}

		// ibc v7-to-v7.1
		// explicitly update the IBC 02-client params, adding the localhost client type
		ctx.Logger().Info("adding localhost client to IBC params")
		params := args.Keepers.IBCKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		args.Keepers.IBCKeeper.ClientKeeper.SetParams(ctx, params)

		// sdk v45-to-v46
		// -- nothing --

		// sdk v46-to-v47
		// Migrate Tendermint consensus parameters from x/params module to a dedicated x/consensus module.
		ctx.Logger().Info("migrating tendermint consensus params")
		baseapp.MigrateParams(ctx, baseAppLegacySS, args.Keepers.ConsensusParamsKeeper)

		// Note: this migration is optional,
		// You can include x/gov proposal migration documented in [UPGRADING.md](https://github.com/cosmos/cosmos-sdk/blob/main/UPGRADING.md)
		// TODO(ajeet): do we need this optional migration?

		ctx.Logger().Info("running module migrations")
		newVm, err := args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
		if err != nil {
			return nil, err
		}

		enabled := args.Keepers.LSCosmosKeeper.GetModuleState(ctx)
		if enabled {
			ctx.Logger().Info("migrating lscsomos module")
			if err = args.Keepers.LSCosmosKeeper.Migrate(ctx); err != nil {
				return nil, err
			}
		}

		ctx.Logger().Info("setting min commission rate to 5%")
		if err = setInitialMinCommissionRate(ctx, args.Keepers); err != nil {
			return nil, err
		}

		ctx.Logger().Info("setting acceptList to empty in oracle params")
		setOraclePairListEmpty(ctx, args.Keepers)

		ctx.Logger().Info("disable auction for MEV")
		if err = disableMEVAuction(ctx, args.Keepers); err != nil {
			return nil, err
		}

		// TODO(ajeet): do we need to set gov -> MinInitialDepositRatio? (default is 0 -> disabled)

		return newVm, nil
	}
}
