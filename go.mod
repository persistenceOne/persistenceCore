module github.com/commitHub/commitBlockchain

require (
	github.com/cosmos/cosmos-sdk v0.35.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/tendermint v0.31.5
	golang.org/x/crypto v0.0.0-20180904163835-0709b304e793

)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
