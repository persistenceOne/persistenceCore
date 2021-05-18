/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
)

// SetTicketIDtoDB : initiates ticketID in Database
func SetTicketIDtoDB(ticket Ticket, kafkaDB *dbm.GoLevelDB, cdc *codec.LegacyAmino, msg []byte) {
	ticketID, Error := cdc.MarshalJSON(ticket)
	if Error != nil {
		panic(Error)
	}

	if Error := kafkaDB.Set(ticketID, msg); Error != nil {
		panic(Error)
	}
}

// AddResponseToDB : Updates response to DB
func AddResponseToDB(ticket Ticket, response []byte, kafkaDB *dbm.GoLevelDB, cdc *codec.LegacyAmino) {
	ticketID, err := cdc.MarshalJSON(ticket)
	if err != nil {
		panic(err)
	}

	err = kafkaDB.SetSync(ticketID, response)
	if err != nil {
		panic(err)
	}
}

// GetResponseFromDB : gives the response from DB
func GetResponseFromDB(ticket Ticket, kafkaDB *dbm.GoLevelDB, cdc *codec.LegacyAmino) []byte {
	ticketID, err := cdc.MarshalJSON(ticket)
	if err != nil {
		panic(err)
	}

	val, _ := kafkaDB.Get(ticketID)

	return val
}

// QueryDB : REST outputs info from DB
func QueryDB(cdc *codec.LegacyAmino, kafkaDB *dbm.GoLevelDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		iDByte, err := cdc.MarshalJSON(vars["ticketID"])
		if err != nil {
			panic(err)
		}

		var response []byte

		check, _ := kafkaDB.Has(iDByte)
		if check {
			response = GetResponseFromDB(Ticket(vars["ticketID"]), kafkaDB, cdc)
		} else {
			output, err := cdc.MarshalJSON("The ticket ID does not exist, it must have been deleted, Query the chain to know")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(fmt.Sprintf("ticket ID does not exist. Error: %s", err.Error())))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(output)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write(response)
	}
}
