package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	queueing "message-queueing"
	"net"
	"os"
	"time"
)

var port = flag.Int("port", 8080, "the port of the application")

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	walFile, err := os.Create(fmt.Sprintf("data/wal-%x", time.Now().UnixMilli()))
	if err != nil {
		panic(err)
	}

	repo, err := queueing.SetupQueueMessageRepository("abc")
	if err != nil {
		panic(err)
	}

	service := queueing.NewQueueService(repo, queueing.NewTimeoutQueue())
	service = queueing.NewWriteAheadLogQueueService(walFile, service)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	queueing.RegisterQueueServiceServer(grpcServer, queueing.NewMessageQueueingServer(service))
	grpcServer.Serve(lis)
}
