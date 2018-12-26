package config

import (
	"os"
	"strings"
)

//DefaultConfig ...
type DefaultConfig struct {
	DB            map[string]*DBConfig       `yaml:"db" json:"db"`
	Redis         map[string]*RedisConfig    `yaml:"redis" json:"redis"`
	Jaeger        map[string]*JaegerConfig   `yaml:"jaeger" json:"jaeger"`
	Rabbitmq      map[string]*RabbitmqConfig `yaml:"rabbitmq" json:"rabbitmq"`
	ExpressConfig *ExpressConfig             `yaml:"express" json:"express"`
	SMSConfig     *SMSConfig                 `yaml:"sms" json:"sms"`
}

//AppConfig ...
type AppConfig struct {
	DB       *DBConfig       `yaml:"db" json:"db"`
	Redis    *RedisConfig    `yaml:"redis" json:"redis"`
	Jaeger   *JaegerConfig   `yaml:"jaeger" json:"jaeger"`
	Rabbitmq *RabbitmqConfig `yaml:"rabbitmq" json:"rabbitmq"`
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
