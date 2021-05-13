package cosmos

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
	tmTypes "github.com/tendermint/tendermint/types"
	"log"
)

func handleTxEvent(chain *relayer.Chain, clientCtx client.Context, txEvent tmTypes.EventDataTx) {
	if txEvent.Result.Code == 0 {
		handleEncodeTx(clientCtx, txEvent.Tx)
	}
}

func handleEncodeTx(clientCtx client.Context, encodedTx []byte) {
	// Should be used if encodedTx is string
	//decodedTx, err := base64.StdEncoding.DecodeString(encodedTx)
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}

	txInterface, err := clientCtx.TxConfig.TxDecoder()(encodedTx)
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
