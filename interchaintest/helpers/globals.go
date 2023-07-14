package helpers

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var SomeoneAddress sdk.AccAddress

func init() {
	addr, err := sdk.AccAddressFromHexUnsafe("0000000000000000000000000000000000000000")
	if err != nil {
		panic("failed to init address")
	}

	SomeoneAddress = addr
}

const (
	PersistenceCore  = "persistenceCore"
	Bech32MainPrefix = "persistence"

	PersistenceBondDenom = "uxprt"
	PersistenceCoinType  = 118
	PersistencePurpose   = 44

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

// SetConfig params at the package state
func SetConfig() {
	cfg := sdk.GetConfig()

	cfg.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	cfg.SetCoinType(PersistenceCoinType)
	cfg.SetPurpose(PersistencePurpose)
}

func debugOutput(t *testing.T, stdout string) {
	if len(stdout) == 0 {
		return
	}

	if true {
		t.Log(stdout)
	}
}
