package config

import (
	"os"
	"strings"
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

//DefaultConfig ...
type DefaultConfig struct {
	DB            map[string]*DBConfig       `yaml:"db"`
	Redis         map[string]*RedisConfig    `yaml:"redis"`
	Jaeger        map[string]*JaegerConfig   `yaml:"jaeger"`
	Rabbitmq      map[string]*RabbitmqConfig `yaml:"rabbitmq"`
	ExpressConfig *ExpressConfig             `yaml:"express"`
}

//AppConfig ...
type AppConfig struct {
	DB       *DBConfig       `yaml:"db"`
	Redis    *RedisConfig    `yaml:"redis"`
	Jaeger   *JaegerConfig   `yaml:"jaeger"`
	Rabbitmq *RabbitmqConfig `yaml:"rabbitmq"`
}

//ExpressConfig ...
type ExpressConfig struct {
	Key      string `yaml:"key"`
	Customer string `yaml:"customer"`
	Secret   string `yaml:"secret"`
}

const (
	//LOCALHOST ...11
	LOCALHOST = "${LOCALHOST}"
)

//EnvVar environment variable
func EnvVar(v string) string {
	if len(v) < 4 {
		return ""
	}
	return v[2 : len(v)-1]
}

func newReplaceMap() map[string]string {
	replaceMap := make(map[string]string)
	replaceMap[LOCALHOST] = os.Getenv(EnvVar(LOCALHOST))
	return replaceMap
}

//NewAppConfig return new app config from default config and old app config.
func NewAppConfig(df *DefaultConfig, af *AppConfig) *AppConfig {
	replaceMap := newReplaceMap()
	config := new(AppConfig)
	config.DB = replaceDB(df, af.DB, replaceMap)
	config.Redis = replaceRedis(df, af.Redis, replaceMap)
	config.Jaeger = replaceJaeger(df, af.Jaeger, replaceMap)
	config.Rabbitmq = replaceRabbitmq(df, af.Rabbitmq, replaceMap)
	return config
}

func replaceDB(df *DefaultConfig, config *DBConfig, replaceMap map[string]string) *DBConfig {
	if config == nil {
		return nil
	}
	tmp := *df.DB[config.Server]
	newConfig := &tmp
	if newConfig != nil {
		newConfig.Server = config.Server
		newConfig.Name = config.Name
		newConfig.Host = replaceLocalhost(newConfig.Host, replaceMap)
	}
	return newConfig
}

func replaceRedis(df *DefaultConfig, config *RedisConfig, replaceMap map[string]string) *RedisConfig {
	if config == nil {
		return nil
	}
	tmp := *df.Redis[config.Server]
	newConfig := &tmp
	if newConfig != nil {
		newConfig.Server = config.Server
		newConfig.Name = config.Name
		newConfig.Host = replaceLocalhost(newConfig.Host, replaceMap)
	}
	return newConfig
}

func replaceJaeger(df *DefaultConfig, config *JaegerConfig, replaceMap map[string]string) *JaegerConfig {
	if config == nil {
		return nil
	}
	tmp := *df.Jaeger[config.Server]
	newConfig := &tmp
	if newConfig != nil {
		newConfig.Server = config.Server
		newConfig.ServiceName = config.ServiceName
		newConfig.Reporter.HostPort = replaceLocalhost(newConfig.Reporter.HostPort, replaceMap)
		newConfig.Sampler.HostPort = replaceLocalhost(newConfig.Sampler.HostPort, replaceMap)
	}
	return newConfig
}

func replaceRabbitmq(df *DefaultConfig, config *RabbitmqConfig, replaceMap map[string]string) *RabbitmqConfig {
	if config == nil {
		return nil
	}
	tmp := *df.Rabbitmq[config.Server]
	newConfig := &tmp
	if newConfig != nil {
		newConfig.Server = config.Server
		newConfig.Host = replaceLocalhost(newConfig.Host, replaceMap)
	}
	return newConfig
}

func replaceLocalhost(value string, replaceMap map[string]string) string {
	return strings.Replace(value, LOCALHOST, replaceMap[LOCALHOST], 1)
}
