module github.com/osmosis-labs/bech32-ibc

go 1.17

require (
	github.com/avast/retry-go v2.6.0+incompatible
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/ibc-go/v2 v2.0.0
	github.com/cosmos/relayer v1.0.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/ory/dockertest/v3 v3.6.2
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.41.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/osmosis-labs/cosmos-sdk v0.44.2-osmo
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2

)
