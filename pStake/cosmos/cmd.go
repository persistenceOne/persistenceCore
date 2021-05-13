package cosmos

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum"
	"github.com/spf13/cobra"
	tmTypes "github.com/tendermint/tendermint/types"
	"log"
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

			homePath, err := cmd.Flags().GetString(constants.FlagHome)
			if err != nil {
				log.Fatalln(err)
			}

			ethereumEndPoint, err := cmd.Flags().GetString(constants.FlagEthereumEndPoint)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("Starting to listen ethereum....")
			go ethereum.StartListening(ethereumEndPoint)

			run(initClientCtx, args[0], timeout, homePath, coinType, args[1])

			return nil
		},
	}
	pStakeCommand.Flags().String(constants.FlagTimeOut, "10s", "timeout time for connecting to rpc")
	pStakeCommand.Flags().Uint32(constants.FlagCoinType, 118, "coin type for wallet")
	pStakeCommand.Flags().String(constants.FlagHome, "./pStake", "home for pStake")
	pStakeCommand.Flags().String(constants.FlagEthereumEndPoint, "wss://goerli.infura.io/ws/v3/e2549c9ec9764e46a7768cc7619a1939", "ethereum node to connect")
	return pStakeCommand
}

func run(initClientCtx client.Context, chainConfigJsonPath, timeout, homePath string, coinType uint32, mnemonics string) {
	err := initializeAndStartChain(chainConfigJsonPath, timeout, homePath, coinType, mnemonics)
	if err != nil {
		log.Fatalf("Error while intiializing and starting chain: %s\n", err.Error())
	}

	log.Println("Starting to listen cosmos txs....")
	txEvents, txCancel, err := Chain.Subscribe(constants.TxEvents)
	if err != nil {
		log.Fatalf("Error while subscribing to tx events: %s\n", err.Error())
	}
	defer txCancel()

	log.Println("Starting to listen cosmos blocks....")
	blockEvents, blockCancel, err := Chain.Subscribe(constants.BlockEvents)
	if err != nil {
		log.Fatalf("Error while subscribing to block events: %s\n", err.Error())
	}
	defer blockCancel()

	for {
		select {
		case txEvent := <-txEvents:
			go handleTxEvent(Chain, initClientCtx, txEvent.Data.(tmTypes.EventDataTx))
		case blockEvent := <-blockEvents:
			go handleNewBlock(Chain, blockEvent.Data.(tmTypes.EventDataNewBlock))
		}
	}

}
