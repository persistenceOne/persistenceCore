package cmd

import (
	"fmt"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum"
	"github.com/persistenceOne/persistenceCore/pStake/status"
	"github.com/persistenceOne/persistenceCore/pStake/tendermint"
	"github.com/spf13/cobra"
)

func GetCmd(initClientCtx client.Context) *cobra.Command {
	pStakeCommand := &cobra.Command{
		Use:   "pStake [path_to_chain_json] [mnemonics]",
		Short: "Start pStake",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			timeout, err := cmd.Flags().GetString(constants.FlagTimeOut)
			if err != nil {
				log.Fatalln(err)
			}

			coinType, err := cmd.Flags().GetUint32(constants.FlagCoinType)
			if err != nil {
				log.Fatalln(err)
			}

			denom, err := cmd.Flags().GetString(constants.FlagDenom)
			if err != nil {
				log.Fatalln(err)
			}
			constants.Denom = denom

			homePath, err := cmd.Flags().GetString(constants.FlagPStakeHome)
			if err != nil {
				log.Fatalln(err)
			}

			tmSleepTime, err := cmd.Flags().GetInt(constants.FlagTendermintSleepTime)
			if err != nil {
				log.Fatalln(err)
			}
			tmSleepDuration := time.Duration(tmSleepTime) * time.Millisecond

			ethereumEndPoint, err := cmd.Flags().GetString(constants.FlagEthereumEndPoint)
			if err != nil {
				log.Fatalln(err)
			}

			ethSleepTime, err := cmd.Flags().GetInt(constants.FlagEthereumSleepTime)
			if err != nil {
				log.Fatalln(err)
			}
			ethSleepDuration := time.Duration(ethSleepTime) * time.Millisecond

			ports, err := cmd.Flags().GetString("ports")
			if err != nil {
				log.Fatalln(err)
			}

			kafkaHome, err := cmd.Flags().GetString(kafka.FlagKafkaHome)
			if err != nil {
				log.Fatalln(err)
			}

			tmStart, err := cmd.Flags().GetInt64(constants.FlagTendermintStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			ethStart, err := cmd.Flags().GetInt64(constants.FlagEthereumStartHeight)
			if err != nil {
				log.Fatalln(err)
			}

			ethPrivateKey, err := cmd.Flags().GetString(constants.FlagEthPrivateKey)
			if err != nil {
				log.Fatalln(err)
			}
			constants.EthAccountPrivateKey, err = crypto.HexToECDSA(ethPrivateKey)
			if err != nil {
				log.Fatal(err)
			}

			ethGasLimit, err := cmd.Flags().GetUint64(constants.FlagEthGasLimit)
			if err != nil {
				log.Fatalln(err)
			}

			db, err := status.InitializeDB(homePath+"/db", tmStart, ethStart)
			if err != nil {
				log.Fatalln(err)
			}
			defer db.Close()

			chain, err := tendermint.InitializeAndStartChain(args[0], timeout, homePath, coinType, args[1])
			if err != nil {
				log.Fatalln(err)
			}
			constants.Address = chain.MustGetAddress()

			ethereumClient, err := ethclient.Dial(ethereumEndPoint)
			if err != nil {
				log.Fatalf("Error while dialing to eth node %s: %s\n", ethereumEndPoint, err.Error())
			}

			protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)
			portsList := strings.Split(ports, ",")
			kafkaState := kafka.NewKafkaState(portsList, kafkaHome)
			go kafkaRoutine(kafkaState, protoCodec)
			server.TrapSignal(kafkaClose(kafkaState))

			log.Println("Starting to listen ethereum....")
			go ethereum.StartListening(ethereumClient, ethSleepDuration, kafkaState, protoCodec)

			//TODO Example: Remove this later
			ethAddress := goEthCommon.HexToAddress("0xac749a63F87Fe0A978Cb1002c2DFe9fdC5Bd52e4")
			ethTxMsg := ethereum.EthTxMsg{
				Address: ethAddress,
				Amount:  big.NewInt(100),
			}
			txhash, err := ethereum.SendTxToEth(ethereumClient, ethTxMsg, ethGasLimit)
			if err != nil {
				log.Fatalf("Error while sending eth txs. %s\n", err.Error())
			} else {
				fmt.Println("ETH TX HASH " + txhash)
			}

			log.Println("Starting to listen tendermint....")
			tendermint.StartListening(initClientCtx.WithHomeDir(homePath), chain, kafkaState, protoCodec, tmSleepDuration)

			return nil
		},
	}
	pStakeCommand.Flags().String(constants.FlagTimeOut, constants.DefaultTimeout, "timeout time for connecting to rpc")
	pStakeCommand.Flags().Uint32(constants.FlagCoinType, constants.DefaultCoinType, "coin type for wallet")
	pStakeCommand.Flags().String(constants.FlagPStakeHome, constants.DefaultPStakeHome, "home for pStake")
	pStakeCommand.Flags().String(constants.FlagEthereumEndPoint, constants.DefaultEthereumEndPoint, "ethereum node to connect")
	pStakeCommand.Flags().String("ports", "localhost:9092", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	pStakeCommand.Flags().String(kafka.FlagKafkaHome, kafka.DefaultKafkaHome, "The kafka config file directory")
	pStakeCommand.Flags().Int(constants.FlagTendermintSleepTime, constants.DefaultTendermintSleepTime, "sleep time between block checking for tendermint in ms")
	pStakeCommand.Flags().Int(constants.FlagEthereumSleepTime, constants.DefaultEthereumSleepTime, "sleep time between block checking for ethereum in ms")
	pStakeCommand.Flags().Int64(constants.FlagTendermintStartHeight, constants.DefaultTendermintStartHeight, fmt.Sprintf("Start checking height on tendermint chain from this height (default %d - starts from where last left)", constants.DefaultTendermintStartHeight))
	pStakeCommand.Flags().Int64(constants.FlagEthereumStartHeight, constants.DefaultEthereumStartHeight, fmt.Sprintf("Start checking height on ethereum chain from this height (default %d - starts from where last left)", constants.DefaultEthereumStartHeight))
	pStakeCommand.Flags().String(constants.FlagDenom, constants.DefaultDenom, "denom name")
	pStakeCommand.Flags().String(constants.FlagEthPrivateKey, "", "private keys of ethereum account which does txs.")
	pStakeCommand.Flags().Uint64(constants.FlagEthGasLimit, constants.DefaultEthGasLimit, "Gas limit for eth txs")
	return pStakeCommand
}
