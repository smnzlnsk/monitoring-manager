package mqttreceiver // import github.com/smnzlnsk/mqttreceiver

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type marshaler struct {
	metricsMarshaler   pmetric.Marshaler
	metricsUnmarshaler pmetric.Unmarshaler
}

func newMarshaler(encoding *component.ID, host component.Host) (*marshaler, error) {
	var metricsMarshaler pmetric.Marshaler = &pmetric.ProtoMarshaler{}
	var metricsUnmarshaler pmetric.Unmarshaler = &pmetric.ProtoUnmarshaler{}

	if encoding != nil {
		ext, ok := host.GetExtensions()[*encoding]
		if !ok {
			return nil, fmt.Errorf("unknown encoding %q", encoding)
		}

		metricsMarshaler, _ = ext.(pmetric.Marshaler)
		metricsUnmarshaler, _ = ext.(pmetric.Unmarshaler)
	}

	m := marshaler{
		metricsMarshaler:   metricsMarshaler,
		metricsUnmarshaler: metricsUnmarshaler,
	}
	return &m, nil
}
