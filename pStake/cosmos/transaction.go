package cosmos

import (
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/golang/protobuf/proto"
	"github.com/persistenceOne/persistenceCore/kafka"
	tmTypes "github.com/tendermint/tendermint/types"
)

func HandleTxEvent(clientCtx *client.Context, txEvent tmTypes.EventDataTx, kafkaState kafka.KafkaState) {
	if txEvent.Result.Code == 0 {
		handleEncodeTx(clientCtx, txEvent.Tx, kafkaState)
	}
}

// handleEncodeTx Should be called if tx is known to be successful
func handleEncodeTx(clientCtx *client.Context, encodedTx []byte, kafkaState kafka.KafkaState) {
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

	fmt.Printf("Memo: %s\n", tx.GetMemo())

	for _, msg := range tx.GetMsgs() {
		switch txMsg := msg.(type) {
		case *banktypes.MsgSend:
			//Convert txMsg to the Msg we want to send forward
			msgBytes, err := proto.Marshal(sdk.Msg(txMsg))
			if err != nil {
				panic(err)
			}
			if true {
				err = kafka.ProducerDeliverMessage(msgBytes, kafka.ToEth, kafkaState.Producer)
				if err != nil {
					log.Print("Failed to add msg to kafka queue: ", err)
				}
				log.Printf("Produced to kafka: %v, for topic %v ", msg.String(), kafka.ToEth)
			} else {
				// reversal queue
				err = kafka.ProducerDeliverMessage(msgBytes, kafka.ToTendermint, kafkaState.Producer)
				if err != nil {
					log.Print("Failed to add msg to kafka queue: ", err)
				}
				log.Printf("Produced to kafka: %v, for topic %v ", msg.String(), kafka.ToTendermint)
			}
		default:

		}
	}
}
