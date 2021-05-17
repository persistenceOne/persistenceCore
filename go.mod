module github.com/persistenceOne/persistenceCore

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.28.0
	github.com/btcsuite/goleveldb v1.0.0
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/relayer v0.9.3
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.10.3
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/shirou/gopsutil v3.21.4+incompatible // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/net v0.0.0-20210510095157-81045d8b478c // indirect
	golang.org/x/sys v0.0.0-20210507161434-a76c4d0a0096 // indirect
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.37.0
	gopkg.in/yaml.v2 v2.4.0
	honnef.co/go/tools v0.1.4
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
