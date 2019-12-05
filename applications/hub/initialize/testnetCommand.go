package initialize

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/keys"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/commitHub/commitBlockchain/applications/hub"
	"github.com/cosmos/cosmos-sdk/codec"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/server"
)

var (
	flagNodeDirectoryPrefix = "node-dir-prefix"
	flagNumberOfValidators  = "v"
	flagOutputDirectory     = "output-dir"
	flagNodeDaemonHome      = "node-daemon-home"
	flagNodeClientHome      = "node-cli-home"
	flagStartingIPAddress   = "starting-ip-address"
)

const nodeDirPerm = 0755

func TestnetCommand(ctx *server.Context, cdc *codec.Codec) *cobra.Command {

	command := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a Commit testnet",
		Long:  `testnet will create "v" number of directories and populate each with necessary files (private validator, genesis, config, etc.). Note, strict routability for addresses is turned off in the config file. Example: gaiad testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			return initializeTestnet(config, cdc)
		},
	}

	command.Flags().Int(flagNumberOfValidators, 4,
		"Number of validators to initialize the testnet with",
	)
	command.Flags().StringP(flagOutputDirectory, "o", "./mytestnet",
		"Directory to store initialization data for the testnet",
	)
	command.Flags().String(flagNodeDirectoryPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)",
	)
	command.Flags().String(flagNodeDaemonHome, "commitNode",
		"Home directory of the node's daemon configuration",
	)
	command.Flags().String(flagNodeClientHome, "commitClient",
		"Home directory of the node's cli configuration",
	)
	command.Flags().String(flagStartingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")

	command.Flags().String(
		client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created",
	)
	command.Flags().String(
		server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)",
	)

	return command
}

func initializeTestnet(config *tmconfig.Config, cdc *codec.Codec) error {
	var chainID string

	outDir := viper.GetString(flagOutputDirectory)
	numberOfValidators := viper.GetInt(flagNumberOfValidators)

	chainID = viper.GetString(client.FlagChainID)
	if chainID == "" {
		chainID = "chain-" + cmn.RandStr(6)
	}

	monikers := make([]string, numberOfValidators)
	nodeIDs := make([]string, numberOfValidators)
	validatorPubKeys := make([]crypto.PubKey, numberOfValidators)

	chainConfig := srvconfig.DefaultConfig()
	chainConfig.MinGasPrices = viper.GetString(server.FlagMinGasPrices)

	var (
		genesisAccounts []hub.application.GenesisAccount
		genesisFiles    []string
	)

	for i := 0; i < numberOfValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirectoryPrefix), i)
		nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
		nodeCliHomeName := viper.GetString(flagNodeClientHome)
		nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
		clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)
		genesisTransactionDirectory := filepath.Join(outDir, "gentxs")

		config.SetRoot(nodeDir)

		err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		err = os.MkdirAll(clientDir, nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		monikers = append(monikers, nodeDirName)
		config.Moniker = nodeDirName

		ip, err := getIP(i, viper.GetString(flagStartingIPAddress))
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		nodeIDs[i], validatorPubKeys[i], err = InitializeNodeValidatorFiles(config)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genesisFiles = append(genesisFiles, config.GenesisFile())

		buf := client.BufferStdin()
		prompt := fmt.Sprintf(
			"Password for account '%s' (default %s):", nodeDirName, hub.application.DefaultKeyPass,
		)

		keyPass, err := client.GetPassword(prompt, buf)
		if err != nil && keyPass != "" {
			return err
		}

		if keyPass == "" {
			keyPass = hub.application.DefaultKeyPass
		}

		addr, secret, err := server.GenerateSaveCoinKey(clientDir, nodeDirName, keyPass, true)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		err = writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint)
		if err != nil {
			return err
		}

		accTokens := sdk.TokensFromTendermintPower(1000)
		accStakingTokens := sdk.TokensFromTendermintPower(500)
		genesisAccounts = append(genesisAccounts, hub.application.GenesisAccount{
			Address: addr,
			Coins: sdk.Coins{
				sdk.NewCoin(fmt.Sprintf("%stoken", nodeDirName), accTokens),
				sdk.NewCoin(sdk.DefaultBondDenom, accStakingTokens),
			},
		})

		valTokens := sdk.TokensFromTendermintPower(100)
		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			validatorPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			staking.NewDescription(nodeDirName, "", "", ""),
			staking.NewCommissionMsg(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			sdk.OneInt(),
		)
		kb, err := keys.NewKeyBaseFromDir(clientDir)
		if err != nil {
			return err
		}
		tx := auth.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
		txBldr := auth.NewTxBuilderFromCLI().WithChainID(chainID).WithMemo(memo).WithKeybase(kb)

		signedTx, err := txBldr.SignStdTx(nodeDirName, hub.application.DefaultKeyPass, tx, false)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		txBytes, err := cdc.MarshalJSON(signedTx)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		err = writeFile(fmt.Sprintf("%v.json", nodeDirName), genesisTransactionDirectory, txBytes)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		gaiaConfigFilePath := filepath.Join(nodeDir, "config/gaiad.toml")
		srvconfig.WriteConfigFile(gaiaConfigFilePath, chainConfig)
	}

	if err := initGenFiles(cdc, chainID, genesisAccounts, genesisFiles, numberOfValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		cdc, config, chainID, monikers, nodeIDs, validatorPubKeys, numberOfValidators,
		outDir, viper.GetString(flagNodeDirectoryPrefix), viper.GetString(flagNodeDaemonHome),
	)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully initialized %d node directories\n", numberOfValidators)
	return nil
}

func initGenFiles(
	cdc *codec.Codec, chainID string, accs []hub.application.GenesisAccount,
	genFiles []string, numValidators int,
) error {

	appGenState := hub.application.NewDefaultGenesisState()
	appGenState.Accounts = accs

	appGenStateJSON, err := codec.MarshalJSONIndent(cdc, appGenState)
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func collectGenFiles(
	cdc *codec.Codec, config *tmconfig.Config, chainID string,
	monikers, nodeIDs []string, valPubKeys []crypto.PubKey,
	numValidators int, outDir, nodeDirPrefix, nodeDaemonHomeName string,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
		gentxsDir := filepath.Join(outDir, "gentxs")
		moniker := monikers[i]
		config.Moniker = nodeDirName

		config.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := initialConfiguration{chainID, gentxsDir, moniker, nodeID, valPubKey}

		genDoc, err := LoadGenesisDoc(cdc, config.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := generateApplicationStateFromInitialConfiguration(cdc, config, initCfg, genDoc)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := config.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		err = ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime)
		if err != nil {
			return err
		}
	}

	return nil
}

func getIP(i int, startingIPAddr string) (string, error) {
	var (
		ip  string
		err error
	)

	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
	} else {
		ip, err = calculateIP(startingIPAddr, i)
		if err != nil {
			return "", err
		}
	}

	return ip, nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := cmn.EnsureDir(writePath, 0700)
	if err != nil {
		return err
	}

	err = cmn.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}
