package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/helpers"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/spf13/cobra"
	tmservice "github.com/tendermint/tendermint/libs/service"
	tmTypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"log"
	"os"
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

			//account, err := cmd.Flags().GetInt(constants.FlagAccount)
			//if err != nil {
			//	log.Fatalln(err)
			//}
			//
			//index, err := cmd.Flags().GetInt(constants.FlagIndex)
			//if err != nil {
			//	log.Fatalln(err)
			//}

			run(initClientCtx, args[0], timeout, coinType, args[1])
			return nil
		},
	}
	pStakeCommand.Flags().String(constants.FlagTimeOut, "10s", "timeout time for connecting to rpc")
	pStakeCommand.Flags().Uint32(constants.FlagCoinType, 118, "coin type for wallet")
	//pStakeCommand.Flags().Int(constants.FlagAccount, 0, "account no. for wallet")
	//pStakeCommand.Flags().Int(constants.FlagIndex, 0, "index of wallet")
	return pStakeCommand
}

func run(initClientCtx client.Context, chainConfigJsonPath, timeout string, coinType uint32, mnemonics string) {
	chain, err := fileInputAdd(chainConfigJsonPath)
	to, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("Error while parsing timeout: %w", err)
	}
	//homePath, err := os.Getwd()
	//if err != nil {
	//	log.Fatalf("Error while getting current directory: %w", err)
	//}

	homePath := "./pStake"

	err = chain.Init(homePath, to, nil, true)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if chain.KeyExists(chain.Key) {
		log.Printf("deleting old key %s\n", chain.Key)
		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			log.Fatalln("could not delete key %s", chain.Key)
		}
	}

	ko, err := helpers.KeyAddOrRestore(chain, chain.Key, coinType, mnemonics)
	if err != nil {
		log.Fatalf("Error while adding keys: %w", err)
	}

	log.Printf("Keys added: %s", ko.Address)

	if err = chain.Start(); err != nil {
		if err != tmservice.ErrAlreadyStarted {
			chain.Error(err)
			return
		}
	}

	txxEvents, txCancel, err := chain.Subscribe(constants.TxEvents)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer txCancel()

	blockEvents, blockCancel, err := chain.Subscribe(constants.BlockEvents)
	if err != nil {
		chain.Error(err)
		return
	}
	defer blockCancel()

	fromAccount, err := chain.GetAddress()
	if err != nil {
		log.Fatalln(err.Error())
	}
	toAccount, err := sdk.AccAddressFromBech32("cosmos120fgcs32s8wus7k80ysfszwl275x4v87wuuxd9")
	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		select {
		case txEvent := <-txxEvents:
			if txEvent.Data.(tmTypes.EventDataTx).Result.Code == 0 {
				go handleEncodeTx(initClientCtx, txEvent.Data.(tmTypes.EventDataTx).Tx)
			}
		case blockEvent := <-blockEvents:
			fmt.Println(blockEvent.Data.(tmTypes.EventDataNewBlock).Block.Height)

			if blockEvent.Data.(tmTypes.EventDataNewBlock).Block.Height%10 == 0 {
				response, ok, err := chain.SendMsg(banktypes.NewMsgSend(fromAccount, toAccount, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1)))))
				if err != nil {
					log.Println(err.Error())
				}
				if !ok {
					fmt.Println("Transaction %s not ok", response.TxHash)
				}
			}
		}
	}

}

func handleEncodeTx(initClientCtx client.Context, encodedTx []byte) {
	// Should be used if encodedTx is string
	//decodedTx, err := base64.StdEncoding.DecodeString(encodedTx)
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}

	txInterface, err := initClientCtx.TxConfig.TxDecoder()(encodedTx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tx, ok := txInterface.(signing.Tx)
	if !ok {
		log.Fatalln("Unable to parse tx")
	}

	fmt.Printf("Memo: %s", tx.GetMemo())

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *banktypes.MsgSend:
			fmt.Println(msg.String())
		default:

		}
	}
}

func fileInputAdd(file string) (*relayer.Chain, error) {
	// If the user passes in a file, attempt to read the chain config from that file
	c := &relayer.Chain{}
	if _, err := os.Stat(file); err != nil {
		return c, err
	}

	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}

	if err = json.Unmarshal(byt, c); err != nil {
		return c, err
	}

	return c, nil
}
