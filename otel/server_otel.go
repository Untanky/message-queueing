package otel

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"time"
)

type otelInterceptor struct {
	numberOfRequests metric.Int64Counter
	numberOfErrors   metric.Int64Counter
	latency          metric.Int64Histogram
}

func NewInterceptor() (*otelInterceptor, error) {
	numberOfRequests, err := meter.Int64Counter("numberOfRequests")
	if err != nil {
		return nil, err
	}

	numberOfErrors, err := meter.Int64Counter("numberOfErrors")
	if err != nil {
		return nil, err
	}

	latency, err := meter.Int64Histogram("requestLatency", metric.WithUnit("μs"))
	if err != nil {
		return nil, err
	}

	return &otelInterceptor{
		numberOfRequests: numberOfRequests,
		numberOfErrors:   numberOfErrors,
		latency:          latency,
	}, nil
}

func (o *otelInterceptor) UnaryInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {
	start := time.Now()
	ctx, span := tracer.Start(ctx, "SubmitMessages")
	defer span.End()

	response, err := handler(ctx, req)

	o.numberOfRequests.Add(ctx, 1)
	if err != nil {
		o.numberOfErrors.Add(ctx, 1)
		span.RecordError(err)
	}

	end := time.Now()
	o.latency.Record(ctx, end.Sub(start).Microseconds())

	return response, err
}
