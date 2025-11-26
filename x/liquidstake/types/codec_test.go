package types_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

func TestRegisterLegacyAminoCodec(t *testing.T) {
	cdc := codec.NewLegacyAmino()

	// Just verify that the function doesn't panic
	require.NotPanics(t, func() {
		types.RegisterLegacyAminoCodec(cdc)
	})
}

func TestRegisterInterfaces(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)

	// Just verify that the function doesn't panic
	require.NotPanics(t, func() {
		types.RegisterInterfaces(registry)
	})
}
