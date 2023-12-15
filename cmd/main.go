package main

import (
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	api "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log"
	queueing "message-queueing"
	"message-queueing/http"
	"message-queueing/otel"
	"message-queueing/replication"
	"net"
	nethttp "net/http"
	"os"
	"os/signal"
	"time"
)

var port = flag.Int("port", 8080, "the port of the application")
var dataDir = flag.String("dataDir", "./data", "the directory to store the data in")

func main() {
	flag.Parse()

	//setupOTel()

	etcdClient, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		},
	)
	if err != nil {
		panic(err)
	}

	controller, err := replication.Open(context.TODO(), etcdClient, *dataDir)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", *dataDir, "abc"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatalf("failed to open file: %w", err)
	}

	storage := queueing.NewIOBlockStorage(file)
	index := queueing.NewNaiveIndex()
	queue := queueing.NewTimeoutQueue()

	queueing.SetupQueueMessageRepository(file, index, queue)

	repo := queueing.NewQueueMessageRepository(storage, index)
	repo = otel.WrapRepository(repo)

	service := queueing.NewQueueService(repo, queue)
	service, err = otel.WrapService(service)
	if err != nil {
		panic(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func(c <-chan os.Signal) {
		<-c
		controller.Close()
		lis.Close()
		os.Exit(2)
	}(signalChannel)

	handler := http.NewServer(service, storage, repo)
	nethttp.Serve(lis, handler)
}

func setupOTel() {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithWriter(log.Writer()), stdouttrace.WithoutTimestamps(),
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
		stdoutmetric.WithoutTimestamps(), stdoutmetric.WithWriter(log.Writer()),
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
