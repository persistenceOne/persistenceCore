package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/spf13/cobra"
)

func GetCmd(initClientCtx client.Context) *cobra.Command {
	pStakeCommand := &cobra.Command{
		Use:   "pStake",
		Short: "Persistence Hub Node Daemon (server)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			run(initClientCtx)

			return nil
		},
	}

	return pStakeCommand
}
func run(initClientCtx client.Context) {

	//height := "15"
	//response, err := queries.GetTxsByHeight(height)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	Tx := "CpABCo0BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm0KLWNvc21vczFkMjBnMGdjd2hyd3Y4ZjI2MjZkeDBua2hhdXUwcnNxazhwZm4wbRItY29zbW9zMTcyN2ZxNng3NmVzcXowM3F0MHZleHM2ajhzZG01NHF4N2txNjNwGg0KBXN0YWtlEgQxMDAwElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQPvtx3smWdI9li0uGECYv8di+QrUI3yM/+4JGOPn1x/dhIECgIIARgBEgQQwJoMGkAimAko9tTisNOSSD5OJ5hDLgT614QmkVwI1LxIgo+y0mrVsGuR/58kQLonSDAOtTBkcnAnZwxsLQAbMQnABSXV"

	fmt.Println(Tx)

	decodedTx, err := base64.StdEncoding.DecodeString(Tx) //response.Result.Txs[0].Tx)
	if err != nil {
		fmt.Println(err.Error())
	}

	txInterface, err := initClientCtx.TxConfig.TxDecoder()(decodedTx)
	if err != nil {
		fmt.Println(err.Error())
	}

	tx, ok := txInterface.(signing.Tx)
	if !ok {
		fmt.Println("Unable to parse tx")
	}
	fmt.Println(tx.GetMsgs()[0].Type())
	//sendMsg := tx.GetMsgs()[0].(banktypes.MsgSend)
	//fmt.Println(sendMsg.FromAddress)
}
