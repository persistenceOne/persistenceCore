package status

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

var db *badger.DB

func InitializeDB(dbPath string, tendermintStart, ethereumStart int64) (*badger.DB, error) {
	dbTemp, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatalln(err)
	}
	db = dbTemp

	if tendermintStart > 0 {
		err = SetCosmosStatus(tendermintStart - 1)
		if err != nil {
			return db, err
		}
	}

	if ethereumStart > 0 {
		err = SetEthereumStatus(ethereumStart - 1)
		if err != nil {
			return db, err
		}
	}

	return db, nil
}
