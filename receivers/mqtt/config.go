package mqttreceiver // import github.com/smnzlnsk/mqttreceiver

import (
	"errors"
	"net"
	"time"

	"go.opentelemetry.io/collector/component"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
	Interval            string        `mapstructure:"interval"`
	ClientID            string        `mapstructure:"client_id"`
	Topic               string        `mapstructure:"topic"`
	EncodingExtensionID *component.ID `mapstructure:"encoding_extension"`
	Broker              BrokerConfig  `mapstructure:"broker"`
}

type BrokerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	// client ID
	if cfg.ClientID == "" {
		return errors.New("client id has to be set")
	}
	// topic
	if cfg.Topic == "" {
		return errors.New("topic has to be set")
	}
	// interval
	interval, _ := time.ParseDuration(cfg.Interval)
	if interval.Seconds() < 1 {
		return errors.New("interval cannot be sub-second")
	}
	if interval.Seconds() > 60 {
		return errors.New("interval cannot be more than a minute")
	}
	// broker.host
	//if !isValidIP(cfg.Broker.Host) {
	//return errors.New("broker.host is invalid")
	//}
	// broker.port
	if cfg.Broker.Port > 65535 || cfg.Broker.Port < 1 {
		return errors.New("broker.port is invalid")
	}
	return nil
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}