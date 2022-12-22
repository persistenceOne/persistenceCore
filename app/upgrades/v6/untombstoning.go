package v6

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type CosMints struct {
	Address     string `json:"address"`
	AmountUxprt string `json:"amount"`
	Delegatee   string `json:"validator_address"`
}

type Validator struct {
	name       string
	valAddress string
	conAddress string
}

// Create new Validator vars for each validator that needs to be untombstoned
var (
	mainnetVals = []Validator{
		{"HashQuark", "persistencevaloper1gydvxcnm95zwdz7h7whpmusy5d5c3ck0p9muc9", "persistencevalcons1dmjc55ve2pe537hu8h8rjrjhp4r536g5jlnlk8"},
		{"fox99", "persistencevaloper1y2svn2zvc0puv3rx6w39aa4zlgj7qe0fz8sh6x", "persistencevalcons1ak5f5ywzmersz4z7e3nsqkem4uvf5jyya62w3c"},
		{"Smart Stake", "persistencevaloper1qtggtsmexluvzulehxs7ypsfl82yk5aznrr2zd", "persistencevalcons1gnevun33uphh9cwkyzau5mcf0fxvuw6cyrf29g"},
		{"Stakin", "persistencevaloper1xykmyvzk88qrlqh3wuw4jckewleyygupsumyj5", "persistencevalcons15fxjrujvsc0le9udjf63504sd4lndcam8ep4cs"},
		{"KuCoin", "persistencevaloper18qgr8va65a50sdmp2yuy4y8w9p4pa2rf76mvmm", "persistencevalcons1m83jqu6q6aqcshnq0yjrdra9nj8rgz79mndh3j"},
	}
	// testnetVals holds the validators to untombstone
	testnetVals = []Validator{
		{"TombRaider", "persistencevaloper1mgd6a660ysram7a0m8ytmjvryneywgm8mg7lcs", "persistence1mgd6a660ysram7a0m8ytmjvryneywgm8jv7z3f"},
	}
)

func mintLostTokens(
	ctx sdk.Context,
	bankKeeper *bankkeeper.BaseKeeper,
	stakingKeeper *stakingkeeper.Keeper,
	mintKeeper *mintkeeper.Keeper,
	cosMints []CosMints,
) error {
	for _, mintRecord := range cosMints {
		cosValAddress, err := sdk.ValAddressFromBech32(mintRecord.Delegatee)
		if err != nil {
			return fmt.Errorf("validator address is not valid bech32: %s", cosValAddress)
		}

		coinAmount, _ := sdk.NewIntFromString(mintRecord.AmountUxprt)

		coin := sdk.NewCoin("uxprt", coinAmount)
		coins := sdk.NewCoins(coin)

		err = mintKeeper.MintCoins(ctx, coins)
		if err != nil {
			return fmt.Errorf("error minting %duxprt to %s: %+v", mintRecord.AmountUxprt, mintRecord.Address, err)
		}

		delegatorAccount, err := sdk.AccAddressFromBech32(mintRecord.Address)
		if err != nil {
			return fmt.Errorf("error converting human address %s to sdk.AccAddress: %+v", mintRecord.Address, err)
		}

		err = bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delegatorAccount, coins)
		if err != nil {
			return fmt.Errorf("error sending minted %duxprt to %s: %+v", mintRecord.AmountUxprt, mintRecord.Address, err)
		}

		sdkAddress, err := sdk.AccAddressFromBech32(mintRecord.Address)
		if err != nil {
			return fmt.Errorf("account address is not valid bech32: %s", mintRecord.Address)
		}

		cosValidator, found := stakingKeeper.GetValidator(ctx, cosValAddress)
		if !found {
			return fmt.Errorf("cos validator '%s' not found", cosValAddress)
		}

		_, err = stakingKeeper.Delegate(ctx, sdkAddress, coin.Amount, stakingtypes.Unbonded, cosValidator, true)
		if err != nil {
			return fmt.Errorf("error delegating minted %duxprt from %s to %s: %+v", mintRecord.AmountUxprt, mintRecord.Address, mintRecord.Delegatee, err)
		}
	}
	return nil
}

func revertTombstone(ctx sdk.Context, slashingKeeper *slashingkeeper.Keeper, validator Validator) error {
	cosValAddress, err := sdk.ValAddressFromBech32(validator.valAddress)
	if err != nil {
		return fmt.Errorf("validator address is not valid bech32: %s", cosValAddress)
	}

	cosConsAddress, err := sdk.ConsAddressFromBech32(validator.conAddress)
	if err != nil {
		return fmt.Errorf("consensus address is not valid bech32: %s", cosValAddress)
	}

	// Revert Tombstone info
	signInfo, ok := slashingKeeper.GetValidatorSigningInfo(ctx, cosConsAddress)

	if !ok {
		return fmt.Errorf("cannot tombstone validator that does not have any signing information: %s", cosConsAddress.String())
	}
	if !signInfo.Tombstoned {
		return fmt.Errorf("cannut untombstone a validator that is not tombstoned: %s", cosConsAddress.String())
	}

	signInfo.Tombstoned = false
	slashingKeeper.SetValidatorSigningInfo(ctx, cosConsAddress, signInfo)

	// Set jail until=now, the validator then must unjail manually
	slashingKeeper.JailUntil(ctx, cosConsAddress, ctx.BlockTime())
	ctx.Logger().Info(fmt.Sprintf("Tombstone successfully reverted for validator: %s: %s", validator.name, validator.valAddress))

	return nil
}

func RevertCosTombstoning(
	ctx sdk.Context,
	slashingKeeper *slashingkeeper.Keeper,
	mintKeeper *mintkeeper.Keeper,
	bankKeeper *bankkeeper.BaseKeeper,
	stakingKeeper *stakingkeeper.Keeper,
) error {
	// Run code on mainnet and testnet for minting lost tokens
	// check the blockheight is more than tombstoning height
	if ctx.ChainID() == "core-1" || ctx.ChainID() == "test-core-1" {
		var Mints []CosMints
		var vals []Validator
		if ctx.ChainID() == "core-1" || ctx.BlockHeight() > 88647536 {
			var cosMints []CosMints
			err := json.Unmarshal([]byte(recordsJsonString), &cosMints)
			if err != nil {
				return fmt.Errorf("error reading COS JSON: %+v", err)
			}
			Mints = append(Mints, cosMints...)
			vals = append(vals, mainnetVals...)
		}
		if ctx.ChainID() == "test-core-1" || ctx.BlockHeight() > 8647536 {
			var cosMints []CosMints
			err := json.Unmarshal([]byte(testnetRecordsJsonString), &cosMints)
			if err != nil {
				return fmt.Errorf("error reading COS JSON: %+v", err)
			}
			Mints = append(Mints, cosMints...)
			vals = append(vals, testnetVals...)
		}

		for _, value := range vals {
			revertTombstone(ctx, slashingKeeper, value)
		}

		mintLostTokens(ctx, bankKeeper, stakingKeeper, mintKeeper, Mints)
	}
	return nil
}
