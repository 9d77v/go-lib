package redis

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"

	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/clients/etcd"
)

//Client redis client with opentracing
type Client struct {
	*redis.Client
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
	return &Client{client}, nil
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
	redisCli = redis
	log.Println("redis inited", redisCli)
	//change to new redis connection when  config changed
	go etcdCli.WatchKey(redisKey, redisConfig, redisCli, func() {
		redis, err := NewClient(redisConfig)
		if err != nil {
			log.Println("redis connect failed")
			return
		}
		redisCli.Close()
		redisCli = redis
		log.Println("redis changed", redisCli)
	})
	return redisCli, err
}
