package cosmos

import (
	"encoding/json"
	"github.com/cosmos/relayer/helpers"
	"github.com/cosmos/relayer/relayer"
	tmservice "github.com/tendermint/tendermint/libs/service"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var Chain *relayer.Chain

func initializeAndStartChain(chainConfigJsonPath, timeout, homePath string, coinType uint32, mnemonics string) (err error) {
	Chain, err = fileInputAdd(chainConfigJsonPath)
	to, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}

	err = Chain.Init(homePath, to, nil, true)
	if err != nil {
		return err
	}

	if Chain.KeyExists(Chain.Key) {
		log.Printf("deleting old key %s\n", Chain.Key)
		err = Chain.Keybase.Delete(Chain.Key)
		if err != nil {
			return err
		}
	}

	ko, err := helpers.KeyAddOrRestore(Chain, Chain.Key, coinType, mnemonics)
	if err != nil {
		return err
	}

	log.Printf("Keys added: %s\n", ko.Address)

	if err = Chain.Start(); err != nil {
		if err != tmservice.ErrAlreadyStarted {
			Chain.Error(err)
			return err
		}
	}
	return err
}

func fileInputAdd(file string) (*relayer.Chain, error) {
	// If the user passes in a file, attempt to read the chain config from that file
	c := &relayer.Chain{}
	if _, err := os.Stat(file); err != nil {
		return c, err
	}

	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}

	if err = json.Unmarshal(byt, c); err != nil {
		return c, err
	}

	return c, nil
}
