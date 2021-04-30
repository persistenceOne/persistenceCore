/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/pkg/errors"
)

func ParseFloat64OrReturnBadRequest(s string, defaultIfEmpty float64) (float64, int, error) {
	if len(s) == 0 {
		return defaultIfEmpty, http.StatusAccepted, nil
	}

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return n, http.StatusBadRequest, err
	}

	return n, http.StatusAccepted, nil
}

func SimulationResponse(cdc *codec.LegacyAmino, gas uint64) ([]byte, error) {
	gasEst := rest.GasEstimateResponse{GasEstimate: gas}
	resp, err := cdc.MarshalJSON(gasEst)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return resp, nil
}
