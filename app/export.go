package app

import (
	storetypes "cosmossdk.io/store/types"
	"fmt"
	stdlog "log"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (app *Application) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string, modulesToExport []string) (servertypes.ExportedApp, error) {
	context := app.NewContextLegacy(true, tmproto.Header{Height: app.LastBlockHeight()})

	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(context, jailWhiteList)
	}

	genesisState, err := app.moduleManager.ExportGenesisForModules(context, app.appCodec, modulesToExport)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	applicationState, Error := codec.MarshalJSONIndent(app.legacyAmino, genesisState)

	if Error != nil {
		return servertypes.ExportedApp{}, Error
	}

	validators, err := staking.WriteValidators(context, app.StakingKeeper)

	return servertypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(context),
	}, err
}

func (app *Application) prepForZeroHeightGenesis(context sdk.Context, jailWhiteList []string) {
	/* Just to be safe, assert the invariants on current state. */
	app.CrisisKeeper.AssertInvariants(context)

	/* Handle fee distribution state. */

	// withdraw all validator commission
	err := app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		valBz, err := app.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}
		_, _ = app.DistributionKeeper.WithdrawValidatorCommission(context, valBz)
		return false
	})
	if err != nil {
		panic(err)
	}

	delegations, err := app.StakingKeeper.GetAllDelegations(context)
	if err != nil {
		panic(err)
	}

	app.withdrawDelegationRewards(context, delegations)

	app.DistributionKeeper.DeleteAllValidatorSlashEvents(context)
	app.DistributionKeeper.DeleteAllValidatorHistoricalRewards(context)

	// set context height to zero
	height := context.BlockHeight()
	context = context.WithBlockHeight(0)

	app.reinitializeValidators(context)
	app.reinitializeDelegations(context, delegations)

	// reset context height
	context = context.WithBlockHeight(height)

	/* Handle staking state. */

	// reset creation height for redelegations & unbonding delegations
	app.resetCreationHeightForDelEntries(context)

	if err = app.resetValidatorBondHeights(context, jailWhiteList); err != nil {
		app.Logger().Error(err.Error())
		return
	}

	/* Handle slashing state. */

	// reset start height on signing infos
	err = app.SlashingKeeper.IterateValidatorSigningInfos(
		context,
		func(validatorConsAddress sdk.ConsAddress, validatorSigningInfo slashingtypes.ValidatorSigningInfo) (stop bool) {
			validatorSigningInfo.StartHeight = 0
			err = app.SlashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
			if err != nil {
				panic(err)
			}
			return false
		},
	)
	if err != nil {
		panic(err)
	}
}

func (app *Application) withdrawDelegationRewards(context sdk.Context, delegations []stakingtypes.Delegation) {
	for _, delegation := range delegations {
		validatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		_, err = app.DistributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
		if err != nil {
			panic(err)
		}
	}
}

func (app *Application) reinitializeValidators(context sdk.Context) {
	err := app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		valBz, err := app.StakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}
		scraps, err := app.DistributionKeeper.GetValidatorOutstandingRewardsCoins(context, valBz)
		if err != nil {
			panic(err)
		}
		feePool, err := app.DistributionKeeper.FeePool.Get(context)
		if err != nil {
			panic(err)
		}
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		err = app.DistributionKeeper.FeePool.Set(context, feePool)
		if err != nil {
			panic(err)
		}

		if err := app.DistributionKeeper.Hooks().AfterValidatorCreated(context, valBz); err != nil {
			// never called as AfterValidatorCreated always returns nil
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}
}

func (app *Application) reinitializeDelegations(context sdk.Context, delegations []stakingtypes.Delegation) {
	for _, delegation := range delegations {
		validatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		if err := app.DistributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress); err != nil {
			// never called as BeforeDelegationCreated always returns nil
			panic(fmt.Errorf("error while incrementing period: %w", err))
		}

		if err := app.DistributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress); err != nil {
			// never called as AfterDelegationModified always returns nil
			panic(fmt.Errorf("error while creating a new delegation period record: %w", err))
		}
	}
}

func (app *Application) resetCreationHeightForDelEntries(context sdk.Context) {
	// iterate through redelegations, reset creation height
	err := app.StakingKeeper.IterateRedelegations(context, func(_ int64, redelegation stakingtypes.Redelegation) (stop bool) {
		for i := range redelegation.Entries {
			redelegation.Entries[i].CreationHeight = 0
		}
		err := app.StakingKeeper.SetRedelegation(context, redelegation)
		if err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	// iterate through unbonding delegations, reset creation height
	err = app.StakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range unbondingDelegation.Entries {
			unbondingDelegation.Entries[i].CreationHeight = 0
		}
		err = app.StakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
		if err != nil {
			panic(err)
		}
		return false
	})
	if err != nil {
		panic(err)
	}
}

func (app *Application) resetValidatorBondHeights(context sdk.Context, jailWhiteList []string) error {
	applyWhiteList := len(jailWhiteList) > 0
	whiteListMap := getWhilteListMap(jailWhiteList)

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := context.KVStore(app.GetKVStoreKey()[stakingtypes.StoreKey])
	iter := storetypes.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(stakingtypes.AddressFromValidatorsKey(iter.Key()))
		validator, err := app.StakingKeeper.GetValidator(context, addr)

		if err != nil {
			panic("Validator not found!")
		}

		validator.UnbondingHeight = 0
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		app.StakingKeeper.SetValidator(context, validator)
		counter++
	}

	if err := iter.Close(); err != nil {
		return fmt.Errorf("error while closing the key-value store reverse prefix iterator: %w", err)
	}

	_, err := app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(context)
	if err != nil {
		stdlog.Fatal(err)
	}

	return nil
}

func getWhilteListMap(jailWhiteList []string) map[string]bool {
	whiteListMap := make(map[string]bool)

	for _, address := range jailWhiteList {
		if _, err := sdk.ValAddressFromBech32(address); err != nil {
			panic(err)
		}

		whiteListMap[address] = true
	}
	return whiteListMap
}
