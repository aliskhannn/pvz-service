run:
	go run cmd/server/main.go

proto:
	protoc --go_out=. --go-grpc_out=. api/proto/*.proto

.PHONY: run proto lint
