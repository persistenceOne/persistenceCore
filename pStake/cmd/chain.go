package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum"
	"github.com/persistenceOne/persistenceCore/pStake/status"
	"github.com/persistenceOne/persistenceCore/pStake/tendermint"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
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

			protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)
			portsList := strings.Split(ports, ",")
			kafkaState := kafka.NewKafkaState(portsList, kafkaHome)
			go kafkaRoutine(kafkaState, protoCodec)
			server.TrapSignal(kafkaClose(kafkaState))

			log.Println("Starting to listen ethereum....")
			go ethereum.StartListening(ethereumEndPoint, ethSleepDuration, kafkaState, protoCodec)
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
	return pStakeCommand
}
