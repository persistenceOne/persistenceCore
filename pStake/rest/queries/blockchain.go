package queries

import (
	"fmt"
	responses2 "github.com/persistenceOne/persistenceCore/pStake/rest/responses"
)

func GetABCI(rpcAddress string) (responses2.ABCIResponse, error) {
	var abci responses2.ABCIResponse
	url := rpcAddress + "/abci_info"
	err := get(url, &abci)
	if err != nil {
		return responses2.ABCIResponse{}, err
	}
	return abci, err
}

func GetTxsByHeight(rpcAddress, height string) (responses2.TxByHeightResponse, error) {
	var txByHeight responses2.TxByHeightResponse
	url := rpcAddress + fmt.Sprintf("/tx_search?query=\"tx.height=%s\"", height)
	err := get(url, &txByHeight)
	if err != nil {
		return responses2.TxByHeightResponse{}, err
	}
	return txByHeight, err
}

func GetTxHash(restAddress, txHash string) (responses2.TxHashResponse, error) {
	var txHashResponse responses2.TxHashResponse
	url := restAddress + "/cosmos/tx/v1beta1/txs/" + txHash
	err := get(url, &txHashResponse)
	if err != nil {
		return responses2.TxHashResponse{}, err
	}
	return txHashResponse, err
}
