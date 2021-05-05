package queries

import (
	"fmt"
	"github.com/persistenceOne/persistenceCore/pStake/responses"
)

func GetABCI(rpcAddress string) (responses.ABCIResponse, error) {
	var abci responses.ABCIResponse
	url := rpcAddress + "/abci_info"
	err := get(url, &abci)
	if err != nil {
		return responses.ABCIResponse{}, err
	}
	return abci, err
}

func GetTxsByHeight(rpcAddress, height string) (responses.TxByHeightResponse, error) {
	var txByHeight responses.TxByHeightResponse
	url := rpcAddress + fmt.Sprintf("/tx_search?query=\"tx.height=%s\"", height)
	err := get(url, &txByHeight)
	if err != nil {
		return responses.TxByHeightResponse{}, err
	}
	return txByHeight, err
}

func GetTxHash(restAddress, txHash string) (responses.TxHashResponse, error) {
	var txHashResponse responses.TxHashResponse
	url := restAddress + "/cosmos/tx/v1beta1/txs/" + txHash
	err := get(url, &txHashResponse)
	if err != nil {
		return responses.TxHashResponse{}, err
	}
	return txHashResponse, err
}
