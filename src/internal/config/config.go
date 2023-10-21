package config

import (
	"github.com/caarlos0/env/v9"
)

type MonitoringConfig struct {
	TraceDestination   string `env:"TRACE_DESTINATION"`
	MetricsDestination string `env:"METRICS_DESTINATION"`
	LogFileName        string `env:"LOG_FILE_NAME"`
}

type GenericConfig struct {
	ServicePort uint16 `env:"SERVICE_PORT" envDefault:"5980"`
	ServiceBind string `env:"SERVICE_BIND" envDefault:"0.0.0.0"`
	Environment string `env:"ENVIRONMENT" envDefault:"dev"`
}

type Config struct {
	MonitoringConfig
	GenericConfig
	SumURL string `env:"SUM_URL" envDefault:"https://data.cityofchicago.org/resource/ydr8-5enu.json"`
}

func InitConfig() (cfg *Config, err error) {
	cfg = &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return
}
