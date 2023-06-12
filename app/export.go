package app

import (
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
	context := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(context, jailWhiteList)
	}

	genesisState := app.moduleManager.ExportGenesisForModules(context, app.applicationCodec, modulesToExport)
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
	applyWhiteList := false

	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, address := range jailWhiteList {
		if _, err := sdk.ValAddressFromBech32(address); err != nil {
			panic(err)
		}

		whiteListMap[address] = true
	}

	app.CrisisKeeper.AssertInvariants(context)

	app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		_, _ = app.DistributionKeeper.WithdrawValidatorCommission(context, val.GetOperator())
		return false
	})

	delegations := app.StakingKeeper.GetAllDelegations(context)
	for _, delegation := range delegations {
		validatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		_, _ = app.DistributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
	}

	app.DistributionKeeper.DeleteAllValidatorSlashEvents(context)

	app.DistributionKeeper.DeleteAllValidatorHistoricalRewards(context)

	height := context.BlockHeight()
	context = context.WithBlockHeight(0)

	app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		scraps := app.DistributionKeeper.GetValidatorOutstandingRewardsCoins(context, val.GetOperator())
		feePool := app.DistributionKeeper.GetFeePool(context)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.DistributionKeeper.SetFeePool(context, feePool)

		app.DistributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator())
		return false
	})

	for _, delegation := range delegations {
		validatorAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		app.DistributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress)
		app.DistributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress)
	}

	context = context.WithBlockHeight(height)

	app.StakingKeeper.IterateRedelegations(context, func(_ int64, redelegation stakingtypes.Redelegation) (stop bool) {
		for i := range redelegation.Entries {
			redelegation.Entries[i].CreationHeight = 0
		}
		app.StakingKeeper.SetRedelegation(context, redelegation)
		return false
	})

	app.StakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range unbondingDelegation.Entries {
			unbondingDelegation.Entries[i].CreationHeight = 0
		}
		app.StakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
		return false
	})

	store := context.KVStore(app.GetKVStoreKey()[stakingtypes.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(stakingtypes.AddressFromValidatorsKey(iter.Key()))
		validator, found := app.StakingKeeper.GetValidator(context, addr)

		if !found {
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
		app.Logger().Error("error while closing the key-value store reverse prefix iterator: ", err)
		return
	}

	_, err := app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(context)
	if err != nil {
		stdlog.Fatal(err)
	}

	app.SlashingKeeper.IterateValidatorSigningInfos(
		context,
		func(validatorConsAddress sdk.ConsAddress, validatorSigningInfo slashingtypes.ValidatorSigningInfo) (stop bool) {
			validatorSigningInfo.StartHeight = 0
			app.SlashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
			return false
		},
	)
}
