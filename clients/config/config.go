package config

import (
	"time"
)

//DBConfig config of relational database
type DBConfig struct {
	Server       string `yaml:"server"`
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	Port         uint   `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxIdleConns uint   `yaml:"max_idle_conns"`
	MaxOpenConns uint   `yaml:"max_open_conns"`
	EnableLog    bool   `yaml:"enable_log"`
}

//RedisConfig config of redis
type RedisConfig struct {
	Server   string `yaml:"server"`
	Host     string `yaml:"host"`
	Name     int    `yaml:"name"`
	Password string `yaml:"password"`
}

//ElasticConfig config of elasticsearch brokers
type ElasticConfig struct {
	URLs []string `yaml:"urls"`
}

//Sampler config of jaeger sampler
type Sampler struct {
	Type            string        `yaml:"type"`
	Param           int           `yaml:"param"`
	HostPort        string        `yaml:"host_port"`
	RefreshInterval time.Duration `yaml:"refresh_interval"`
}

//Reporter config of jaeger reporter
type Reporter struct {
	LogSpans      bool          `yaml:"log_spans"`
	HostPort      string        `yaml:"host_port"`
	FlushInterval time.Duration `yaml:"flush_interval"`
	QueueSize     int           `yaml:"queue_size"`
}

//JaegerConfig config of jeager
type JaegerConfig struct {
	Server      string    `yaml:"server"`
	ServiceName string    `yaml:"service_name"`
	Sampler     *Sampler  `yaml:"sampler"`
	Reporter    *Reporter `yaml:"reporter"`
}

//RabbitmqConfig config of rabbitmq
type RabbitmqConfig struct {
	Server   string `yaml:"server"`
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

//ExpressConfig ...
type ExpressConfig struct {
	Key      string `yaml:"key"`
	Customer string `yaml:"customer"`
	Secret   string `yaml:"secret"`
}
