protoc:
	rm -f ./internal/pb/*.go
	mkdir -p ./internal/pb
	protoc -I ./proto \
	--go_out ./internal/pb --go_opt paths=source_relative \
	--go-grpc_out ./internal/pb --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./internal/pb --grpc-gateway_opt paths=source_relative \
	proto/*.proto
server:
	go run cmd/main.go
.PHONY: protoc
