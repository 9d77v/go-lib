package redis

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	opentracing "github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"

	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/clients/etcd"
	"github.com/9d77v/go-lib/clients/jaeger"
)

//Client redis client with opentracing
type Client struct {
	*redis.Client
	Tracer opentracing.Tracer
	Closer io.Closer
}

//NewClient get redis connection
func NewClient(config *config.RedisConfig) (*Client, error) {
	if config == nil {
		return nil, errors.New("db config is not exist")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       config.Name,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("redis connection failed")
	}
	return &Client{client, nil, nil}, nil
}

//NewClientFromEtcd init redis from etcd config and watch config to update redis
func NewClientFromEtcd(etcdCli *etcd.Client) (redisCli *Client, err error) {
	appName := os.Getenv("APP_NAME")
	profile := os.Getenv("PROFILE")
	redisKey := etcdCli.GetEtcdKey(profile, appName, "redis")
	redisConfig := new(config.RedisConfig)
	err = etcdCli.GetValue(5*time.Second, redisKey, redisConfig)
	if err != nil {
		log.Println("redis config is not exist:", err)
	}
	redis, err := NewClient(redisConfig)
	if err != nil {
		log.Println("redis connect failed")
	}
	tracer, closer, err := jaeger.InitTracerFromEtcd(etcdCli, "redis")
	redis.Tracer = tracer
	redis.Closer = closer
	redisCli = redis
	log.Println("redis inited", redisCli)
	//change to new redis connection when  config changed
	go etcdCli.WatchKey(redisKey, redisConfig, func() {
		redis, err := NewClient(redisConfig)
		if err != nil {
			log.Println("redis connect failed")
			return
		}
		redis.Closer = redisCli.Closer
		redis.Tracer = redisCli.Tracer
		redisCli.Close()
		redisCli = redis
		log.Println("redis changed", redisCli)
	})
	return redisCli, err
}

//ClientWithContext ...
func (c *Client) ClientWithContext(ctx context.Context) *Client {
	if ctx == nil {
		return c
	}
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		return c
	}
	// clone using context
	copy := c.WithContext(c.Context())
	copy.WrapProcess(func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			span := c.Tracer.StartSpan("redis", opentracing.ChildOf(parentSpan.Context()))
			tags.PeerService.Set(span, "redis")
			tags.DBType.Set(span, "redis")
			span.SetTag("db.method", cmd.Name())
			span.LogKV("cmd", cmd.Name(), cmd.Args())
			defer span.Finish()

			return oldProcess(cmd)
		}
	})
	return &Client{copy, nil, nil}
}
