module github.com/binance-chain/bep3-deputy

go 1.13

require (
	github.com/aws/aws-sdk-go v1.27.0
	github.com/binance-chain/go-sdk v1.2.0
	github.com/binance-chain/ledger-cosmos-go v0.9.9 // indirect
	github.com/cespare/cp v1.1.1 // indirect
	github.com/ethereum/go-ethereum v1.9.10
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/golang/mock v1.4.3
	github.com/gorilla/mux v1.7.3
	github.com/jinzhu/gorm v1.9.10
	github.com/kava-labs/cosmos-sdk v0.38.3-stable.0.20200520223313-bfbe25d175da
	github.com/kava-labs/go-sdk v0.1.6
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.32.7
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
)

replace github.com/zondax/hid => github.com/binance-chain/hid v0.9.1-0.20190807012304-e1ffd6f0a3cc
