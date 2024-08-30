package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type      = component.MustNewType("mqttreceiver")
	ScopeName = "github.com/smnzlnsk/monitoring-manager/receivers/mqtt"
)
