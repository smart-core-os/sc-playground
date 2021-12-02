module github.com/smart-core-os/sc-playground

go 1.16

require (
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/google/go-cmp v0.5.6
	github.com/improbable-eng/grpc-web v0.14.0
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/rs/cors v1.7.0
	github.com/smart-core-os/sc-api/go v1.0.0-beta.26
	github.com/smart-core-os/sc-golang v0.0.0-20211025094731-4496ddbbd1f0
	github.com/soheilhy/cmux v0.1.5
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	nhooyr.io/websocket v1.8.7 // indirect
)

replace (
	github.com/smart-core-os/sc-api/go => github.com/smart-core-os/sc-api/go v1.0.0-beta.26.0.20211129083553-12fbc9816ef7
	github.com/smart-core-os/sc-golang => github.com/smart-core-os/sc-golang v0.0.0-20211202091942-cfaaba2b5da4
)
