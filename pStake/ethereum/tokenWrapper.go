package ethereum

import (
	"context"
	"crypto/ecdsa"
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

func SendTxToEth(client *ethclient.Client, ethTxMsg EthTxMsg, gasLimit uint64) (string, error) {
	ctx := context.Background()
	publicKey := constants.EthAccountPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
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

	tx, err := instance.GenerateUTokens(auth, ethTxMsg.Address, ethTxMsg.Amount)
	if err != nil {
		return "", err
	}
	return tx.Hash().String(), err

}
