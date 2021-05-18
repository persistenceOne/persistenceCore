package status

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
)

const (
	COSMOS   = "COSMOS"
	ETHEREUM = "ETHEREUM"
)

type Status struct {
	Name            string
	LastCheckHeight int64
}

func GetCosmosStatus() (Status, error) {
	var status Status
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(COSMOS))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &status)
			return err
		})
		return err
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func SetCosmosStatus(height int64) error {
	status := Status{
		Name:            COSMOS,
		LastCheckHeight: height,
	}
	b, err := json.Marshal(status)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(status.Name), b)
	})
	if err != nil {
		return err
	}
	return nil
}

func GetEthereumStatus() (Status, error) {
	var status Status
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(ETHEREUM))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &status)
			return err
		})
		return err
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func SetEthereumStatus(height int64) error {
	status := Status{
		Name:            ETHEREUM,
		LastCheckHeight: height,
	}
	b, err := json.Marshal(status)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(status.Name), b)
	})
	if err != nil {
		return err
	}
	return nil
}
