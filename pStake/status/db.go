package status

import (
	"github.com/btcsuite/goleveldb/leveldb"
	"log"
)

var db *leveldb.DB

func InitializeDB(dbPath string, tendermintStart, ethereumStart int64) (*leveldb.DB, error) {
	dbTemp, err := leveldb.OpenFile(dbPath, nil)
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
		err = SetEthereumStatus(ethereumStart)
		if err != nil {
			return db, err
		}
	}

	return db, nil
}
