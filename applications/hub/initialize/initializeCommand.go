package initialize

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/commitHub/commitBlockchain/applications/hub"
)

const (
	flagOverwrite    = "overwrite"
	flagClientHome   = "home-client"
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

type printInfo struct {
	Moniker    string          `json:"moniker"`
	ChainID    string          `json:"chain_id"`
	NodeID     string          `json:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message"`
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s\n", string(out))
	return nil
}

func InitializeCommand(ctx *server.Context, cdc *codec.Codec) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:   "initialize [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			configuration := ctx.Config
			configuration.SetRoot(viper.GetString(cli.HomeFlag))

			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			nodeID, _, err := initializeNodeValidatorFiles(configuration)
			if err != nil {
				return err
			}

			configuration.Moniker = args[0]

			var appState json.RawMessage
			genFile := configuration.GenesisFile()

			if appState, err = initializeEmptyGenesis(cdc, genFile, viper.GetBool(flagOverwrite)); err != nil {
				return err
			}

			if err = exportGenesisFile(genFile, chainID, nil, appState); err != nil {
				return err
			}

			toPrint := printInfo{configuration.Moniker, chainID, nodeID, "", appState}

			cfg.WriteConfigFile(filepath.Join(configuration.RootDir, "configuration", "configuration.toml"), configuration)

			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, hub.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}
