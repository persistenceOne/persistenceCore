package responses

type WebSocketTx struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  struct {
		Query string `json:"query"`
		Data  struct {
			Type  string `json:"type"`
			Value struct {
				TxResult struct {
					Height string `json:"height"`
					Tx     string `json:"tx"`
					Result struct {
						Code      int    `json:"code,omitempty"`
						Log       string `json:"log"`
						GasWanted string `json:"gas_wanted"`
						GasUsed   string `json:"gas_used"`
						Codespace string `json:"codespace"`
					} `json:"result"`
				} `json:"TxResult"`
			} `json:"value"`
		} `json:"data"`
		Events struct {
			TmEvent  []string `json:"tm.event"`
			TxHash   []string `json:"tx.hash"`
			TxHeight []string `json:"tx.height"`
		} `json:"events"`
	} `json:"result"`
}
