package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/websocket"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/queries"
	"github.com/persistenceOne/persistenceCore/pStake/responses"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func GetCmd(initClientCtx client.Context) *cobra.Command {
	pStakeCommand := &cobra.Command{
		Use:   "pStake",
		Short: "Persistence Hub Node Daemon (server)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var ports, err = cmd.Flags().GetString("ports")
			fmt.Println(ports, err)
			if err != nil {
				return err
			}
			go kafkaRoutine(ports)
			run(initClientCtx)

			return nil
		},
	}

	pStakeCommand.Flags().String("ports", "PLAINTEXT://localhost:9092", "ports kafka brokers are running on, --ports https://192.100.10.10:443,https://192.100.10.11:443")

	return pStakeCommand
}

// kafkaRoutine: starts kafka in a separate goRoutine, consumers will each start in different go routines
// no need to store any db, producers and consumers are inside kafkaState struct.
// use kafka.KafkaProducerDeliverMessage() -> to produce message
// use kafka.KafkaTopicConsumer -> to consume messages.
func kafkaRoutine(ports string) {
	portsList := strings.Split(ports, ",")
	_ = kafka.NewKafkaState(portsList)

	time.Sleep(1000000000)

	go consumeMsgSend()
	// go consume other messages

	fmt.Println("started consumers")
}
func consumeMsgSend() {
	for {
		//consume logic here.
		time.Sleep(kafka.SleepRoutine)
	}
}

var lastBlockHeight int64 = 1 //6048016

func setAppLastBlockHeight(height int64) {
	lastBlockHeight = height
}

func getAppLastBlockHeight() int64 {
	return lastBlockHeight
}

func run(initClientCtx client.Context) {

	if constants.HttpMethod {
		httpMethod(initClientCtx)
	} else {
		webSocketMethod(initClientCtx)
	}

}

func handleSendCoinMsg(msg *banktypes.MsgSend) {
	fmt.Println(msg.String())
}

func handleEncodeTx(initClientCtx client.Context, encodedTx string) {
	decodedTx, err := base64.StdEncoding.DecodeString(encodedTx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	txInterface, err := initClientCtx.TxConfig.TxDecoder()(decodedTx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tx, ok := txInterface.(signing.Tx)
	if !ok {
		log.Fatalln("Unable to parse tx")
	}

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *banktypes.MsgSend:
			handleSendCoinMsg(msg)
		default:

		}
	}
}

func httpMethod(initClientCtx client.Context) {
	for {
		abciResponse, err := queries.GetABCI()
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		lastBlockHeight, err := strconv.ParseInt(abciResponse.Result.Response.LastBlockHeight, 10, 64)
		if err != nil {
			log.Fatalln(err)
		}

		if lastBlockHeight > getAppLastBlockHeight() {

			checkHeight := getAppLastBlockHeight() + 1
			response, err := queries.GetTxsByHeight(strconv.FormatInt(checkHeight, 10))
			if err != nil {
				log.Println(err)
				time.Sleep(2 * time.Second)
				continue
			}

			for _, txResponse := range response.Result.Txs {
				fmt.Println(checkHeight)
				fmt.Println(txResponse.TxResult.Code)
				if txResponse.TxResult.Code == 0 {
					handleEncodeTx(initClientCtx, txResponse.Tx)
				}
			}
			setAppLastBlockHeight(checkHeight)
		} else {
			time.Sleep(3500 * time.Millisecond)
		}
	}
}

func webSocketMethod(initClientCtx client.Context) {
	wsURL := url.URL{Scheme: "wss", Host: constants.WEBOSCKET_URL, Path: "/websocket"}
	if !constants.TLS_ENABLED {
		wsURL = url.URL{Scheme: "ws", Host: constants.WEBOSCKET_URL, Path: "/websocket"}
	}
	log.Printf("connecting to %s", wsURL.String())

	connection, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	err = connection.WriteMessage(websocket.TextMessage, []byte(`{"method":"subscribe", "id":"dontcare","jsonrpc":"2.0","params":["tm.event='Tx'"]}`))
	if err != nil {
		log.Println("write:", err)
		return
	}

	go func() {
		for {
			var webSocketTx responses.WebSocketTx
			err := connection.ReadJSON(&webSocketTx)
			if err != nil {
				log.Println("WebSocket Connection Error:", err)
				return
			}
			if webSocketTx.Result.Data.Value.TxResult.Result.Code == 0 {
				handleEncodeTx(initClientCtx, webSocketTx.Result.Data.Value.TxResult.Tx)
			}
		}
	}()

	connection.SetCloseHandler(func(code int, text string) error {
		log.Printf("Close Handler: code - %d, message: %s", code, text)
		connection.WriteMessage(websocket.TextMessage, []byte(`{"method":"subscribe", "id":"dontcare","jsonrpc":"2.0","params":["tm.event='Tx'"]}`))
		return nil
	})

	defer func(connection *websocket.Conn) {
		err := connection.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(connection)

	select {}

}
