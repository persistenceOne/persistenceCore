package queries

import (
	"fmt"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/responses"
)

func GetABCI() (responses.ABCIResponse, error) {
	var abci responses.ABCIResponse
	url := constants.RPC_URL + "/abci_info"
	err := get(url, &abci)
	if err != nil {
		return responses.ABCIResponse{}, err
	}
	return abci, err
}

func GetTxsByHeight(height string) (responses.TxByHeightResponse, error) {
	var txByHeight responses.TxByHeightResponse
	url := constants.RPC_URL + fmt.Sprintf("/tx_search?query=\"tx.height=%s\"", height)
	err := get(url, &txByHeight)
	if err != nil {
		return responses.TxByHeightResponse{}, err
	}
	return txByHeight, err
}

func GetTxHash(txHash string) (responses.TxHashResponse, error) {
	var txHashResponse responses.TxHashResponse
	url := constants.REST_URL + "/cosmos/tx/v1beta1/txs/" + txHash
	err := get(url, &txHashResponse)
	if err != nil {
		return responses.TxHashResponse{}, err
	}
	return txHashResponse, err
}
