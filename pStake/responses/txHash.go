package responses

import (
	"time"
)

type TxResponse struct {
	Height string `json:"height"`
	Txhash string `json:"txhash"`
	Code   int    `json:"code"`
	Tx     struct {
		Body struct {
			Messages []interface{} `json:"messages"`
			Memo     string        `json:"memo"`
		} `json:"body"`
	} `json:"tx"`
	Timestamp time.Time `json:"timestamp"`
}

type TxHashResponse struct {
	TxResponse TxResponse `json:"tx_response"`
}
