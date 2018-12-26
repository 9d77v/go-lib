package cli

import (
	"os"
	"testing"
	"time"
)

// func TestNewConfig(t *testing.T) {
// 	os.Setenv("LOCALHOST", "192.168.1.1")
// 	localhost := os.Getenv("LOCALHOST")
// 	t.Log("ENV LOCALHOST", localhost)
// 	configMap := make(map[string]*config.AppConfig)
// 	orderConfig := new(config.AppConfig)
// 	orderDB := new(config.DBConfig)
// 	orderDB.Server = "default"
// 	orderDB.User = "postgres"
// 	orderDB.Driver = "postgres"
// 	orderDB.Host = localhost
// 	orderDB.Port = 5432
// 	orderDB.MaxIdleConns = 10
// 	orderDB.MaxOpenConns = 100
// 	orderDB.Name = "order"
// 	orderDB.Password = "123456"
// 	orderDB.EnableLog = true
// 	orderRedis := new(config.RedisConfig)
// 	orderRedis.Name = 1
// 	orderRedis.Server = "default"
// 	orderRedis.Host = localhost + ":3306"
// 	orderJaeger := new(config.JaegerConfig)
// 	orderJaegerSampler := new(config.Sampler)
// 	orderJaegerSampler.HostPort = localhost + ":6831"
// 	orderJaegerSampler.RefreshInterval = 10
// 	orderJaegerReporter := new(config.Reporter)
// 	orderJaegerReporter.HostPort = localhost + ":6831"
// 	orderJaegerReporter.FlushInterval = 5
// 	orderJaegerReporter.QueueSize = 1000
// 	orderJaegerReporter.LogSpans = true
// 	orderJaeger.Server = "default"
// 	orderJaeger.ServiceName = "order"
// 	orderJaeger.Sampler = orderJaegerSampler
// 	orderJaeger.Reporter = orderJaegerReporter
// 	orderConfig.DB = orderDB
// 	orderConfig.Redis = orderRedis
// 	orderConfig.Jaeger = orderJaeger

// 	userConfig := new(config.AppConfig)
// 	userDB := new(config.DBConfig)
// 	userDB.Server = "default"
// 	userDB.User = "postgres"
// 	userDB.Driver = "postgres"
// 	userDB.Host = localhost
// 	userDB.Port = 5432
// 	userDB.MaxIdleConns = 10
// 	userDB.MaxOpenConns = 100
// 	userDB.Name = "user"
// 	userDB.Password = "123456"
// 	userDB.EnableLog = true
// 	userRedis := new(config.RedisConfig)
// 	userRedis.Name = 0
// 	userRedis.Server = "default"
// 	userRedis.Host = localhost + ":3306"
// 	userJaeger := new(config.JaegerConfig)
// 	userJaegerSampler := new(config.Sampler)
// 	userJaegerSampler.HostPort = localhost + ":6831"
// 	userJaegerSampler.RefreshInterval = 10
// 	userJaegerReporter := new(config.Reporter)
// 	userJaegerReporter.HostPort = localhost + ":6831"
// 	userJaegerReporter.FlushInterval = 5
// 	userJaegerReporter.QueueSize = 1000
// 	userJaegerReporter.LogSpans = true
// 	userJaeger.Server = "default"
// 	userJaeger.ServiceName = "user"
// 	userJaeger.Sampler = userJaegerSampler
// 	userJaeger.Reporter = userJaegerReporter
// 	userConfig.DB = userDB
// 	userConfig.Redis = userRedis
// 	userConfig.Jaeger = userJaeger

// 	configMap["order"] = orderConfig
// 	configMap["user"] = userConfig

// 	t.Log("order", orderConfig.DB)
// 	t.Log("user", userConfig.DB)
// 	type args struct {
// 		configPath string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want map[string]*config.AppConfig
// 	}{
// 		{"test read config from dev directory", args{"dev"}, configMap},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewConfig(tt.args.configPath); !reflect.DeepEqual(got, tt.want) {
// 				for k, v := range got {
// 					t.Log(k, *v)
// 				}
// 				for k, v := range tt.want {
// 					t.Log(k, *v)
// 				}
// 				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestSyncConfig(t *testing.T) {
	os.Setenv("LOCALHOST", "192.168.1.21")
	type args struct {
		etcdServer     string
		dialTimeout    time.Duration
		requestTimeout time.Duration
		configPath     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test init config",
			args{
				dialTimeout:    5 * time.Second,
				requestTimeout: 10 * time.Second,
				configPath:     "dev",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncConfig(tt.args.dialTimeout, tt.args.requestTimeout, tt.args.configPath); (err != nil) != tt.wantErr {
				t.Errorf("SyncConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
