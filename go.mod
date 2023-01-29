module github.com/dragos-rebegea/evaluare-tool

go 1.15

require (
	github.com/ElrondNetwork/elrond-go v1.3.36
	github.com/ElrondNetwork/elrond-go-core v1.1.16-0.20220711092037-f35a3a0faf0f
	github.com/ElrondNetwork/elrond-go-logger v1.0.7
	github.com/ElrondNetwork/elrond-sdk-erdgo v1.0.24-0.20220812120443-2d2ef08d5f8a
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v0.0.0-20190301062745-f9e10995c85a
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.8.0
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/stretchr/testify v1.7.1
	github.com/urfave/cli v1.22.9
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gorm.io/driver/mysql v1.4.5
	gorm.io/gorm v1.23.8
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.40 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.40

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.40 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.40

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.40 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.40

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.54-rc3 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.54-rc3
