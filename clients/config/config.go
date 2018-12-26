package config

import (
	"time"
)

//DBConfig config of relational database
type DBConfig struct {
	Server       string `yaml:"server" json:"server"`
	Driver       string `yaml:"driver" json:"driver"`
	Host         string `yaml:"host" json:"host"`
	Port         uint   `yaml:"port" json:"port"`
	User         string `yaml:"user" json:"user"`
	Password     string `yaml:"password" json:"password"`
	Name         string `yaml:"name" json:"name"`
	MaxIdleConns uint   `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns uint   `yaml:"max_open_conns" json:"max_open_conns"`
	EnableLog    bool   `yaml:"enable_log" json:"enable_log"`
}

//RedisConfig config of redis
type RedisConfig struct {
	Server   string `yaml:"server" json:"server"`
	Host     string `yaml:"host" json:"host"`
	Name     int    `yaml:"name" json:"name"`
	Password string `yaml:"password" json:"password"`
}

//ElasticConfig config of elasticsearch brokers
type ElasticConfig struct {
	URLs []string `yaml:"urls" json:"urls"`
}

//Sampler config of jaeger sampler
type Sampler struct {
	Type            string        `yaml:"type" json:"type"`
	Param           int           `yaml:"param" json:"param"`
	HostPort        string        `yaml:"host_port" json:"host_port"`
	RefreshInterval time.Duration `yaml:"refresh_interval" json:"refresh_interval"`
}

//Reporter config of jaeger reporter
type Reporter struct {
	LogSpans      bool          `yaml:"log_spans" json:"log_spans"`
	HostPort      string        `yaml:"host_port" json:"host_port"`
	FlushInterval time.Duration `yaml:"flush_interval" json:"flush_interval"`
	QueueSize     int           `yaml:"queue_size" json:"queue_size"`
}

//JaegerConfig config of jeager
type JaegerConfig struct {
	Server      string    `yaml:"server" json:"server"`
	ServiceName string    `yaml:"service_name" json:"service_name"`
	Sampler     *Sampler  `yaml:"sampler" json:"sampler"`
	Reporter    *Reporter `yaml:"reporter" json:"reporter"`
}

//RabbitmqConfig config of rabbitmq
type RabbitmqConfig struct {
	Server   string `yaml:"server" json:"server"`
	Host     string `yaml:"host" json:"host"`
	Port     uint   `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
}

//ExpressConfig ...
type ExpressConfig struct {
	Key      string `yaml:"key" json:"key"`
	Customer string `yaml:"customer" json:"customer"`
	Secret   string `yaml:"secret" json:"secret"`
}
