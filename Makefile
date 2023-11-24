.DEFAULT_GOAL := buildApp

clean:
	rm ./cmd/cmd ./data/*

buildGRPC: message-queueing.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative message-queueing.proto

buildApp: buildGRPC
	go build -o cmd ./cmd
