package responses

type TxForTxByHeight struct {
	Hash     string `json:"hash"`
	Height   string `json:"height"`
	Index    int    `json:"index"`
	TxResult struct {
		Code      int    `json:"code"`
		Data      string `json:"data"`
		Log       string `json:"log"`
		Info      string `json:"info"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
		Events    []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Index bool   `json:"index"`
			} `json:"attributes"`
		} `json:"events"`
		Codespace string `json:"codespace"`
	} `json:"tx_result"`
	Tx string `json:"tx"`
}

type TxByHeightResult struct {
	Txs        []TxForTxByHeight `json:"txs"`
	TotalCount string            `json:"total_count"`
}

type TxByHeightResponse struct {
	Result TxByHeightResult `json:"result"`
}
