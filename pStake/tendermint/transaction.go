package tendermint

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceCore/kafka/utils"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum"
	"log"
	"math/big"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

func handleTxEvent(clientCtx client.Context, txEvent tmTypes.EventDataTx, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) {
	if txEvent.Result.Code == 0 {
		_ = handleEncodeTx(clientCtx, txEvent.Tx, kafkaState, protoCodec)
	}
}

func handleTxSearchResult(clientCtx client.Context, txSearchResult *tmCoreTypes.ResultTxSearch, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
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
func handleEncodeTx(clientCtx client.Context, encodedTx []byte, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
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

	memo := strings.TrimSpace(tx.GetMemo())
	validMemo := goEthCommon.IsHexAddress(memo)
	var ethAddress goEthCommon.Address
	if validMemo {
		ethAddress = goEthCommon.HexToAddress(memo)
	}

	for _, msg := range tx.GetMsgs() {
		switch txMsg := msg.(type) {
		case *banktypes.MsgSend:
			var amount *big.Int
			for _, coin := range txMsg.Amount {
				if coin.Denom == constants.Denom {
					amount = coin.Amount.BigInt()
					break
				}
			}
			sendToEth := txMsg.ToAddress == constants.Address.String() && amount != nil && validMemo
			if sendToEth {
				ethTxMsg := ethereum.EthTxMsg{
					Address: ethAddress,
					Amount:  amount,
				}
				msgBytes, err := json.Marshal(ethTxMsg)
				if err != nil {
					panic(err)
				}
				err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, kafkaState.Producer)
				if err != nil {
					log.Print("Failed to add msg to kafka queue: ", err)
				}
				log.Printf("Produced to kafka: %v, for topic %v ", msg.String(), utils.ToEth)
			} else {
				msg := &banktypes.MsgSend{
					FromAddress: txMsg.ToAddress,
					ToAddress:   txMsg.FromAddress,
					Amount:      txMsg.Amount,
				}
				msgBytes, err := protoCodec.MarshalInterface(sdk.Msg(msg))
				if err != nil {
					panic(err)
				}
				err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, kafkaState.Producer)
				if err != nil {
					log.Print("Failed to add msg to kafka queue: ", err)
				}
				log.Printf("Produced to kafka: %v, for topic %v ", msg.String(), utils.ToTendermint)
			}
		default:

		}
	}

	return nil
}
