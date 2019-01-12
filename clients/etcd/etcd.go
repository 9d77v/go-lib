package etcd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	yaml "gopkg.in/yaml.v2"
)

//Client .
type Client struct {
	*clientv3.Client
}

//NewClient .
func NewClient(dialTimeout time.Duration) (*Client, error) {
	etcdServer := os.Getenv("ETCD_SERVER")
	if etcdServer == "" {
		etcdServer = "localhost:7500"
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdServer},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &Client{cli}, nil
}

//GetValue Get a key
func (e *Client) GetValue(dailTimeout time.Duration, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), dailTimeout)
	resp, err := e.Get(ctx, key)
	cancel()
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return errors.New("value not found")
	}
	for _, ev := range resp.Kvs {
		if string(ev.Key) == key {
			err = yaml.Unmarshal(ev.Value, value)
			return err
		}
	}
	return nil
}

//GetEtcdKey etcd key
func (e *Client) GetEtcdKey(profile, appName, item string) string {
	return fmt.Sprintf("/config/%s/%s/%s", profile, appName, item)
}

//InitClient init client by remote changed config
type InitClient func()

//WatchKey .
func (e *Client) WatchKey(key string, config interface{}, initClient InitClient) {
	rch := e.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if string(ev.Kv.Key) == key {
				err := yaml.Unmarshal(ev.Kv.Value, config)
				if err != nil {
					log.Println("config parse fialed")
					return
				}
				initClient()
			}
		}
	}
}

//SyncKey sync yaml content to etcd
func (e *Client) SyncKey(requestTimeout time.Duration, key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := e.Get(ctx, key)
	cancel()
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		_, err := e.Put(ctx, key, string(value))
		cancel()
		if err != nil {
			return err
		}
	} else {
		for _, ev := range resp.Kvs {
			if string(ev.Key) == key && string(value) != string(ev.Value) {
				ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
				_, err := e.Put(ctx, key, string(value))
				cancel()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
