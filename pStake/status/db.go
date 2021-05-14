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

	err = SetCosmosStatus(tendermintStart)
	if err != nil {
		return db, err
	}

	err = SetEthereumStatus(ethereumStart)
	if err != nil {
		return db, err
	}

	return db, nil
}
