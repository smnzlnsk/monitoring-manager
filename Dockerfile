# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN apk add --no-cache git make

WORKDIR /app

RUN go install go.opentelemetry.io/collector/cmd/builder@v0.109.0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN builder --config=config/opentelemetry-collector-builder/remote-manifest.yaml --skip-strict-versioning

FROM alpine:latest

WORKDIR /otel

RUN apk add --no-cache make curl grep bash

COPY --from=builder /app/build/monitoringmanager .
COPY --from=builder /app/config/opentelemetry-collector/opentelemetry-config.yaml .

# Make the binary executable
RUN chmod +x /otel/monitoringmanager

# Use the standard OpenTelemetry collector command format
ENTRYPOINT ["/otel/monitoringmanager", "--config=/otel/opentelemetry-config.yaml"]
