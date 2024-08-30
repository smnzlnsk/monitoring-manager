package mqttreceiver // import github.com/smnzlnsk/mqttreceiver

import (
	"context"
	"time"

	"github.com/smnzlnsk/mqttreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	defaultInterval   = time.Second * 1
	defaultClientID   = "test-mqtt-receiver-client"
	defaultTopic      = "telemetry/metrics"
	defaultEncoding   = "otlp"
	defaultBroker     = "127.0.0.1"
	defaultBrokerPort = 1883
)

// NewFactory creates a factory for mqtt-exporter
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Interval: string(defaultInterval),
		ClientID: defaultClientID,
		Topic:    defaultTopic,
		Broker: BrokerConfig{
			Host: defaultBroker,
			Port: defaultBrokerPort,
		},
	}
}

func createMetricsReceiver(
	ctx context.Context,
	set receiver.Settings,
	cfg component.Config,
	con consumer.Metrics,
) (receiver.Metrics, error) {
	logger := set.Logger
	config := (cfg.(*Config))
	mr, err := newMQTTReceiver(config, logger, con)
	if err != nil {
		return nil, err
	}
	return mr, nil
}
