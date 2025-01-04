# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN apk add --no-cache git make

WORKDIR /app

COPY . .

RUN go mod download && \
  go install go.opentelemetry.io/collector/cmd/builder@v0.109.0

RUN builder --config=config/opentelemetry-collector-builder/dev-manifest.yaml --skip-strict-versioning

FROM alpine:latest

WORKDIR /otel

COPY --from=builder /app/build/monitoringmanager .
COPY --from=builder /app/config/opentelemetry-collector/opentelemetry-config.yaml .

ENTRYPOINT [ "./monitoringmanager", "--config=opentelemetry-config.yaml"]
