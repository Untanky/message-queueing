# Build the application from source
ARG BUILDARCH
FROM --platform="linux/amd64" golang:1.21 AS build-stage

ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /message-queueing ./cmd

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /message-queueing /message-queueing

EXPOSE 8080

ENTRYPOINT ["/message-queueing"]
