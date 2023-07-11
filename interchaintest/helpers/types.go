package helpers

// Go based data types for querying on the contract.
// Execute types are not needed here. We just use strings. Could add though in the future and to_string it

// EntryPoint
type QueryMsg struct {
	// IBCHooks
	GetCount      *GetCountQuery      `json:"get_count,omitempty"`
	GetTotalFunds *GetTotalFundsQuery `json:"get_total_funds,omitempty"`
}

type GetTotalFundsQuery struct {
	// {"get_total_funds":{"addr":"persistence1..."}}
	Addr string `json:"addr"`
}

type GetTotalFundsResponse struct {
	// {"data":{"total_funds":[{"denom":"ibc/04F5F501207C3626A2C14BFEF654D51C2E0B8F7CA578AB8ED272A66FE4E48097","amount":"1"}]}}
	Data *GetTotalFundsObj `json:"data"`
}

type GetTotalFundsObj struct {
	TotalFunds []WasmCoin `json:"total_funds"`
}

type WasmCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type GetCountQuery struct {
	// {"get_total_funds":{"addr":"persistence1..."}}
	Addr string `json:"addr"`
}

type GetCountResponse struct {
	// {"data":{"count":0}}
	Data *GetCountObj `json:"data"`
}

type GetCountObj struct {
	Count int64 `json:"count"`
}
