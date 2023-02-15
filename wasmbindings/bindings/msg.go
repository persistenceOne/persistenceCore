package bindings

type CheckersMsg struct {
	UpdateExchangeRate *UpdateExchangeRate `json:"update_move_count,omitempty"`
}

type UpdateExchangeRate struct {
	Symbol string `json:"symbol"`
	Rate   uint64 `json:"rate"`
}
