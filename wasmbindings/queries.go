package wasmbindings

import oraclekeeper "github.com/persistenceOne/persistence-sdk/v2/x/oracle/keeper"

type QueryPlugin struct {
	oracleKeeper *oraclekeeper.Keeper
}

func NewQueryPlugin(oracleKeeper *oraclekeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{oracleKeeper: oracleKeeper}
}
