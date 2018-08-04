package jaeger

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/clients/etcd"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
)

//InitGlobalTracer ...
func InitGlobalTracer(config *config.JaegerConfig) (io.Closer, error) {
	if config == nil {
		return nil, errors.New("jaeger config is not exist")
	}
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:              jaeger.SamplerTypeConst,
			Param:             1,
			SamplingServerURL: config.Sampler.HostPort,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.Reporter.LogSpans,
			LocalAgentHostPort: config.Reporter.HostPort,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	return cfg.InitGlobalTracer(
		config.ServiceName,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Observer(rpcmetrics.NewObserver(jMetricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
}

//NewJaegerConfig get config from etcd
func NewJaegerConfig(etcdCli *etcd.Client) (string, *config.JaegerConfig) {
	appName := os.Getenv("APP_NAME")
	profile := os.Getenv("PROFILE")
	jaegerKey := etcdCli.GetEtcdKey(profile, appName, "jaeger")
	jaegerConfig := new(config.JaegerConfig)
	err := etcdCli.GetValue(5*time.Second, jaegerKey, jaegerConfig)
	if err != nil {
		log.Println("jaeger config is not exist:", err)
		jaegerConfig = nil
	}
	return jaegerKey, jaegerConfig
}

//InitGlobalTracerFromEtcd ...
func InitGlobalTracerFromEtcd(etcdCli *etcd.Client) (closer io.Closer, err error) {
	jaegerKey, jaegerConfig := NewJaegerConfig(etcdCli)
	c, err := InitGlobalTracer(jaegerConfig)
	if err != nil {
		log.Println("jaeger connect failed")
	}
	closer = c
	log.Println("jaeger inited")
	//change to new db connection when  config changed
	go etcdCli.WatchKey(jaegerKey, jaegerConfig, closer, func() {
		c, err := InitGlobalTracer(jaegerConfig)
		if err != nil {
			log.Println("jaeger connect failed")
			return
		}
		closer.Close()
		closer = c
		log.Println("jaeger changed")
	})
	return closer, err
}
