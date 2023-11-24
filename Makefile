.DEFAULT_GOAL := buildApp

clean:
	rm ./cmd/cmd ./data/*

buildGRPC: message_queueing.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative message_queueing.proto

test: buildGRPC
	go test ./...

buildApp: buildGRPC
	go build -o cmd ./cmd

buildImage: buildGRPC
	docker buildx build --load -t message-queueing:latest .
