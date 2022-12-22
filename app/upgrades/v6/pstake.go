package v6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	lscosmoskeeper "github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
)

// MintPstakeTokens send stk/uatom to persistence1zl42hd5h9c7z4ej43fhss9nvgm6nuad0js8z6n (for https://www.mintscan.io/cosmos/txs/DE691EC8EBB5A79E2AB421291660111E893823CA0CC9EBDED5E3C72B503067C3 sending amount to reward address)
// at c_value 0.999142233051758540 = 44961400stk/uatom (`curl -X GET -H "Content-Type: application/json" -H "x-cosmos-block-height: 8650000" 'https://rest.core.persistence.one/pstake/lscosmos/v1beta1/c_value'`)
func MintPstakeTokens(ctx sdk.Context, k *lscosmoskeeper.Keeper) error {
	if ctx.ChainID() != "core-1" && ctx.ChainID() != "test-core-1" {
		return nil
	}

	atomTVU := k.GetDepositAccountAmount(ctx).
		Add(k.GetIBCTransferTransientAmount(ctx)).
		Add(k.GetDelegationTransientAmount(ctx)).
		Add(k.GetStakedAmount(ctx)).
		Add(k.GetHostDelegationAccountAmount(ctx))

	mintedAmount := k.GetMintedAmount(ctx)
	mintDenom := k.GetHostChainParams(ctx).MintDenom
	if atomTVU.LTE(mintedAmount) {
		return nil
	}

	toNewMint := atomTVU.Sub(mintedAmount)

	if ctx.ChainID() == "core-1" {
		if toNewMint.GT(sdk.NewInt(44961400)) {
			mischiefUserAddress := sdk.MustAccAddressFromBech32("persistence1zl42hd5h9c7z4ej43fhss9nvgm6nuad0js8z6n")
			toSendUser := sdk.NewInt(44961400)
			err := k.MintTokens(ctx, sdk.NewCoin(mintDenom, toSendUser), mischiefUserAddress)
			if err != nil {
				k.Logger(ctx).Error("Failed to mint and send 44961400stk/uatom to persistence1zl42hd5h9c7z4ej43fhss9nvgm6nuad0js8z6n")
				return err
			}
			pstakeFeeAddress := sdk.MustAccAddressFromBech32(k.GetHostChainParams(ctx).PstakeParams.PstakeFeeAddress)
			remainingAmount := toNewMint.Sub(toSendUser)
			err = k.MintTokens(ctx, sdk.NewCoin(mintDenom, remainingAmount), pstakeFeeAddress)
			if err != nil {
				k.Logger(ctx).Error("Failed to mint and send remainingAmount to pstakeFeeAddress")
				return err
			}
		}
		return nil
	}
	if ctx.ChainID() == "test-core-1" {
		pstakeFeeAddress := sdk.MustAccAddressFromBech32(k.GetHostChainParams(ctx).PstakeParams.PstakeFeeAddress)
		err := k.MintTokens(ctx, sdk.NewCoin(mintDenom, toNewMint), pstakeFeeAddress)
		if err != nil {
			k.Logger(ctx).Error("Failed to mint and send toNewMint to pstakeFeeAddress")
			return err
		}
		return nil
	}

	return nil
}
