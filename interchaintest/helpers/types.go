package helpers

import "cosmossdk.io/math"

type QueryMsg struct {
	// IBCHooks
	GetCount      *GetCountQuery      `json:"get_count,omitempty"`
	GetTotalFunds *GetTotalFundsQuery `json:"get_total_funds,omitempty"`

	// Superfluid LP
	GetTotalAmountLocked *GetTotalAmountLockedQuery `json:"total_amount_locked,omitempty"`
}

type ExecMsg struct {
	LockLstAsset *LockLstAssetMsg `json:"lock_lst_asset,omitempty"`
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

type GetTotalAmountLockedQuery struct {
	// {"total_amount_locked": {"user":"persistence1..."}}
	AssetInfo AssetInfo `json:"asset_info"`
	User      string    `json:"user"`
}

type GetTotalAmountLockedResponse struct {
	Data math.Int `json:"data"`
}

type Asset struct {
	Amount math.Int  `json:"amount"`
	Info   AssetInfo `json:"info"`
}

type AssetInfo struct {
	NativeToken NativeTokenInfo `json:"native_token"`
}

type NativeTokenInfo struct {
	Denom string `json:"denom"`
}

type LockLstAssetMsg struct {
	Asset Asset `json:"asset"`
}

type SuperFluidInstantiateMsg struct {
	VaultAddress          string      `json:"vault_addr"`
	Owner                 string      `json:"owner"`
	AllowedLockableTokens []AssetInfo `json:"allowed_lockable_tokens"`
}
