package helpers

import "cosmossdk.io/math"

type QueryMsg struct {
	// IBCHooks
	GetCount      *GetCountQuery      `json:"get_count,omitempty"`
	GetTotalFunds *GetTotalFundsQuery `json:"get_total_funds,omitempty"`

	// Superfluid LP
	GetLockedLstForUser *GetLockedLstForUserQuery `json:"locked_lst_for_user,omitmepty"`
}

type ExecMsg struct {
	LockLstAssetForUser *LockLstAssetForUserMsg `json:"lock_lst_asset_for_user,omitmepty"`
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

type GetLockedLstForUserQuery struct {
	// {"locked_lst_for_user": {"user":"persistence1..."}}
	Asset Asset  `json:"asset"`
	User  string `json:"user"`
}

type GetLockedLstForUserResponse struct {
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

type LockLstAssetForUserMsg struct {
	Asset Asset  `json:"asset"`
	User  string `json:"user"`
}
