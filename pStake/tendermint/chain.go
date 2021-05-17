package tendermint

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

func InitializeAndStartChain(chainConfigJsonPath, timeout, homePath string, coinType uint32, mnemonics string) (*relayer.Chain, error) {
	chain, err := fileInputAdd(chainConfigJsonPath)
	to, err := time.ParseDuration(timeout)
	if err != nil {
		return chain, err
	}

	err = chain.Init(homePath, to, nil, true)
	if err != nil {
		return chain, err
	}

	if chain.KeyExists(chain.Key) {
		log.Printf("deleting old key %s\n", chain.Key)
		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			return chain, err
		}
	}

	ko, err := helpers.KeyAddOrRestore(chain, chain.Key, coinType, mnemonics)
	if err != nil {
		return chain, err
	}

	log.Printf("Keys added: %s\n", ko.Address)

	if err = chain.Start(); err != nil {
		if err != tmservice.ErrAlreadyStarted {
			chain.Error(err)
			return chain, err
		}
	}
	return chain, nil
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
