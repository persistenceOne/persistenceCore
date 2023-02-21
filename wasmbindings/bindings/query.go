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
	ExchangeRate string `json:"exchange_rate"`
}

type GetAllExchangeRateResponse struct {
	ExchangeRate []string `json:"all_exchange_rates"`
}
