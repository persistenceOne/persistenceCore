module github.com/persistenceOne/persistenceCore

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.6
	github.com/cosmos/ibc-go/v2 v2.0.2
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210828152312-66f60bf46e71
	google.golang.org/grpc v1.42.0
	gopkg.in/yaml.v2 v2.4.0
	honnef.co/go/tools v0.0.1-2020.1.4
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.7

replace github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2

replace github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.3
