package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers the necessary x/liquidstake interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgLiquidStake{}, "liquidstake/MsgLiquidStake", nil)
	cdc.RegisterConcrete(&MsgStakeToLP{}, "liquidstake/MsgStakeToLP", nil)
	cdc.RegisterConcrete(&MsgLiquidUnstake{}, "liquidstake/MsgLiquidUnstake", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "liquidstake/MsgUpdateParams", nil)
	cdc.RegisterConcrete(&MsgUpdateWhitelistedValidators{}, "liquidstake/MsgUpdateWhitelistedValidators", nil)
	cdc.RegisterConcrete(&MsgSetModulePaused{}, "liquidstake/MsgSetModulePaused", nil)
}

// RegisterInterfaces registers the x/liquidstake interfaces types with the interface registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgLiquidStake{},
		&MsgStakeToLP{},
		&MsgLiquidUnstake{},
		&MsgUpdateParams{},
		&MsgUpdateWhitelistedValidators{},
		&MsgSetModulePaused{},
	)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/liquidstake module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding as Amino
	// is still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/liquidstake and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
