package cosmos

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceCore/kafka"
	tmTypes "github.com/tendermint/tendermint/types"
	"log"
)

func HandleNewBlock(chain *relayer.Chain, clientCtx *client.Context, blockEvent tmTypes.EventDataNewBlock, kafkaState kafka.KafkaState) {
	fmt.Printf("Cosmos New Block: %d\n", blockEvent.Block.Height)

	fromAccount, err := chain.GetAddress()
	if err != nil {
		log.Fatalln(err.Error())
	}
	toAccount, err := sdk.AccAddressFromBech32("cosmos120fgcs32s8wus7k80ysfszwl275x4v87wuuxd9")
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, transaction := range blockEvent.Block.Txs {
		handleEncodeTx(clientCtx, transaction, kafkaState)
	}

	if blockEvent.Block.Height%10 == 0 {
		response, ok, err := chain.SendMsg(banktypes.NewMsgSend(fromAccount, toAccount, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1)))))
		if err != nil {
			log.Println(err.Error())
		}
		if !ok {
			fmt.Println("Transaction %s not ok", response.TxHash)
		}
	}

}
