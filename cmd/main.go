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
	"log"
	queueing "message-queueing"
	"message-queueing/http"
	"message-queueing/otel"
	"net"
	nethttp "net/http"
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

	file, err := os.OpenFile(fmt.Sprintf("data/%s", "abc"), os.O_CREATE|os.O_RDWR, 0600)
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

	handler := http.NewServer(service)
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
