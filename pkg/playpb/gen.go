package playpb

//go:generate protomod protoc -- -I=../.. --go_out=paths=source_relative:../.. --go-grpc_out=paths=source_relative:../.. pkg/playpb/playground.proto
