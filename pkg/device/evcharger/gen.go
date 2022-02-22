package evcharger

//go:generate protomod protoc -- -I../../.. --go_out=paths=source_relative:../../.. --go-grpc_out=paths=source_relative:../../../ --wrapper_out=paths=source_relative:../../../ --router_out=paths=source_relative:../../../ pkg/device/evcharger/evcharger.proto
