package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceCore/pStake/abi"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
)

type EthTxMsg struct {
	Address common.Address `json:"address"`
	Amount  *big.Int       `json:"amount"`
}

func SendTxToEth(client *ethclient.Client, ethTxMsgs []EthTxMsg, gasLimit uint64) (string, error) {
	if len(ethTxMsgs) == 0 {
		return "", fmt.Errorf("number of txs to be send to ethereum is 0")
	}
	ctx := context.Background()
	publicKey := constants.EthAccountPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalln("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(constants.EthAccountPrivateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	contractAddress := common.HexToAddress(constants.TokenWrapperAddress)
	instance, err := abi.NewAbi(contractAddress, client)
	if err != nil {
		return "", err
	}

	addresses := make([]common.Address, len(ethTxMsgs))
	amounts := make([]*big.Int, len(ethTxMsgs))
	for i, ethTxMsg := range ethTxMsgs {
		addresses[i] = ethTxMsg.Address
		amounts[i] = ethTxMsg.Amount
	}

	tx, err := instance.GenerateUTokensInBatch(auth, addresses, amounts)
	if err != nil {
		return "", err
	}
	return tx.Hash().String(), err

}
