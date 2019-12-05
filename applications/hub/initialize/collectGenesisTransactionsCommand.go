package initialize

import (
	"encoding/json"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/commitHub/commitBlockchain/applications/hub"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const (
	flagGenesisTransactionDirectory = "genesisTransactionDirectory"
)

type initialConfiguration struct {
	ChainID   string
	GenTxsDir string
	Name      string
	NodeID    string
	ValPubKey crypto.PubKey
}

func CollectGenesisTransactionsCommand(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-genesis-transactions",
		Short: "Collect genesis transactions and output a genesis.json file",
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			name := viper.GetString(client.FlagName)
			nodeID, valPubKey, err := InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			genDoc, err := LoadGenesisDoc(cdc, config.GenesisFile())
			if err != nil {
				return err
			}

			genTxsDir := viper.GetString(flagGenesisTransactionDirectory)
			if genTxsDir == "" {
				genTxsDir = filepath.Join(config.RootDir, "config", "gentx")
			}

			toPrint := printInfo{config.Moniker, genDoc.ChainID, nodeID, genTxsDir, json.RawMessage("")}
			initCfg := initialConfiguration{genDoc.ChainID, genTxsDir, name, nodeID, valPubKey}

			appMessage, err := generateApplicationStateFromInitialConfiguration(cdc, config, initCfg, genDoc)
			if err != nil {
				return err
			}

			toPrint.AppMessage = appMessage

			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, hub.application.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagGenesisTransactionDirectory, "",
		"override default \"gentx\" directory from which collect and execute "+
			"genesis transactions; default [--home]/config/gentx/")
	return cmd
}

func generateApplicationStateFromInitialConfiguration(cdc *codec.Codec, config *cfg.Config, initConfiguration initialConfiguration, genDoc types.GenesisDoc) (appState json.RawMessage, err error) {

	genFile := config.GenesisFile()
	var (
		appGenTxs       []auth.StdTx
		persistentPeers string
		genTxs          []json.RawMessage
		jsonRawTx       json.RawMessage
	)

	appGenTxs, persistentPeers, err = hub.application.CollectStdTxs(
		cdc, config.Moniker, initConfiguration.GenTxsDir, genDoc,
	)
	if err != nil {
		return
	}

	genTxs = make([]json.RawMessage, len(appGenTxs))
	config.P2P.PersistentPeers = persistentPeers

	for i, stdTx := range appGenTxs {
		jsonRawTx, err = cdc.MarshalJSON(stdTx)
		if err != nil {
			return
		}
		genTxs[i] = jsonRawTx
	}

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	appState, err = hub.application.GaiaAppGenStateJSON(cdc, genDoc, genTxs)
	if err != nil {
		return
	}

	err = ExportGenesisFile(genFile, initConfiguration.ChainID, nil, appState)
	return
}
