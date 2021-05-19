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
)

//func handleTxEvent(clientCtx client.Context, txEvent tmTypes.EventDataTx, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) {
//	if txEvent.Result.Code == 0 {
//		_ = processTx(clientCtx, txEvent.Tx, kafkaState, protoCodec)
//	}
//}

func handleTxSearchResult(clientCtx client.Context, txSearchResult *tmCoreTypes.ResultTxSearch, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
	for _, tx := range txSearchResult.Txs {
		err := processTx(clientCtx, tx, kafkaState, protoCodec)
		if err != nil {
			log.Printf("Failed to process tendermint tx: %s\n", tx.Hash.String())
			return err
		}
	}
	return nil
}

func processTx(clientCtx client.Context, txQueryResult *tmCoreTypes.ResultTx, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) error {
	if txQueryResult.TxResult.GetCode() == 0 {
		// Should be used if txQueryResult.Tx is string
		//decodedTx, err := base64.StdEncoding.DecodeString(txQueryResult.Tx)
		//if err != nil {
		//	log.Fatalln(err.Error())
		//}

		txInterface, err := clientCtx.TxConfig.TxDecoder()(txQueryResult.Tx)
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

		for i, msg := range tx.GetMsgs() {
			switch txMsg := msg.(type) {
			case *banktypes.MsgSend:
				var amount *big.Int
				for _, coin := range txMsg.Amount {
					if coin.Denom == constants.PSTakeDenom {
						amount = coin.Amount.BigInt()
						break
					}
				}
				if txMsg.ToAddress == constants.PSTakeAddress.String() && amount != nil && validMemo {
					log.Printf("TM Tx: %s, Msg Index: %d\n", txQueryResult.Hash.String(), i)
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
	}

	return nil
}
