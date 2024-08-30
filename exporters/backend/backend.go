package backendexporter // import github.com/smnzlnsk/monitoringbackendexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type backend struct {
	config *Config
	logger *zap.Logger
	*marshaler
	host   component.Host
	cancel context.CancelFunc
}

func newBackend(cfg *Config, logger *zap.Logger) (*backend, error) {
	receiver := &backend{
		config: cfg,
		logger: logger,
	}
	return receiver, nil
}

func (b *backend) start(ctx context.Context, host component.Host) error {
	ctx = context.Background()
	ctx, b.cancel = context.WithCancel(ctx)
	marshaler, err := newMarshaler(nil, host)
	if err != nil {
		return err
	}
	b.marshaler = marshaler
	b.host = host

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (b *backend) shutdown(ctx context.Context) error {
	if b.cancel != nil {
		b.cancel()
	}
	return nil
}

func (b *backend) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	_, err := b.metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		b.logger.Error("error marshaling metric data")
		return err
	}
	b.logger.Info("received metric data")
	return nil
}
