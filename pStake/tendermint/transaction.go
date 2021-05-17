package tendermint

import (
	"log"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceCore/kafka"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

func handleTxEvent(clientCtx client.Context, txEvent tmTypes.EventDataTx, kafkaState kafka.KafkaState, protoCodec *codec.ProtoCodec) {
	if txEvent.Result.Code == 0 {
		_ = handleEncodeTx(clientCtx, txEvent.Tx, kafkaState, protoCodec)
	}
}

func handleTxSearchResult(clientCtx client.Context, txSearchResult *tmCoreTypes.ResultTxSearch, kafkaState kafka.KafkaState, protoCodec *codec.ProtoCodec) error {
	for _, tx := range txSearchResult.Txs {
		if tx.TxResult.Code == 0 {
			err := handleEncodeTx(clientCtx, tx.Tx, kafkaState, protoCodec)
			if err != nil {
				log.Printf("Failed to process tendermint tx: %s\n", tx.Hash)
				return err
			}
		}
	}
	return nil
}

// handleEncodeTx Should be called if tx is known to be successful
func handleEncodeTx(clientCtx client.Context, encodedTx []byte, kafkaState kafka.KafkaState, protoCodec *codec.ProtoCodec) error {
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

	validMemo := goEthCommon.IsHexAddress(strings.TrimSpace(tx.GetMemo()))

	for _, msg := range tx.GetMsgs() {
		switch txMsg := msg.(type) {
		case *banktypes.MsgSend:
			//TODO Convert txMsg to the Msg we want to send forward
			msgBytes, err := protoCodec.MarshalInterface(sdk.Msg(txMsg))
			if err != nil {
				panic(err)
			}
			if txMsg.ToAddress == Chain.MustGetAddress().String() {
				if validMemo {
					err = kafka.ProducerDeliverMessage(msgBytes, kafka.ToEth, kafkaState.Producer)
					if err != nil {
						log.Print("Failed to add msg to kafka queue: ", err)
					}
					log.Printf("Produced to kafka: %v, for topic %v ", msg.String(), kafka.ToEth)
				} else {
					//TODO Convert txMsg to the Msg we want to sent to tendermint reversal queue
					err = kafka.ProducerDeliverMessage(msgBytes, kafka.ToTendermint, kafkaState.Producer)
					if err != nil {
						log.Print("Failed to add msg to kafka queue: ", err)
					}
					log.Printf("Produced to kafka: %v, for topic %v ", msg.String(), kafka.ToTendermint)
				}
			}
		default:

		}
	}

	return nil
}
