package main

import (
	"encoding/json"
	"io"

	"github.com/commitHub/commitBlockchain/applications/hub"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tendermintABSITypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tendermintTypes "github.com/tendermint/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/applications/hub"
	"github.com/commitHub/commitBlockchain/applications/hub/initialize"
)

const flagInvalidCheckPeriod = "invalid-check-period"

var invalidCheckPeriod uint

func main() {
	codec := hub.MakeCodec()

	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(sdkTypes.Bech32PrefixAccAddr, sdkTypes.Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(sdkTypes.Bech32PrefixValAddr, sdkTypes.Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(sdkTypes.Bech32PrefixConsAddr, sdkTypes.Bech32PrefixConsPub)
	configuration.Seal()

	context := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCommand := &cobra.Command{
		Use:               "hubNode",
		Short:             "Commit Hub Node Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(context),
	}

	rootCommand.AddCommand(initialize.InitializeCommand(context, codec))
	rootCommand.AddCommand(initialize.CollectGenesisTransactionsCommand(context, codec))
	rootCommand.AddCommand(initialize.TestnetCommand(context, codec))
	rootCommand.AddCommand(initialize.GenesisTransactionCommand(context, codec))
	rootCommand.AddCommand(initialize.AddGenesisAccountCommand(context, codec))
	rootCommand.AddCommand(initialize.ValidateGenesisCommand(context, codec))
	rootCommand.AddCommand(client.NewCompletionCmd(rootCommand, true))

	server.AddCommands(context, codec, rootCommand, newApplication, exportApplicationStateAndValidators)

	executor := cli.PrepareBaseCmd(rootCommand, "CA", hub.application.DefaultNodeHome)
	rootCommand.PersistentFlags().UintVar(
		&invalidCheckPeriod,
		flagInvalidCheckPeriod,
		0,
		"Assert registered invariants every N blocks",
	)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApplication(logger log.Logger, db tendermintDB.DB, traceStore io.Writer) tendermintABSITypes.Application {
	return hub.NewCommitHubApplication(logger, db, traceStore, true, invalidCheckPeriod, basehub.application.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))), basehub.application.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)))
}

func exportApplicationStateAndValidators(logger log.Logger, db tendermintDB.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string) (json.RawMessage, []tendermintTypes.GenesisValidator, error) {

	if height != -1 {
		genesisApplication := hub.NewCommitHubApplication(logger, db, traceStore, false, uint(1))
		err := genesisApplication.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return genesisApplication.ExportApplicationStateAndValidators(forZeroHeight, jailWhiteList)
	} else {
		genesisApplication := hub.NewCommitHubApplication(logger, db, traceStore, true, uint(1))
		return genesisApplication.ExportApplicationStateAndValidators(forZeroHeight, jailWhiteList)
	}
}
