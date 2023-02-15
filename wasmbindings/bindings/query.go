package bindings

type OracleQuery struct {
	GetExchangeRate     *GetExchangeRate     `json:"get_exchange_rate"`
	GetAllExchangeRates *GetAllExchangeRates `json:"get_all_exchange_rates"`
}

type GetExchangeRate struct {
	Symbol string `json:"symbol"`
}

type GetAllExchangeRates struct{}

type GetExchangeRateResponse struct {
	ExchangeRate uint64 `json:"exchange_rate"`
}

type GetAllExchangeRateResponse struct {
	ExchangeRate []uint64 `json:"exchange_rate"`
}
