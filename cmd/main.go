package main

import (
	"flag"
	"fmt"
	api "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"log"
	queueing "message-queueing"
	"message-queueing/otel"
	"net"
	"os"
	"time"
)

var port = flag.Int("port", 8080, "the port of the application")

func main() {
	flag.Parse()

	setupOTel()

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

	repo = otel.WrapRepository(repo)

	service := queueing.NewQueueService(repo, queueing.NewTimeoutQueue())
	service = queueing.NewWriteAheadLogQueueService(walFile, service)
	service, err = otel.WrapService(service)
	if err != nil {
		panic(err)
	}

	server, err := otel.WrapQueueServiceServer(NewMessageQueueingServer(service))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	queueing.RegisterQueueServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}

func setupOTel() {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(), stdouttrace.WithWriter(log.Writer()), stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		panic(err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("message-queueing"),
		),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(r),
	)

	meterExporter, err := stdoutmetric.New(
		stdoutmetric.WithoutTimestamps(), stdoutmetric.WithPrettyPrint(), stdoutmetric.WithWriter(log.Writer()),
	)
	if err != nil {
		panic(err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(meterExporter, sdkmetric.WithInterval(time.Minute))),
		sdkmetric.WithResource(r),
	)

	api.SetTracerProvider(tracerProvider)
	api.SetMeterProvider(meterProvider)
}
