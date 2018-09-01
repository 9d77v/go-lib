package gorm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/clients/etcd"
	"github.com/9d77v/go-lib/clients/jaeger"
	"github.com/jinzhu/gorm"
	opentracing "github.com/opentracing/opentracing-go"

	//postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	parentSpanGormKey = "opentracingParentSpan"
	spanGormKey       = "opentracingSpan"
	spanGormTracer    = "opentracingTracer"
)

//Client gorm client
type Client struct {
	*gorm.DB
	Tracer opentracing.Tracer
	Closer io.Closer
}

//NewClient gorm client
func NewClient(config *config.DBConfig) (*Client, error) {
	if config == nil {
		return nil, errors.New("db config is not exist")
	}
	//support postgres
	if config.Driver != "postgres" {
		return nil, errors.New("unsupport driver,now only support postgres")
	}
	//auto create database
	dbURL := fmt.Sprintf("host=%s port=%d user=%s sslmode=disable password=%s",
		config.Host, config.Port, config.User, config.Password)
	dbInit, err := gorm.Open(config.Driver, dbURL)
	if err != nil {
		return nil, err
	}
	defer dbInit.Close()
	initSQL := fmt.Sprintf("CREATE DATABASE \"%s\" WITH  OWNER =%s ENCODING = 'UTF8' CONNECTION LIMIT=-1;",
		config.Name, config.User)
	err = dbInit.Exec(initSQL).Error
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, err
	}
	dbWithNameURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		config.Host, config.Port, config.User, config.Name, config.Password)
	//global database connection
	db, err := gorm.Open(config.Driver, dbWithNameURL)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db.LogMode(config.EnableLog)
	db.DB().SetMaxIdleConns(int(config.MaxIdleConns))
	db.DB().SetMaxOpenConns(int(config.MaxOpenConns))

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	return &Client{db, nil, nil}, err
}

//NewDBConfig get config from etcd
func NewDBConfig(etcdCli *etcd.Client) (string, *config.DBConfig) {
	appName := os.Getenv("APP_NAME")
	profile := os.Getenv("PROFILE")
	dbKey := etcdCli.GetEtcdKey(profile, appName, "db")
	dbConfig := new(config.DBConfig)
	err := etcdCli.GetValue(5*time.Second, dbKey, dbConfig)
	if err != nil {
		log.Println("db config is not exist:", err)
		dbConfig = nil
	}
	return dbKey, dbConfig
}

//NewClientFromEtcd init gorm from etcd config and watch config to update gorm
func NewClientFromEtcd(etcdCli *etcd.Client, values ...interface{}) (dbCli *Client, err error) {
	dbKey, dbConfig := NewDBConfig(etcdCli)
	db, err := NewClient(dbConfig)
	if err != nil {
		log.Println("db connect failed")
	}
	db.AutoMigrate(values...)
	tracer, closer, err := jaeger.InitTracerFromEtcd(etcdCli, dbConfig.Driver)
	db.Tracer = tracer
	db.Closer = closer
	dbCli = db
	dbCli.addGormCallbacks()
	log.Println("db inited")
	//change to new db connection when  config changed
	go etcdCli.WatchKey(dbKey, dbConfig, func() {
		log.Println("db wait for change")
		db, err := NewClient(dbConfig)
		if err != nil {
			log.Println("db connect failed")
			return
		}
		db.AutoMigrate(values...)
		db.Closer = dbCli.Closer
		db.Tracer = dbCli.Tracer
		dbCli.Close()
		dbCli = db
		dbCli.addGormCallbacks()
		log.Println("db changed")
	})
	return dbCli, err
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
	c.DB = c.DB.Set(parentSpanGormKey, parentSpan)
	c.DB = c.DB.Set(spanGormTracer, c.Tracer)
	return c
}

func (c *Client) addGormCallbacks() {
	callbacks := newCallbacks()
	registerCallbacks(c.DB, "create", callbacks)
	registerCallbacks(c.DB, "query", callbacks)
	registerCallbacks(c.DB, "update", callbacks)
	registerCallbacks(c.DB, "delete", callbacks)
	registerCallbacks(c.DB, "row_query", callbacks)
}
